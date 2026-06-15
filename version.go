package main

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
)

func initVersion(mux *http.ServeMux) (err error) {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return errors.New("unable to read build info from binary")
	}

	version := bi.Main.Version
	dirty := false

	for _, kv := range bi.Settings {
		switch kv.Key {
		case "vcs.revision":
			version = kv.Value
		case "vcs.modified":
			if kv.Value == "true" {
				dirty = true
			}
		}
	}
	if dirty {
		version = version + "+dirty"
	}

	mux.HandleFunc("GET /version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(w, version)
	})

	return nil
}
