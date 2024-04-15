package needforheatmanualserver

import (
	"errors"
	"io/fs"

	"golang.org/x/text/language"
)

var (
	ErrFallbackInvalid = errors.New("fallback is invalid")
)

// Parse all the available languages of files in a folder.
func ParseLanguageFiles(fsys fs.FS, folder string) ([]language.Tag, error) {
	var langs []language.Tag

	entries, err := fs.ReadDir(fsys, folder)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		tag, err := language.Parse(entry.Name())
		if err != nil {
			// Not a valid language tag.
			continue
		}

		langs = append(langs, tag)
	}

	return langs, nil
}

// Choose which file to show from options based on the Language-Accept header.
// Fallback will be chosen when the Language-Accept header does not contain an available language.
//
// The file name will be returned.
func ChooseFile(options []language.Tag, fallback language.Tag, acceptLangHeader string) (string, error) {
	if optionsContainFallback(options, fallback) {
		options = setCorrectFallbackOrder(options, fallback)
	}

	matcher := language.NewMatcher(options)

	tag, _ := language.MatchStrings(matcher, acceptLangHeader)

	return tag.String(), nil
}

func optionsContainFallback(options []language.Tag, fallback language.Tag) bool {
	var present bool

	for _, tag := range options {
		if tag == fallback {
			present = true
			break
		}
	}

	return present
}

func setCorrectFallbackOrder(options []language.Tag, fallback language.Tag) []language.Tag {
	if options[0] == fallback {
		return options
	}

	for i := range options {
		if options[i] == fallback {
			// Swap i and first index in slice.
			options[0], options[i] = options[i], options[0]
			break
		}
	}

	return options
}
