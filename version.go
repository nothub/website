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

	router.GET("/version", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "%s", bi.Main.Version)
	})

	return nil
}
