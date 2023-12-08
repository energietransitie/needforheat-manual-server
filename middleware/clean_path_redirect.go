package middleware

import (
	"net/http"
	"path"

	"github.com/go-chi/chi"
)

// CleanPathRedirect cleans the path and redirects if it changed.
func CleanPathRedirect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())

		routePath := rctx.RoutePath
		if routePath == "" {
			if r.URL.RawPath != "" {
				routePath = r.URL.RawPath
			} else {
				routePath = r.URL.Path
			}
		}

		cleanPath := path.Clean(routePath)
		cleanPath += "/"

		rctx.RoutePath = cleanPath
		r.URL.Path = cleanPath

		next.ServeHTTP(w, r)
	})
}
