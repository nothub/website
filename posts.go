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

	chroma "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/gin-gonic/gin"

	gmfigure "github.com/mangoumbrella/goldmark-figure"
	"github.com/yuin/goldmark"
	gmhl "github.com/yuin/goldmark-highlighting/v2"
	gmmeta "github.com/yuin/goldmark-meta"
	gmext "github.com/yuin/goldmark/extension"
	gmparser "github.com/yuin/goldmark/parser"
	gmanchor "go.abhg.dev/goldmark/anchor"
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

type anchorTexter struct{}

func (*anchorTexter) AnchorText(h *gmanchor.HeaderInfo) []byte {
	if h.Level > 2 {
		return nil
	}
	return []byte("¶")
}

func initPosts(router *gin.Engine) (err error) {
	log.Println("loading posts")

	gm := goldmark.New(goldmark.WithParserOptions(gmparser.WithAutoHeadingID()), goldmark.WithExtensions(
		gmext.Footnote,
		gmext.Strikethrough,
		gmext.Table,
		gmfigure.Figure,
		gmmeta.New(),
		&gmanchor.Extender{Texter: &anchorTexter{}},
		gmhl.NewHighlighting(
			gmhl.WithStyle("gruvbox"),
			gmhl.WithFormatOptions(chroma.WithLineNumbers(true)),
		)))

	dir, err := fs.ReadDir("posts")
	if err != nil {
		log.Fatalln(err.Error())
	}

	var posts = make(map[string]Post)

	for _, entry := range dir {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			log.Printf("skipping %s\n", entry.Name())
			continue
		}

		slug := strings.TrimSuffix(entry.Name(), ".md")

		byts, err := fs.ReadFile("posts/" + entry.Name())
		if err != nil {
			log.Fatalln(err.Error())
		}

		var buf bytes.Buffer
		ctx := gmparser.NewContext()
		err = gm.Convert(byts, &buf, gmparser.WithContext(ctx))
		if err != nil {
			log.Fatalln(err.Error())
		}

		meta, err := parseMeta(gmmeta.Get(ctx))
		if err != nil {
			log.Fatalln(err.Error())
		}

		if !meta.Draft || optLoadDrafts {
			log.Printf("registering post: %s\n", slug)
			posts[slug] = Post{
				Meta:    meta,
				Content: template.HTML(buf.String()),
			}
		} else {
			log.Printf("skipping draft: %s\n", slug)
		}
	}

	// TODO: sort posts descending by date

	for slug, post := range posts {
		for _, tag := range post.Meta.Tags {
			linkTag(tag, "Post: "+post.Meta.Title, "/posts/"+slug)
		}
	}

	router.GET("/posts", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "posts.gohtml", posts)
	})

	router.GET("/posts/*path", func(ctx *gin.Context) {
		path := strings.TrimSpace(strings.TrimPrefix(ctx.Param("path"), "/"))
		switch path {
		case "":
			// rewrite /posts/ to /posts
			ctx.Request.URL.Path = "/posts"
			router.HandleContext(ctx)
		case "rss.xml":
			// redirect /posts/rss.xml to /rss.xml
			ctx.Redirect(http.StatusPermanentRedirect, "/rss.xml")
		default:
			if post, ok := posts[path]; ok {
				ctx.HTML(http.StatusOK, "post.gohtml", post)
			} else {
				ctx.AbortWithStatus(http.StatusNotFound)
			}
		}
	})

	router.GET("/rss.xml", func(ctx *gin.Context) {
		// TODO
		ctx.AbortWithStatus(http.StatusTeapot)
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
