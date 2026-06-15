package main

import (
	"html/template"
	"log"
	"log/slog"
	"net/http"
)

type Ref struct {
	Name string
	Link string
}

var tags = make(map[string][]Ref)

func linkTag(tag string, name string, link string) {
	tags[tag] = append(tags[tag], Ref{Name: name, Link: link})
}

func initTags(mux *http.ServeMux, tmpl *template.Template) (err error) {
	log.Println("linking tags")

	mux.HandleFunc("GET /tags", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.ExecuteTemplate(w, "tags.gohtml", tags); err != nil {
			slog.Warn("template error", "err", err)
		}
	})

	mux.HandleFunc("GET /tags/{tag}", func(w http.ResponseWriter, r *http.Request) {
		tag := r.PathValue("tag")
		refs, ok := tags[tag]
		if !ok {
			writeError(w, r, http.StatusNotFound, tmpl)
			return
		}
		if err := tmpl.ExecuteTemplate(w, "tags.gohtml", map[string][]Ref{tag: refs}); err != nil {
			slog.Warn("template error", "err", err)
		}
	})

	return nil
}
