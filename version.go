package main

import (
	"errors"
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

	// TODO: GET /version → plain text version string
	_ = version

	return nil
}
