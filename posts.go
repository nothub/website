package main

import (
	"bytes"
	"errors"
	"fmt"
	"go.abhg.dev/goldmark/anchor"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	chroma "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/gin-gonic/gin"
	"github.com/yuin/goldmark"
	gmhl "github.com/yuin/goldmark-highlighting/v2"
	gmmeta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

type Post struct {
	Meta    Meta
	Content template.HTML
}

type Meta struct {
	Title string
	Desc  string
	Date  time.Time
	Tags  []string
	Draft bool
}

func (meta Meta) DateString() string {
	return meta.Date.Format(time.DateOnly)
}

func initPosts(router *gin.Engine) (err error) {
	log.Println("loading posts")

	gm := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithExtensions(
			gmmeta.New(),
			extension.Strikethrough,
			extension.Footnote,
			&anchor.Extender{},
			gmhl.NewHighlighting(
				gmhl.WithStyle("gruvbox"),
				gmhl.WithFormatOptions(
					chroma.WithLineNumbers(true),
				),
			),
		),
	)

	dir, err := fs.ReadDir("posts")
	if err != nil {
		log.Fatalln(err.Error())
	}

	var posts = make(map[string]Post)

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
		context := parser.NewContext()
		err = gm.Convert(md, &buf, parser.WithContext(context))
		if err != nil {
			log.Fatalln(err.Error())
		}

		meta, err := parseMeta(gmmeta.Get(context))
		if err != nil {
			log.Fatalln(err.Error())
		}
		log.Printf("%++v\n", meta)

		posts[slug] = Post{
			Meta:    meta,
			Content: template.HTML(buf.String()),
		}
	}

	// TODO: sort posts descending by date

	for slug, post := range posts {
		for _, tag := range post.Meta.Tags {
			linkTag(tag, "Post: "+post.Meta.Title, "/posts/"+slug)
		}
	}

	router.GET("/posts", func(c *gin.Context) {
		c.HTML(http.StatusOK, "posts.gohtml", posts)
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
				c.HTML(http.StatusOK, "post.gohtml", post)
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

func parseMeta(rawMeta map[string]any) (meta Meta, err error) {
	if err = validateMetaEntry[string]("title", rawMeta); err != nil {
		return meta, err
	}
	meta.Title = rawMeta["title"].(string)

	if err = validateMetaEntry[string]("description", rawMeta); err != nil {
		return meta, err
	}
	meta.Desc = rawMeta["description"].(string)

	if err = validateMetaEntry[string]("date", rawMeta); err != nil {
		return meta, err
	}
	dateVal, err := time.Parse(time.DateOnly, fmt.Sprint(rawMeta["date"]))
	if err != nil {
		return meta, errors.New("meta contains wrong date format: " + err.Error())
	}
	meta.Date = dateVal

	if _, ok := rawMeta["tags"]; !ok {
		return meta, errors.New(fmt.Sprintf("meta entry %q is missing", "tags"))
	}
	rawTags := rawMeta["tags"].([]any)
	for val, rawTag := range rawTags {
		_, typeOk := rawTag.(string)
		if !typeOk {
			return meta, errors.New(fmt.Sprintf("tag %q has wrong type", val))
		}
		meta.Tags = append(meta.Tags, rawTag.(string))
	}

	if err := validateMetaEntry[bool]("draft", rawMeta); err != nil {
		return meta, err
	}
	meta.Draft = rawMeta["draft"].(bool)

	return meta, nil
}

func validateMetaEntry[T any](name string, meta map[string]any) error {
	if val, ok := meta[name]; ok {
		_, typeOk := val.(T)
		if !typeOk {
			return errors.New(fmt.Sprintf("meta entry %q has wrong type", name))
		}
	} else {
		return errors.New(fmt.Sprintf("meta entry %q is missing", name))
	}
	return nil
}
