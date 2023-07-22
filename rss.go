package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func initRss(router *gin.Engine) (err error) {
	// TODO

	router.GET("/rss.xml", func(c *gin.Context) {
		// TODO
		c.AbortWithStatus(http.StatusTeapot)
	})

	return nil
}
