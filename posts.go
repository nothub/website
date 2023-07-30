package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuin/goldmark"
	gmmeta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/toc"
)

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
			&toc.Extender{Title: "ToC", MaxDepth: 1},
		),
	)

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

		// TODO:
		//   https://github.com/yuin/goldmark-highlighting
		//   https://github.com/abhinav/goldmark-anchor

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

		posts[slug] = template.HTML(buildHeader(meta) + "<hr>\n" + buf.String())
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

func buildHeader(postMeta *metaData) string {
	// TODO: do this in a sane way with goldmark ast transformer
	str := strings.Builder{}
	str.WriteString("<section>\n")
	str.WriteString(fmt.Sprintf("<h1>%s</h1>\n", postMeta.title))
	str.WriteString("<time>" + postMeta.date.Format(time.DateOnly) + "</time>\n")
	str.WriteString("[ ")
	var tagLinks []string
	for _, tag := range postMeta.tags {
		tagLinks = append(tagLinks, "<a class=\"tag\" href=\"/tags/"+tag+"\">"+tag+"</a>")
	}
	str.WriteString(strings.Join(tagLinks, " "))
	str.WriteString(" ]\n")
	str.WriteString("<p>" + postMeta.desc + "</p>\n")
	str.WriteString("</section>\n")
	return str.String()
}

type metaData struct {
	title string
	desc  string
	date  time.Time
	tags  []string
	draft bool
}

func parseMeta(meta map[string]any) (*metaData, error) {
	var data metaData

	if err := validateMetaEntry[string]("title", meta); err != nil {
		return nil, err
	}
	data.title = meta["title"].(string)

	if err := validateMetaEntry[string]("description", meta); err != nil {
		return nil, err
	}
	data.desc = meta["description"].(string)

	if err := validateMetaEntry[string]("date", meta); err != nil {
		return nil, err
	}
	dateVal, err := time.Parse(time.DateOnly /* time.DateOnly */, fmt.Sprint(meta["date"]))
	if err != nil {
		return nil, errors.New("meta contains wrong date format: " + err.Error())
	}
	data.date = dateVal

	if _, ok := meta["tags"]; !ok {
		return nil, errors.New("meta entry tags is missing")
	}
	rawTags := meta["tags"].([]any)
	for _, rawTag := range rawTags {
		_, typeOk := rawTag.(string)
		if !typeOk {
			return nil, errors.New(fmt.Sprintf("meta entry tag has wrong type"))
		}
		data.tags = append(data.tags, rawTag.(string))
	}

	if err := validateMetaEntry[bool]("draft", meta); err != nil {
		return nil, err
	}
	data.draft = meta["draft"].(bool)

	return &data, nil
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
