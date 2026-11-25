package main

import (
	"context"
	"embed"
	"errors"
	"html/template"
	iofs "io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nothub/website/internal/middleware"
)

//go:embed data/* jslinux/* posts/* static/* templates/*
var fs embed.FS

func buildHandler(tmpl *template.Template, logger *slog.Logger, trustProxy bool) http.Handler {
	mux := http.NewServeMux()

	if err := initPosts(mux, tmpl); err != nil {
		log.Fatalln(err.Error())
	}
	if err := initProjects(mux, tmpl); err != nil {
		log.Fatalln(err.Error())
	}
	if err := initTags(mux, tmpl); err != nil {
		log.Fatalln(err.Error())
	}
	if err := initVersion(mux); err != nil {
		log.Fatalln(err.Error())
	}

	mux.Handle("GET /{$}", http.RedirectHandler("/about", http.StatusMovedPermanently))
	mux.Handle("GET /shell", http.RedirectHandler("/jslinux/", http.StatusMovedPermanently))

	mux.HandleFunc("GET /about", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.ExecuteTemplate(w, "about.gohtml", nil); err != nil {
			slog.Warn("template error", "err", err)
		}
	})

	staticFS, _ := iofs.Sub(fs, "static")
	staticServer, err := revalidatingFileServer(staticFS)
	if err != nil {
		log.Fatalln(err.Error())
	}
	jslinuxFS, _ := iofs.Sub(fs, "jslinux")
	jslinuxFiles, err := revalidatingFileServer(jslinuxFS)
	if err != nil {
		log.Fatalln(err.Error())
	}
	jslinuxServer := http.StripPrefix("/jslinux", jslinuxFiles)

	mux.Handle("GET /static/", http.StripPrefix("/static", staticServer))

	mux.HandleFunc("GET /robots.txt", func(w http.ResponseWriter, r *http.Request) {
		data, _ := fs.ReadFile("static/robots.txt")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write(data)
	})
	mux.HandleFunc("GET /sitemap.xml", func(w http.ResponseWriter, r *http.Request) {
		data, _ := fs.ReadFile("static/sitemap.xml")
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		w.Write(data)
	})

	mux.HandleFunc("GET /teapot", func(w http.ResponseWriter, r *http.Request) {
		setCacheHeader(w)
		w.WriteHeader(http.StatusTeapot)
		w.Write([]byte("🫖"))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, r, http.StatusNotFound, tmpl)
	})

	// jslinux must be served at /jslinux/ (with trailing slash) so that
	// relative URLs in its JS resolve correctly via window.location.pathname.
	// Dispatch it before TrailingSlash so the slash is never stripped.
	muxHandler := middleware.TrailingSlash(middleware.Clacks(mux))
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/jslinux" {
			http.Redirect(w, r, "/jslinux/", http.StatusMovedPermanently)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/jslinux/") {
			jslinuxServer.ServeHTTP(w, r)
			return
		}
		muxHandler.ServeHTTP(w, r)
	})

	var handler http.Handler = middleware.Recovery(logger)(inner)
	handler = middleware.RealIP(logger, trustProxy)(handler)
	handler = middleware.Logger(logger)(handler)
	return handler
}

func main() {
	log.SetFlags(0)
	initFlags()

	tmpl := template.Must(template.New("").ParseFS(fs, "templates/*.gohtml"))

	required := []string{"error.gohtml", "about.gohtml", "posts.gohtml", "post.gohtml", "projects.gohtml", "tags.gohtml"}
	for _, name := range required {
		if tmpl.Lookup(name) == nil {
			log.Fatalf("required template %q not found in templates/", name)
		}
	}

	handler := buildHandler(tmpl, slogger, optTrustProxy)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		err := srv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("graceful shutdown complete")
		} else {
			log.Fatalf("server error: %s\n", err.Error())
		}
	}()

	signal.Ignore(syscall.SIGHUP)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	sig := <-signals
	log.Printf("signal received: %s\n", sig.String())
	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	err := srv.Shutdown(ctx)
	cancel()
	if err != nil {
		log.Fatalf("server shutdown error: %s\n", err)
	}
	os.Exit(0)
}
