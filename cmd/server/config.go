package main

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/energietransitie/twomes-manual-server/parser"
	"golang.org/x/text/language"
)

const (
	SourceEnvDefault string = "./source"
)

var (
	ErrFallbackLangEnvNotSet = errors.New("environment variable TWOMES_FALLBACK_LANG was not set")
)

// Config contains the configuration for the server.
type Config struct {
	// Source is where the manuals are pulled from.
	//
	// Set by environment variable TWOMES_MANUAL_SOURCE.
	//
	// This can be a local directory or a git repository.
	// A local directory has to be a regular path (e.g. 'source' or './source').
	// A git repository has to start with https:// (e.g. 'https://github.com/energietransitie/twomes-presence-detector-firmware').
	Source fs.FS

	// FallbackLanguage sets the fallback language for when
	// a client's Accept-Language header does not contain any available language.
	//
	// Set by environment variable TWOMES_FALLBACK_LANG.
	//
	// This must be a valid language code (e.g. nl-NL or en-GB).
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
	sourceEnv, ok := os.LookupEnv("TWOMES_MANUAL_SOURCE")
	if !ok {
		sourceEnv = SourceEnvDefault
	}

	if strings.Contains(sourceEnv, "https://") {
		log.Println("using git repository", sourceEnv, "as manual source")
		return parser.NewLabRepoSource(sourceEnv, nil)
	}
	log.Println("using local directory", sourceEnv, "as manual source")
	return parser.NewLabDirSource(sourceEnv)
}

func parseFallbackLangEnv() (language.Tag, error) {
	fallbackLangEnv, ok := os.LookupEnv("TWOMES_FALLBACK_LANG")
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
