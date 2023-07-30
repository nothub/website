package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

var gopkg = func(c *gin.Context) {
	if c.Query("go-get") != "1" {
		// not a go package request
		c.Next()
		return
	}

	modPath := strings.TrimPrefix(c.Request.URL.Path, "/")
	modPath = strings.TrimSuffix(modPath, "/")
	modPath = strings.TrimSpace(modPath)

	if len(modPath) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.HTML(http.StatusOK, "gopkg.gohtml", map[string]any{
		"pkg": modPath,
	})

	c.Done()
}
