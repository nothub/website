package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func initPosts(router *gin.Engine) (err error) {
	// TODO

	router.GET("/posts/*path", func(c *gin.Context) {
		path := c.Param("path")
		log.Printf("path=%q\n", path)
		// TODO
		c.AbortWithStatus(http.StatusTeapot)
	})

	return nil
}
