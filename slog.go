package main

import (
	"log/slog"
	"os"
)

var slogger = slog.New(slog.NewJSONHandler(os.Stderr, nil))

// TODO: slog middleware → log status, method, path, query, ip, latency, ua per request
