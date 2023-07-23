package main

import (
	"context"
	"embed"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//go:embed data/* static/* templates/*
var fs embed.FS

func main() {
	gin.DisableConsoleColor()
	router := gin.Default()
	router.SetHTMLTemplate(template.Must(template.New("").
		ParseFS(fs, "templates/*.tmpl")))

	// go module vanity url redirects
	router.Use(gopkg)

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/about")
	})

	router.GET("/about", func(c *gin.Context) {
		c.HTML(http.StatusOK, "about.tmpl", nil)
	})

	if err := initProjects(router); err != nil {
		log.Fatalln(err.Error())
	}

	if err := initReads(router); err != nil {
		log.Fatalln(err.Error())
	}

	if err := initPosts(router); err != nil {
		log.Fatalln(err.Error())
	}

	if err := initRss(router); err != nil {
		log.Fatalln(err.Error())
	}

	router.GET("/static/*path", func(c *gin.Context) {
		c.FileFromFS(c.Request.URL.Path, http.FS(fs))
	})

	router.GET("/robots.txt", func(c *gin.Context) {
		c.Request.URL.Path = "/static/robots.txt"
		router.HandleContext(c)
	})

	router.GET("/sitemap.xml", func(c *gin.Context) {
		c.Request.URL.Path = "/static/sitemap.xml"
		router.HandleContext(c)
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		err := srv.ListenAndServe()
		if err == http.ErrServerClosed {
			log.Println("graceful shutdown complete")
		} else if err == nil {
			log.Fatalln("http.Server stopped with nil error")
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
			// TODO: reload config
		case syscall.SIGINT, syscall.SIGTERM:
			log.Println("shutting down server...")
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			err := srv.Shutdown(ctx)
			if cancel != nil {
				cancel()
			}
			if err != nil {
				log.Fatal("server shutdown error: ", err)
			}
			os.Exit(0)
		default:
			log.Fatalf("unhandled signal: %s", sig.String())
		}
	}
}
