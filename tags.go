package main

import (
	"log"
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

func initTags(mux *http.ServeMux) (err error) {
	log.Println("linking tags")

	// TODO: GET /tags → render tags.gohtml with all tags
	// TODO: GET /tags/{tag} → render tags.gohtml filtered by tag; redirect /tags/ to /tags; 404 for unknown tags

	return nil
}
