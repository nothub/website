package main

import (
	"log/slog"
	"os"
)

var slogger = slog.New(slog.NewJSONHandler(os.Stderr, nil))
