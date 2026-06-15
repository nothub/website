package main

import (
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

// TODO: setCacheHeader → Cache-Control: public, max-age=604800, immutable
// TODO: setClacksHeader → middleware that sets X-Clacks-Overhead: GNU Terry Pratchett
