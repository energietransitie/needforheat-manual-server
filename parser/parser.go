// Package parser implements types and methods to parse markdown files in a filesystem to HTML
// to a destination filesystem that can be used to serve Twomes manuals.
package parser

import (
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/energietransitie/twomes-manual-server/wfs"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

const (
	htmlTemplateFileName = "template.html"
	fallbackManualTitle  = "Twomes manual"
)

var (
	ErrTemplateNotFound = errors.New("template file could not be found")
)

// HTMLTemplate contains data for filling a template.html.
type HTMLTemplate struct {
	Language string
	Title    string
	Body     template.HTML
}

// A Parser can parse manuals written in markdown to html files.
//
// Following the folder structure specification, the parser will parse all markdown files to HTML
// while checking languages and creating a structure that can be served by a [Server].
type Parser struct {
	destFS fs.FS
}

// Create a new Parser that uses sourceFS as its filesystem to parse manuals.
func New(destFS fs.FS) *Parser {
	parser := &Parser{
		destFS: destFS,
	}

	parser.eraseDest()
	return parser
}

// Parse files following the folder structure specification
// to a filesystem that a [Server] can use to serve HTML.
//
// destFS has to be a writable filessytem.
func (p *Parser) Parse(sourceFS fs.FS) error {
	if sourceFS == nil {
		return nil
	}

	err := p.parseRecursive(sourceFS, ".")
	if err != nil {
		return err
	}

	return nil
}

// Erase the destination filesystem.
func (p *Parser) eraseDest() error {
	return wfs.RemoveAll(p.destFS, ".")
}

// Parse the source folder recursively and process each file or folder appropriately.
func (p *Parser) parseRecursive(sourceFS fs.FS, pathName string) error {
	entries, err := fs.ReadDir(sourceFS, pathName)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fullPath := path.Join(pathName, entry.Name())
		if isReadmeFile(entry) {
			continue
		} else if isMarkdownFile(entry) {
			err = p.parseMdToHTML(sourceFS, fullPath)
			if err != nil {
				return err
			}
		} else if isDetailsFile(entry) {
			err = p.getRepoManual(sourceFS, fullPath)
			if err != nil {
				return err
			}
		} else if isDisplayNamesFile(entry) {
			err = p.copyFileToDest(sourceFS, fullPath)
			if err != nil {
				return err
			}
		} else if isAssetFolder(entry) {
			err = p.copyDirToDest(sourceFS, fullPath)
			if err != nil {
				return err
			}
		} else if entry.IsDir() {
			err = p.parseRecursive(sourceFS, fullPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Parse a markdown at filepath in p.SourceFS and parse it to HTML.
//
// The generated HTML-file will be called index.html in a folder named after the language code.
// The language code is taken from the markdown file's name.
// The folder will be placed in the same spot as in p.sourceFS, except not in a language directory.
func (p *Parser) parseMdToHTML(sourceFS fs.FS, filePath string) error {
	md, err := fs.ReadFile(sourceFS, filePath)
	if err != nil {
		return err
	}

	mdParser := parser.NewWithExtensions(parser.CommonExtensions)
	doc := mdParser.Parse(md)

	htmlRenderer := html.NewRenderer(html.RendererOptions{Flags: html.CommonFlags})

	renderedHTML := markdown.Render(doc, htmlRenderer)

	templateFilePath, err := p.findTemplateFile(sourceFS, filePath)
	if err != nil {
		return err
	}

	t, err := template.New(htmlTemplateFileName).ParseFS(sourceFS, templateFilePath)
	if err != nil {
		return err
	}

	destFilePath, err := GetDestinationFilePath(sourceFS, filePath)
	if err != nil {
		return err
	}

	destinationHTMLPath := createDestinationPath(destFilePath)

	err = wfs.MkdirAll(p.destFS, path.Dir(destinationHTMLPath), fs.ModePerm)
	if err != nil {
		return err
	}

	file, err := wfs.CreateFile(p.destFS, destinationHTMLPath)
	if err != nil {
		return err
	}
	defer file.Close()

	language := strings.TrimSuffix(path.Base(filePath), ".md")
	title := findTitle(md)

	templateData := HTMLTemplate{
		Language: language,
		Title:    title,
		Body:     template.HTML(renderedHTML),
	}

	return t.Execute(file, templateData)
}

// Get the repo manuals for the device based on details.json file at filePath.
func (p *Parser) getRepoManual(sourceFS fs.FS, filePath string) error {
	file, err := fs.ReadFile(sourceFS, filePath)
	if err != nil {
		return err
	}

	details := struct {
		Repo string `json:"firmware_repository"`
	}{}
	err = json.Unmarshal(file, &details)
	if err != nil {
		return err
	}

	// TODO: possibly support authentication.
	deviceRepo, err := NewDeviceRepoSource(details.Repo, nil)
	if err != nil {
		return err
	}

	return p.Parse(deviceRepo)
}

// Return the filepath of the template file that should be used for the file at the specified filePath.
func (p *Parser) findTemplateFile(sourceFS fs.FS, filePath string) (string, error) {
	splitFilePath := strings.Split(filePath, string(os.PathSeparator))
	if len(splitFilePath) <= 1 {
		return "", ErrTemplateNotFound
	}

	parentDir := path.Join(splitFilePath[:len(splitFilePath)-1]...)
	testFilePath := path.Join(parentDir, htmlTemplateFileName)

	if !fileExists(sourceFS, testFilePath) {
		// Look for template file in the next directory up.
		return p.findTemplateFile(sourceFS, parentDir)
	}

	return testFilePath, nil
}

// Copy file at filePath from p.sourceFS to p.destFS.
func (p *Parser) copyFileToDest(sourceFS fs.FS, filePath string) error {
	sourceFile, err := sourceFS.Open(filePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFilePath, err := GetDestinationFilePath(sourceFS, filePath)
	if err != nil {
		return err
	}

	destFile, err := wfs.CreateFile(p.destFS, destFilePath)
	if errors.Is(err, fs.ErrNotExist) {
		info, err := fs.Stat(sourceFS, path.Dir(filePath))
		if err != nil {
			return err
		}
		err = wfs.MkdirAll(p.destFS, path.Dir(destFilePath), info.Mode())
		if err != nil {
			return err
		}
		destFile, err = wfs.CreateFile(p.destFS, destFilePath)
		if err != nil {
			return err
		}
	}
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// Copy dir at path from p.sourceFS to p.destFS.
func (p *Parser) copyDirToDest(sourceFS fs.FS, dirPath string) error {
	entries, err := fs.ReadDir(sourceFS, dirPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fullPath := path.Join(dirPath, entry.Name())

		if entry.IsDir() {
			p.copyDirToDest(sourceFS, fullPath)
		} else {
			p.copyFileToDest(sourceFS, fullPath)
		}
	}

	return nil
}

// Find the title of a markdown file.
//
// input is the bytes read from a markdown file.
// If a title could not be found, a fallback will be used.
func findTitle(input []byte) string {
	// We only need to replace the first occurance, so use Replace, instead of ReplaceAll.
	osIndependentInputString := strings.Replace(string(input), "\r\n", "\n", 1)
	lines := strings.Split(osIndependentInputString, "\n")
	if len(lines) <= 0 {
		return fallbackManualTitle
	}

	title, ok := strings.CutPrefix(lines[0], "# ")
	if ok {
		return strings.TrimSpace(title)
	}
	return fallbackManualTitle
}

// Returns if file at filePath exists.
func fileExists(sourceFS fs.FS, filePath string) bool {
	_, err := fs.Stat(sourceFS, filePath)
	return err == nil
}

// Returns if d is a README file.
func isReadmeFile(d fs.DirEntry) bool {
	return !d.IsDir() && d.Name() == "README.md"
}

// Returns if d is a markdown file.
func isMarkdownFile(d fs.DirEntry) bool {
	return path.Ext(d.Name()) == ".md"
}

func isDisplayNamesFile(d fs.DirEntry) bool {
	return !d.IsDir() && d.Name() == "display_names.json"
}

func isDetailsFile(d fs.DirEntry) bool {
	return !d.IsDir() && d.Name() == "details.json"
}

// Returns if d is an asset directory as decribed in the specification.
func isAssetFolder(d fs.DirEntry) bool {
	return d.IsDir() && d.Name() == "assets"
}

// Returns the destinationPath for an HTML file, based on the markdown source path.
func createDestinationPath(sourcePath string) string {
	dirPath, sourceFile := path.Split(sourcePath)

	// Remove languages directory.
	dirPath = strings.TrimSuffix(dirPath, "languages/")

	// Add folder with language code name.
	languageCode := strings.TrimSuffix(sourceFile, ".md")
	dirPath = path.Join(dirPath, languageCode)

	return path.Join(dirPath, "index.html")
}