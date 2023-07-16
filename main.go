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

//go:embed assets/* templates/*
var fs embed.FS

func main() {
	gin.DisableConsoleColor()
	r := gin.Default()
	r.SetHTMLTemplate(template.Must(template.New("").ParseFS(fs, "templates/*.tmpl")))

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", map[string]any{
			"PageTitle": "titleXy",
			"Todos": []map[string]any{
				{
					"Title": "dumm",
					"Done":  false,
				},
				{
					"Title": "blöd",
					"Done":  true,
				},
			},
		})
	})

	// https://webfinger.net/

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
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
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	for {
		sig := <-signals
		switch sig {
		case syscall.SIGHUP:
			log.Println("reloading config...")
			// TODO
		case syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM:
			log.Println("shutting down server...")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err := srv.Shutdown(ctx)
			cancel()
			if err != nil {
				log.Fatal("server shutdown error: ", err)
			}
			os.Exit(0)
		default:
			log.Fatalf("unhandled signal: %s", sig.String())
		}
	}
}
