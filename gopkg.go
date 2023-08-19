package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

var gopkg = func(ctx *gin.Context) {
	if ctx.Query("go-get") != "1" {
		// not a go package request
		ctx.Next()
		return
	}

	modPath := strings.TrimPrefix(ctx.Request.URL.Path, "/")
	modPath = strings.TrimSuffix(modPath, "/")
	modPath = strings.TrimSpace(modPath)

	if len(modPath) == 0 {
		ctx.AbortWithStatus(http.StatusNotFound)
	}

	ctx.HTML(http.StatusOK, "gopkg.gohtml", map[string]any{
		"pkg":    modPath,
		"gopher": gopher(),
	})
}
