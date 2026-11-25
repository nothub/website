package main

import (
	"crypto/sha256"
	"encoding/hex"
	"html/template"
	iofs "io/fs"
	"log/slog"
	"net/http"
	"strings"

	"github.com/elnormous/contenttype"
)

func setCacheHeader(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "public, max-age=604800, immutable")
}

func writeError(w http.ResponseWriter, r *http.Request, code int, tmpl *template.Template) {
	available := []contenttype.MediaType{
		contenttype.NewMediaType("text/html"),
	}
	w.Header().Add("Vary", "Accept")
	_, _, err := contenttype.GetAcceptableMediaType(r, available)
	if err != nil {
		http.Error(w, http.StatusText(code), code)
		return
	}
	w.WriteHeader(code)
	if execErr := tmpl.ExecuteTemplate(w, "error.gohtml", code); execErr != nil {
		slog.Warn("error template execution failed", "err", execErr)
	}
}

// revalidatingFileServer serves fsys with a content-hash ETag on every file and
// "Cache-Control: no-cache", so browsers pick up a deploy at once and get a 304
// otherwise. Embedded files have no modtime, so without the ETag every
// revalidation would download the whole file again.
func revalidatingFileServer(fsys iofs.FS) (http.Handler, error) {
	etags := make(map[string]string)
	err := iofs.WalkDir(fsys, ".", func(p string, d iofs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := iofs.ReadFile(fsys, p)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(data)
		etag := `"` + hex.EncodeToString(sum[:16]) + `"`
		etags["/"+p] = etag
		if strings.HasSuffix(p, "index.html") {
			etags["/"+strings.TrimSuffix(p, "index.html")] = etag
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	files := http.FileServerFS(fsys)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if etag, ok := etags[r.URL.Path]; ok {
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("ETag", etag)
		}
		files.ServeHTTP(w, r)
	}), nil
}
