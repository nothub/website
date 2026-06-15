package main

import (
	"context"
	"embed"
	"errors"
	"html/template"
	iofs "io/fs"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nothub/website/internal/middleware"
)

//go:embed assets/* data/* jslinux/* posts/* static/* templates/*
var fs embed.FS

func buildHandler(tmpl *template.Template, logger *slog.Logger, clacksRand *rand.Rand, trustProxy bool) http.Handler {
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
	mux.Handle("GET /shell", http.RedirectHandler("/jslinux", http.StatusMovedPermanently))

	mux.HandleFunc("GET /about", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.ExecuteTemplate(w, "about.gohtml", nil); err != nil {
			slog.Warn("template error", "err", err)
		}
	})

	staticFS, _ := iofs.Sub(fs, "static")
	assetsFS, _ := iofs.Sub(fs, "assets")
	jslinuxFS, _ := iofs.Sub(fs, "jslinux")

	mux.HandleFunc("GET /static/", func(w http.ResponseWriter, r *http.Request) {
		setCacheHeader(w)
		http.StripPrefix("/static", http.FileServerFS(staticFS)).ServeHTTP(w, r)
	})
	mux.HandleFunc("GET /assets/", func(w http.ResponseWriter, r *http.Request) {
		setCacheHeader(w)
		http.StripPrefix("/assets", http.FileServerFS(assetsFS)).ServeHTTP(w, r)
	})
	mux.HandleFunc("GET /jslinux/", func(w http.ResponseWriter, r *http.Request) {
		setCacheHeader(w)
		http.StripPrefix("/jslinux", http.FileServerFS(jslinuxFS)).ServeHTTP(w, r)
	})

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

	var handler http.Handler = middleware.Clacks(clacksRand)(mux)
	handler = middleware.TrailingSlash(handler)
	handler = middleware.Recovery(logger)(handler)
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

	clacksRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	handler := buildHandler(tmpl, slogger, clacksRand, optTrustProxy)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	go func() {
		err := srv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("graceful shutdown complete")
		} else {
			log.Fatalf("server error: %s\n", err.Error())
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	for {
		sig := <-signals
		log.Printf("signal received: %s\n", sig.String())
		switch sig {
		case syscall.SIGHUP:
			log.Fatalf("unhandled signal: %s\n", sig.String())
		case syscall.SIGINT, syscall.SIGTERM:
			log.Println("shutting down server...")
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			err := srv.Shutdown(ctx)
			cancel()
			if err != nil {
				log.Fatalf("server shutdown error: %s\n", err)
			}
			os.Exit(0)
		default:
			log.Printf("unhandled signal: %s\n", sig.String())
		}
	}
}
