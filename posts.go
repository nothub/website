package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func initPosts(router *gin.Engine) (err error) {
	log.Println("loading posts")

	dir, err := fs.ReadDir("posts")
	if err != nil {
		log.Fatalln(err.Error())
	}

	var posts = make(map[string]template.HTML)

	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}
		slug := strings.TrimSuffix(entry.Name(), ".html")
		log.Printf("loading post: %s\n", slug)

		by, err := fs.ReadFile("posts/" + entry.Name())
		if err != nil {
			log.Fatalln(err.Error())
		}

		posts[slug] = template.HTML(by)
	}

	for _, entry := range dir {
		if !entry.IsDir() {
			continue
		}
		slug := strings.TrimSuffix(entry.Name(), ".html")
		if _, ok := posts[slug]; ok {
			log.Printf("found data for post: %s\n", slug)
		}
	}

	router.GET("/posts", func(c *gin.Context) {
		c.HTML(http.StatusOK, "posts.tmpl", posts)
	})

	router.GET("/posts/*slug", func(c *gin.Context) {
		slug := strings.TrimPrefix(c.Param("slug"), "/")

		// redirect to posts overview
		if len(strings.TrimSpace(slug)) == 0 {
			c.Request.URL.Path = "/posts"
			router.HandleContext(c)
			return
		}

		if post, ok := posts[slug]; ok {
			c.HTML(http.StatusOK, "post.tmpl", post)
		} else {
			c.AbortWithStatus(http.StatusNotFound)
		}
	})

	return nil
}
