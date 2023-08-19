package slogging

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var slogger = slog.New(slog.NewTextHandler(os.Stderr, nil))

var Gin = func(c *gin.Context) {
	start := time.Now()
	c.Next()
	latency := time.Now().Sub(start)
	attrs := []slog.Attr{
		slog.Int("status", c.Writer.Status()),
		slog.String("method", c.Request.Method),
		slog.String("path", c.Request.URL.Path),
		slog.String("ip", c.ClientIP()),
		slog.Duration("latency", latency),
		slog.String("ua", c.Request.UserAgent()),
	}
	if len(c.Errors) > 0 {
		attrs = append(attrs, slog.String("err", c.Errors.String()))
	}
	switch {
	case c.Writer.Status() >= http.StatusInternalServerError:
		slogger.Error("gin", attrs)
	default:
		slogger.Info("gin", attrs)
	}
}
