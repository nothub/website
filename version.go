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

	ver := bi.Main.Version
	for _, kv := range bi.Settings {
		switch kv.Key {
		case "vcs.modified":
			if kv.Value == "true" {
				ver = ver + "+DIRTY"
			}
		}
	}

	router.GET("/version", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "%s", ver)
	})

	return nil
}
