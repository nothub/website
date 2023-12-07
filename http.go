package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

func setCacheHeader(ctx *gin.Context) {
	ctx.Writer.Header().Set("Cache-Control", "public, max-age=604800, immutable")
}
