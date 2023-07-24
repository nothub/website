package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/yuin/goldmark"
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
		slug := strings.TrimSuffix(entry.Name(), ".md")
		log.Printf("loading post: %s\n", slug)

		md, err := fs.ReadFile("posts/" + entry.Name())
		if err != nil {
			log.Fatalln(err.Error())
		}

		var buf bytes.Buffer
		err = goldmark.Convert(md, &buf)
		if err != nil {
			log.Fatalln(err.Error())
		}

		posts[slug] = template.HTML(buf.Bytes())
	}

	router.GET("/posts", func(c *gin.Context) {
		c.HTML(http.StatusOK, "posts.tmpl", posts)
	})

	router.GET("/posts/*path", func(c *gin.Context) {
		path := strings.TrimSpace(strings.TrimPrefix(c.Param("path"), "/"))
		switch path {
		case "":
			// rewrite /posts/ to /posts
			c.Request.URL.Path = "/posts"
			router.HandleContext(c)
		case "rss.xml":
			// redirect /posts/rss.xml to /rss.xml
			c.Redirect(http.StatusPermanentRedirect, "/rss.xml")
		default:
			if post, ok := posts[path]; ok {
				c.HTML(http.StatusOK, "post.tmpl", post)
			} else {
				c.AbortWithStatus(http.StatusNotFound)
			}
		}
	})

	router.GET("/rss.xml", func(c *gin.Context) {
		// TODO
		c.AbortWithStatus(http.StatusTeapot)
	})

	return nil
}
