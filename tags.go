package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

// tags holds a list of internal links per tag
var tags = make(map[string][]string)

func linkTag(tag string, link string) {
	tags[tag] = append(tags[tag], link)
}

func initTags(router *gin.Engine) (err error) {
	log.Println("linking tags")

	router.GET("/tags", func(c *gin.Context) {
		c.HTML(http.StatusOK, "tags.gohtml", tags)
	})

	router.GET("/tags/*path", func(c *gin.Context) {
		path := strings.TrimSpace(strings.TrimPrefix(c.Param("path"), "/"))
		switch path {
		case "":
			// rewrite /tags/ to /tags
			c.Request.URL.Path = "/tags"
			router.HandleContext(c)
		default:
			if links, ok := tags[path]; ok {
				c.HTML(http.StatusOK, "tags.gohtml", map[string]any{
					path: links,
				})
			} else {
				c.AbortWithStatus(http.StatusNotFound)
			}
		}
	})

	return nil
}
