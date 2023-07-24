package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
	"html/template"
	"log"
	"net/http"
)

func initPosts(router *gin.Engine) (err error) {
	log.Println("rendering posts")

	var gm = goldmark.New(goldmark.WithRendererOptions(html.WithUnsafe()))

	dir, err := fs.ReadDir("posts")
	if err != nil {
		log.Fatalln(err.Error())
	}

	var posts = make(map[string]template.HTML)

	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}
		md, err := fs.ReadFile("posts/" + entry.Name())
		if err != nil {
			log.Fatalln(err.Error())
		}
		slug := entry.Name()
		// TODO: sanitize slug
		posts[slug] = renderPost(md, gm)
	}

	for _, entry := range dir {
		if !entry.IsDir() {
			continue
		}
		if slug, ok := posts[entry.Name()]; ok {
			log.Printf("found data for post: %s\n", slug)
		}
	}

	router.GET("/posts", func(c *gin.Context) {
		c.HTML(http.StatusOK, "posts.tmpl", map[string]map[string]template.HTML{
			"Posts": posts,
		})
	})

	router.GET("/posts/*path", func(c *gin.Context) {
		path := c.Param("path")
		if post, ok := posts[path]; ok {
			c.HTML(http.StatusOK, "post.tmpl", post)
		}
	})

	return nil
}

func renderPost(md []byte, gm goldmark.Markdown) template.HTML {
	// TODO:
	//   https://github.com/yuin/goldmark-meta / https://github.com/abhinav/goldmark-frontmatter
	//   https://github.com/yuin/goldmark-highlighting
	//   https://github.com/abhinav/goldmark-anchor
	//   https://github.com/abhinav/goldmark-toc

	var buf bytes.Buffer
	err := gm.Convert(md, &buf)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println(string(buf.Bytes()))

	return template.HTML(buf.Bytes())
}
