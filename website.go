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

//go:embed stuff/* templates/*
var fs embed.FS

func main() {
	gin.DisableConsoleColor()
	router := gin.Default()
	router.SetHTMLTemplate(template.Must(template.New("").
		ParseFS(fs, "templates/*.tmpl")))

	router.Use(gopkg)

	router.GET("/", func(c *gin.Context) {
		c.Request.URL.Path = "/about"
		router.HandleContext(c)
	})

	router.GET("/about", func(c *gin.Context) {
		c.HTML(http.StatusOK, "about.tmpl", nil)
	})

	router.GET("/posts/*path", func(c *gin.Context) {
		// TODO
		c.AbortWithStatus(http.StatusTeapot)
	})

	router.GET("/reads/*path", func(c *gin.Context) {
		// TODO
		c.AbortWithStatus(http.StatusTeapot)
	})

	router.GET("/projects/*path", func(c *gin.Context) {
		// TODO
		c.AbortWithStatus(http.StatusTeapot)
	})

	router.GET("/stuff/*path", func(c *gin.Context) {
		c.FileFromFS(c.Request.URL.Path, http.FS(fs))
	})

	router.GET("/rss.xml", func(c *gin.Context) {
		// TODO
		c.AbortWithStatus(http.StatusTeapot)
	})

	router.GET("/robots.txt", func(c *gin.Context) {
		c.Request.URL.Path = "/stuff/robots.txt"
		router.HandleContext(c)
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		err := srv.ListenAndServe()
		if err == nil {
		} else if err == http.ErrServerClosed {
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
