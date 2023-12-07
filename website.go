package main

import (
	"context"
	"embed"
	"errors"
	"github.com/gin-gonic/gin"
	"html/template"
	"hub.lol/website/slogging"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//go:embed assets/* data/* posts/* static/* templates/*
var fs embed.FS

func main() {
	// disable logging decoration
	log.SetFlags(0)

	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(slogging.Gin)

	// default recovery handler
	router.Use(gin.Recovery())

	router.SetHTMLTemplate(template.Must(template.New("").
		ParseFS(fs, "templates/*.gohtml")))

	// go module vanity url redirects
	if err := initGoPkgs(router); err != nil {
		log.Fatalln(err.Error())
	}

	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusPermanentRedirect, "/about")
	})

	router.GET("/about", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "about.gohtml", nil)
	})

	if err := initPosts(router); err != nil {
		log.Fatalln(err.Error())
	}

	if err := initReads(router); err != nil {
		log.Fatalln(err.Error())
	}

	if err := initProjects(router); err != nil {
		log.Fatalln(err.Error())
	}

	if err := initTags(router); err != nil {
		log.Fatalln(err.Error())
	}

	router.GET("/assets/*path", func(ctx *gin.Context) {
		setCacheHeader(ctx)
		ctx.FileFromFS(ctx.Request.URL.Path, http.FS(fs))
	})

	router.GET("/static/*path", func(ctx *gin.Context) {
		setCacheHeader(ctx)
		ctx.FileFromFS(ctx.Request.URL.Path, http.FS(fs))
	})

	router.GET("/robots.txt", func(ctx *gin.Context) {
		ctx.Request.URL.Path = "/static/robots.txt"
		router.HandleContext(ctx)
	})

	router.GET("/sitemap.xml", func(ctx *gin.Context) {
		ctx.Request.URL.Path = "/static/sitemap.xml"
		router.HandleContext(ctx)
	})

	router.GET("/teapot", func(ctx *gin.Context) {
		setCacheHeader(ctx)
		ctx.String(http.StatusTeapot, "🫖")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
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
