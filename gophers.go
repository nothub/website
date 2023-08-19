package main

import (
	"log"
)

var gophers []string

func init() {
	dir, err := fs.ReadDir("static/gophers")
	if err != nil {
		log.Fatalf("unable to read gopher dir: %s\n", err.Error())
	}
	for _, file := range dir {
		if file.IsDir() {
			continue
		}
		gophers = append(gophers, "/static/gophers/"+file.Name())
	}
}

func gopher() string {
	return gophers[random.Intn(len(gophers))]
}
