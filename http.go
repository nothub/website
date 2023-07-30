package main

import (
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

func setDefaultHeader(req *http.Request) {
	req.Header.Set("User-Agent", "github.com/nothub/website")
}
