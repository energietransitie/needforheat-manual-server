package needforheatmanualserver

import (
	"errors"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/go-chi/chi"
	"golang.org/x/text/language"
)

const (
	genericCampaign    string = "generic"
	manufacturerManual string = "manufacturer"
)

type ServerOptions struct {
	FallbackLanguage language.Tag
}

// A Server is a wrapper for an HTTP server that serves manuals from a filesystem.
//
// The Chi library is used and is fully compatible with net/http.
type Server struct {
	*chi.Mux
	fsys    fs.FS
	options ServerOptions
}

// Create a new server that uses fsys as its filesystem to serve manuals.
func NewServer(fsys fs.FS, options ServerOptions) *Server {
	r := chi.NewRouter()

	server := &Server{
		Mux:     r,
		fsys:    fsys,
		options: options,
	}

	r.Handle("/campaigns/{manual_type_name}/", Handler(server.handleCampaignGenericRedirect))

	r.Handle("/campaigns/{campaign_name}/{manual_type_name}/", Handler(server.handleLanguageRedirect))

	r.Handle("/campaigns/{campaign_name}/{manual_type_name}/*", http.FileServer(http.FS(server.fsys)))

	r.Handle("/devices/{device_type_name}/", Handler(server.handleDisplayName))

	r.Handle("/devices/{device_type_name}/{manual_type_name}/", Handler(server.handleDeviceGenericRedirect))

	languageRedirectWithManufacturerFallback := manufacturerFallbackMiddleware(server.handleLanguageRedirect)
	r.Handle("/devices/{device_type_name}/{manual_type_name}/{campaign_name}/", languageRedirectWithManufacturerFallback)

	r.Handle("/devices/{device_type_name}/{manual_type_name}/{campaign_name}/*", http.FileServer(http.FS(server.fsys)))

	//EnergyQuery
	r.Handle("/energy_queries/{energy_query_type_name}/", Handler(server.handleDisplayName))

	r.Handle("/energy_queries/{energy_query_type_name}/{manual_type_name}/", Handler(server.handleDeviceGenericRedirect))

	r.Handle("/energy_queries/{energy_query_type_name}/{manual_type_name}/{campaign_name}/", languageRedirectWithManufacturerFallback)

	r.Handle("/energy_queries/{energy_query_type_name}/{manual_type_name}/{campaign_name}/*", http.FileServer(http.FS(server.fsys)))

	//Cloud_feeds
	r.Handle("/cloud_feeds/{cloud_feed_type_name}/", Handler(server.handleDisplayName))

	r.Handle("/cloud_feeds/{cloud_feed_type_name}/{manual_type_name}/", Handler(server.handleDeviceGenericRedirect))

	r.Handle("/cloud_feeds/{cloud_feed_type_name}/{manual_type_name}/{campaign_name}/", languageRedirectWithManufacturerFallback)

	r.Handle("/cloud_feeds/{cloud_feed_type_name}/{manual_type_name}/{campaign_name}/*", http.FileServer(http.FS(server.fsys)))

	return server
}

func (s *Server) handleCampaignGenericRedirect(w http.ResponseWriter, r *http.Request) error {
	urlPath := strings.Trim(r.URL.Path, "/")

	splitPath := strings.Split(urlPath, "/")

	if len(splitPath) != 2 {
		// Invalid URL.
		return NewHandlerError(nil, http.StatusNotFound)
	}

	splitPath = insertInSlice(splitPath, 1, genericCampaign)

	redirectPath := "/" + path.Join(splitPath...) + "/"

	http.Redirect(w, r, redirectPath, http.StatusFound)
	return nil
}

// Handle serving display_names.json for a requested device.
func (s *Server) handleDisplayName(w http.ResponseWriter, r *http.Request) error {
	urlPath := strings.Trim(r.URL.Path, "/")

	file, err := fs.ReadFile(s.fsys, urlPath+"/display_names.json")
	if err != nil {
		return NewHandlerError(err, http.StatusNotFound)
	}

	w.Write(file)
	return nil
}

func (s *Server) handleDeviceGenericRedirect(w http.ResponseWriter, r *http.Request) error {
	urlPath := strings.Trim(r.URL.Path, "/")

	redirectPath := "/" + path.Join(urlPath, genericCampaign) + "/"

	http.Redirect(w, r, redirectPath, http.StatusFound)
	return nil
}

// Handle redirection to correct language based on Accept-Language header.
func (s *Server) handleLanguageRedirect(w http.ResponseWriter, r *http.Request) error {
	urlPath := strings.Trim(r.URL.Path, "/")
	availableLangs, err := ParseLanguageFiles(s.fsys, urlPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return NewHandlerError(err, http.StatusNotFound)
		}

		return NewHandlerError(err, http.StatusInternalServerError)
	}

	acceptLang := r.Header.Get("Accept-Language")

	lang, err := ChooseFile(availableLangs, s.options.FallbackLanguage, acceptLang)
	if err != nil {
		return NewHandlerError(err, http.StatusInternalServerError)
	}

	redirectPath := "/" + path.Join(urlPath, lang) + "/"

	http.Redirect(w, r, redirectPath, http.StatusFound)
	return nil
}

// Middleware that will fallback from 'generic' to 'campaign' when generic was not found.
func manufacturerFallbackMiddleware(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		err := next(w, r)
		if err != nil {
			if err, ok := err.(*HandlerError); ok {
				if err.Code != http.StatusNotFound {
					return err
				}

				urlPath := strings.Trim(r.URL.Path, "/")

				splitURLPath := strings.Split(urlPath, "/")
				splitURLPath[len(splitURLPath)-1] = manufacturerManual

				redirectPath := "/" + path.Join(splitURLPath...) + "/"
				http.Redirect(w, r, redirectPath, http.StatusFound)
			}

			return err
		}
		return nil
	}
}

// Insert val into slice at index i.
// All elements from index i and up will be shifted.
func insertInSlice(s []string, i int, val string) []string {
	newSlice := make([]string, 0, len(s)+1)

	newSlice = append(newSlice, s[:i]...)
	newSlice = append(newSlice, val)
	newSlice = append(newSlice, s[i:]...)

	return newSlice
}
