package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

type Ref struct {
	Name string
	Link string
}

var tags = make(map[string][]Ref)

func linkTag(tag string, name string, link string) {
	tags[tag] = append(tags[tag], Ref{Name: name, Link: link})
}

func initTags(router *gin.Engine) (err error) {
	log.Println("linking tags")

	router.GET("/tags", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "tags.gohtml", tags)
	})

	router.GET("/tags/*path", func(ctx *gin.Context) {
		path := strings.TrimSpace(strings.TrimPrefix(ctx.Param("path"), "/"))
		switch path {
		case "":
			// rewrite /tags/ to /tags
			ctx.Request.URL.Path = "/tags"
			router.HandleContext(ctx)
		default:
			if links, ok := tags[path]; ok {
				ctx.HTML(http.StatusOK, "tags.gohtml", map[string]any{
					path: links,
				})
			} else {
				ctx.AbortWithStatus(http.StatusNotFound)
			}
		}
	})

	return nil
}
