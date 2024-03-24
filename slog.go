package main

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var slogger = slog.New(slog.NewTextHandler(os.Stderr, nil))

var slogGin = func(ctx *gin.Context) {
	start := time.Now()
	ctx.Next()
	latency := time.Now().Sub(start)
	attrs := []any{
		slog.Int("status", ctx.Writer.Status()),
		slog.String("method", ctx.Request.Method),
		slog.String("path", ctx.Request.URL.Path),
		slog.String("query", ctx.Request.URL.RawQuery),
		slog.String("ip", ctx.ClientIP()),
		slog.Duration("latency", latency),
		slog.String("ua", ctx.Request.UserAgent()),
	}
	if len(ctx.Errors) > 0 {
		attrs = append(attrs, slog.String("err", ctx.Errors.String()))
	}
	switch {
	case ctx.Writer.Status() >= http.StatusInternalServerError:
		slogger.Error("gin", attrs...)
	default:
		slogger.Info("gin", attrs...)
	}
}
