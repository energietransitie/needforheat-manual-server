package parser

import (
	"io/fs"
	"os"
)

// LabDirSource is a local directory that contains manuals made by a lab.
type LabDirSource struct {
	fs.FS
}

// Create a new source filesystem from a directory at path.
func NewLabDirSource(path string) (fs.FS, error) {
	return LabDirSource{os.DirFS(path)}, nil
}

// Get the path to copy a file to at the destination filesystem.
// Use filePath at sourceFS to determine the path the file should be copied to at the destination filesystem.
func (dir LabDirSource) GetDestinationFilePath(filePath string) string {
	return filePath
}

// Get the path to copy a directory to at the destination filesystem.
// Use dirPath at sourceFS to determine the path the directory should be copied to at the destination filesystem.
func (dir LabDirSource) GetDestinationDirPath(dirPath string) string {
	return dirPath
}
