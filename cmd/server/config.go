package main

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/energietransitie/needforheat-manual-server/parser"
	"golang.org/x/text/language"
)

const (
	SourceEnvDefault string = "./source"
)

var (
	ErrFallbackLangEnvNotSet = errors.New("environment variable NFH_FALLBACK_LANG was not set")
)

// Config contains the configuration for the server.
type Config struct {
	// Source is where the manuals are pulled from.
	//
	// Set by environment variable NFH_MANUAL_SOURCE.
	//
	// This can be a local directory or a git repository.
	// A local directory has to be a regular path (e.g. 'source' or './source').
	// A git repository has to start with https:// and end with .git (e.g. 'https://github.com/energietransitie/twomes-presence-detector-firmware.git').
	//
	// The default branch is used, unless you set NFH_MANUAL_SOURCE_BRANCH to a branch name.
	Source fs.FS

	// FallbackLanguage sets the fallback language for when
	// a client's Accept-Language header does not contain any available language.
	//
	// Set by environment variable NFH_FALLBACK_LANG.
	//
	// This must be a valid language code (e.g. nl-NL or en-US).
	FallbackLanguage language.Tag
}

// Get configuration from environment variables.
// An error is returned if an environment variable is not set
// and there is no default setting for it or if a setting was invalid.
func getConfig() (*Config, error) {
	source, err := parseSourceEnv()
	if err != nil {
		return nil, err
	}

	fallbackLang, err := parseFallbackLangEnv()
	if err != nil {
		return nil, err
	}

	return &Config{
		Source:           source,
		FallbackLanguage: fallbackLang,
	}, nil
}

func parseSourceEnv() (fs.FS, error) {
	sourceEnv, ok := os.LookupEnv("NFH_MANUAL_SOURCE")
	if !ok {
		sourceEnv = SourceEnvDefault
	}

	sourceBranchEnv := os.Getenv("NFH_MANUAL_SOURCE_BRANCH")

	if strings.Contains(sourceEnv, "https://") {
		log.Println("using git repository", sourceEnv, "as manual source")
		if sourceBranchEnv != "" {
			log.Println("using branch", sourceBranchEnv)
		}
		return parser.NewLabRepoSource(sourceEnv, sourceBranchEnv, nil)
	}
	log.Println("using local directory", sourceEnv, "as manual source")
	return parser.NewLabDirSource(sourceEnv)
}

func parseFallbackLangEnv() (language.Tag, error) {
	fallbackLangEnv, ok := os.LookupEnv("NFH_FALLBACK_LANG")
	if !ok {
		return language.Tag{}, ErrFallbackLangEnvNotSet
	}

	lang, err := language.Parse(fallbackLangEnv)
	if err != nil {
		return language.Tag{}, err
	}

	log.Println("using", fallbackLangEnv, "as fallback language")
	return lang, nil
}
