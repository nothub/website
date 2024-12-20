package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
)

func initVersion(router *gin.Engine) (err error) {

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

	router.GET("/version", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "%s", version)
	})

	return nil
}
