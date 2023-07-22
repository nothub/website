package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func initReads(router *gin.Engine) (err error) {
	// TODO

	router.GET("/reads/*path", func(c *gin.Context) {
		path := c.Param("path")
		log.Printf("path=%q\n", path)
		// TODO
		c.AbortWithStatus(http.StatusTeapot)
	})

	return nil
}
