package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func initPosts(router *gin.Engine) (err error) {
	// TODO

	// https://github.com/yuin/goldmark
	// https://github.com/yuin/goldmark-meta / https://github.com/abhinav/goldmark-frontmatter
	// https://github.com/yuin/goldmark-highlighting
	// https://github.com/abhinav/goldmark-anchor
	// https://github.com/abhinav/goldmark-toc

	router.GET("/posts/*path", func(c *gin.Context) {
		path := c.Param("path")
		log.Printf("path=%q\n", path)
		// TODO
		c.AbortWithStatus(http.StatusTeapot)
	})

	return nil
}
