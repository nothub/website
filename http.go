package main

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/elnormous/contenttype"
)

func setCacheHeader(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "public, max-age=604800, immutable")
}

func writeError(w http.ResponseWriter, r *http.Request, code int, tmpl *template.Template) {
	available := []contenttype.MediaType{
		contenttype.NewMediaType("text/html"),
	}
	_, _, err := contenttype.GetAcceptableMediaType(r, available)
	if err != nil {
		http.Error(w, http.StatusText(code), code)
		return
	}
	w.WriteHeader(code)
	if execErr := tmpl.ExecuteTemplate(w, "error.gohtml", code); execErr != nil {
		slog.Warn("error template execution failed", "err", execErr)
	}
}
