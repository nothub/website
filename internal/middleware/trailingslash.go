package middleware

import (
	"net/http"
	"net/url"
	"path"
	"strings"
)

func TrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && strings.HasSuffix(r.URL.Path, "/") {
			// This runs before the mux cleans the path. Cleaning here collapses
			// a leading "//", which would otherwise make the Location a
			// protocol-relative URL pointing at another host.
			target := url.URL{Path: path.Clean(r.URL.Path), RawQuery: r.URL.RawQuery}
			http.Redirect(w, r, target.String(), http.StatusMovedPermanently)
			return
		}
		next.ServeHTTP(w, r)
	})
}
