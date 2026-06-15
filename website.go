package main

import (
	"context"
	"embed"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/elnormous/contenttype"
)

//go:embed assets/* data/* jslinux/* posts/* static/* templates/*
var fs embed.FS

func main() {
	// disable logging decoration
	log.SetFlags(0)

	mux := http.NewServeMux()

	// TODO: slog middleware
	// TODO: recovery middleware
	// TODO: setClacksHeader middleware
	// TODO: set up HTML templates

	// TODO: GET / → redirect to /about
	// TODO: GET /shell → redirect to /jslinux
	// TODO: GET /about → render about.gohtml

	if err := initPosts(mux); err != nil {
		log.Fatalln(err.Error())
	}

	if err := initVersion(mux); err != nil {
		log.Fatalln(err.Error())
	}

	if err := initProjects(mux); err != nil {
		log.Fatalln(err.Error())
	}

	if err := initTags(mux); err != nil {
		log.Fatalln(err.Error())
	}

	// TODO: GET /assets/ → serve embedded FS with cache header
	// TODO: GET /static/ → serve embedded FS with cache header
	// TODO: GET /jslinux/ → serve embedded FS with cache header
	// TODO: GET /robots.txt → serve static/robots.txt
	// TODO: GET /sitemap.xml → serve static/sitemap.xml
	// TODO: GET /teapot → 418

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
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
			// TODO: reload config & data
			log.Fatalf("unhandled signal: %s\n", sig.String())
		case syscall.SIGINT, syscall.SIGTERM:
			log.Println("shutting down server...")
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			err := srv.Shutdown(ctx)
			if cancel != nil {
				cancel()
			}
			if err != nil {
				log.Fatalf("server shutdown error: %s\n", err)
			}
			os.Exit(0)
		default:
			log.Printf("unhandled signal: %s\n", sig.String())
		}
	}
}
