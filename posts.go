package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	iofs "io/fs"
	"log"
	"log/slog"
	"net/http"
	"slices"
	"time"

	chroma "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/elnormous/contenttype"
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

type PostEntry struct {
	Slug string
	Post
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

func newGoldmark() goldmark.Markdown {
	return goldmark.New(goldmark.WithParserOptions(gmparser.WithAutoHeadingID()), goldmark.WithExtensions(
		gmext.Footnote,
		gmext.Strikethrough,
		gmext.Table,
		gmfigure.Figure,
		gmmeta.New(),
		&gmanchor.Extender{Texter: &anchorTexter{}},
		gmhl.NewHighlighting(
			gmhl.WithStyle("gruvbox"),
			gmhl.WithFormatOptions(
				chroma.TabWidth(4),
			),
		)))
}

func loadPosts(fsys iofs.FS, gm goldmark.Markdown, loadDrafts bool) ([]PostEntry, error) {
	dir, err := iofs.ReadDir(fsys, "posts")
	if err != nil {
		return nil, err
	}

	posts := make(map[string]Post)

	for _, entry := range dir {
		if !entry.IsDir() {
			log.Printf("skipping non-directory entry in posts/: %s\n", entry.Name())
			continue
		}

		slug := entry.Name()

		byts, err := iofs.ReadFile(fsys, "posts/"+slug+"/index.md")
		if err != nil {
			log.Printf("skipping posts/%s: no index.md (%s)\n", slug, err)
			continue
		}

		var buf bytes.Buffer
		ctx := gmparser.NewContext()
		if err := gm.Convert(byts, &buf, gmparser.WithContext(ctx)); err != nil {
			return nil, err
		}

		meta, err := parseMeta(gmmeta.Get(ctx))
		if err != nil {
			return nil, err
		}

		if !meta.Draft || loadDrafts {
			log.Printf("registering post: %s\n", slug)
			posts[slug] = Post{Meta: meta, Content: template.HTML(buf.String())}
		} else {
			log.Printf("skipping draft: %s\n", slug)
		}
	}

	sorted := make([]PostEntry, 0, len(posts))
	for slug, p := range posts {
		sorted = append(sorted, PostEntry{Slug: slug, Post: p})
	}
	slices.SortFunc(sorted, func(a, b PostEntry) int {
		return b.Meta.Date.Compare(a.Meta.Date)
	})
	return sorted, nil
}

func initPosts(mux *http.ServeMux, tmpl *template.Template) (err error) {
	log.Println("loading posts")

	sorted, err := loadPosts(fs, newGoldmark(), optLoadDrafts)
	if err != nil {
		log.Fatalln(err.Error())
	}

	posts := make(map[string]Post, len(sorted))
	for _, e := range sorted {
		posts[e.Slug] = e.Post
		for _, tag := range e.Meta.Tags {
			linkTag(tag, "Post: "+e.Meta.Title, "/posts/"+e.Slug)
		}
	}

	mux.HandleFunc("GET /posts", func(w http.ResponseWriter, r *http.Request) {
		available := []contenttype.MediaType{
			contenttype.NewMediaType("text/html"),
			contenttype.NewMediaType("application/rss+xml"),
			contenttype.NewMediaType("application/atom+xml"),
		}
		accepted, _, _ := contenttype.GetAcceptableMediaType(r, available)
		mt := accepted.Type + "/" + accepted.Subtype
		if mt == "application/rss+xml" || mt == "application/atom+xml" {
			serveRSS(w, sorted)
			return
		}
		if err := tmpl.ExecuteTemplate(w, "posts.gohtml", sorted); err != nil {
			slog.Warn("template error", "tmpl", "posts.gohtml", "err", err)
		}
	})

	mux.Handle("GET /posts/rss.xml", http.RedirectHandler("/rss.xml", http.StatusMovedPermanently))

	mux.HandleFunc("GET /rss.xml", func(w http.ResponseWriter, r *http.Request) {
		serveRSS(w, sorted)
	})

	mux.HandleFunc("GET /posts/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		p, ok := posts[slug]
		if !ok {
			writeError(w, r, http.StatusNotFound, tmpl)
			return
		}
		available := []contenttype.MediaType{
			contenttype.NewMediaType("text/html"),
			contenttype.NewMediaType("text/markdown"),
		}
		accepted, _, _ := contenttype.GetAcceptableMediaType(r, available)
		if accepted.Type+"/"+accepted.Subtype == "text/markdown" {
			raw, err := fs.ReadFile("posts/" + slug + "/index.md")
			if err != nil {
				writeError(w, r, http.StatusInternalServerError, tmpl)
				return
			}
			w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
			w.Write(raw)
			return
		}
		if err := tmpl.ExecuteTemplate(w, "post.gohtml", p); err != nil {
			slog.Warn("template error", "tmpl", "post.gohtml", "err", err)
		}
	})

	mux.HandleFunc("GET /posts/{slug}/", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		if _, ok := posts[slug]; !ok {
			writeError(w, r, http.StatusNotFound, tmpl)
			return
		}
		postFS, err := iofs.Sub(fs, "posts/"+slug)
		if err != nil {
			writeError(w, r, http.StatusNotFound, tmpl)
			return
		}
		file := r.URL.Path[len("/posts/"+slug+"/"):]
		if file == "" || file == "index.md" {
			writeError(w, r, http.StatusNotFound, tmpl)
			return
		}
		http.StripPrefix("/posts/"+slug, http.FileServerFS(postFS)).ServeHTTP(w, r)
	})

	return nil
}

func serveRSS(w http.ResponseWriter, posts []PostEntry) {
	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	w.Write([]byte(buildRSS(posts)))
}

func buildRSS(posts []PostEntry) string {
	var b bytes.Buffer
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	b.WriteString("\n")
	b.WriteString(`<rss version="2.0">`)
	b.WriteString("\n")
	b.WriteString(`  <!-- In memory of Aaron Swartz (1986–2013) — https://www.aaronsw.com -->`)
	b.WriteString("\n")
	b.WriteString("  <channel>\n")
	b.WriteString("    <title>hub.lol</title>\n")
	b.WriteString("    <link>https://hub.lol</link>\n")
	for _, e := range posts {
		b.WriteString("    <item>\n")
		b.WriteString("      <title>" + xmlEscape(e.Meta.Title) + "</title>\n")
		b.WriteString("      <link>https://hub.lol/posts/" + e.Slug + "</link>\n")
		b.WriteString("      <description>" + xmlEscape(e.Meta.Desc) + "</description>\n")
		b.WriteString("      <pubDate>" + e.Meta.Date.Format(time.RFC1123Z) + "</pubDate>\n")
		b.WriteString("    </item>\n")
	}
	b.WriteString("  </channel>\n")
	b.WriteString("</rss>\n")
	return b.String()
}

func xmlEscape(s string) string {
	var b bytes.Buffer
	for _, c := range s {
		switch c {
		case '&':
			b.WriteString("&amp;")
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		case '"':
			b.WriteString("&quot;")
		case '\'':
			b.WriteString("&apos;")
		default:
			b.WriteRune(c)
		}
	}
	return b.String()
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
			return meta, errors.New(fmt.Sprintf("tag at index %d has wrong type", val))
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
