package middleware

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestClacksHeaderSent(t *testing.T) {
	handler := Clacks(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("body"))
	}))

	seen := 0
	for range 1000 {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
		// Result() holds the headers as they were when the response was written.
		if w.Result().Header.Get("X-Clacks-Overhead") == "GNU Terry Pratchett" {
			seen++
		}
	}
	if seen == 0 {
		t.Error("X-Clacks-Overhead never sent in 1000 responses")
	}
}

func TestClacksConcurrent(t *testing.T) {
	handler := Clacks(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			for range 100 {
				handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
			}
		})
	}
	wg.Wait()
}
