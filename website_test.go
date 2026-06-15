package main

import (
	"encoding/xml"
	"html/template"
	"log/slog"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func testHandler(t *testing.T) http.Handler {
	t.Helper()
	tmpl := template.Must(template.New("").ParseFS(fs, "templates/*.gohtml"))
	rng := rand.New(rand.NewSource(0))
	return buildHandler(tmpl, slog.Default(), rng, false)
}

func TestRootRedirect(t *testing.T) {
	h := testHandler(t)
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusMovedPermanently {
		t.Errorf("GET / status = %d, want 301", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "/about" {
		t.Errorf("GET / Location = %q, want /about", loc)
	}
}

func TestShellRedirect(t *testing.T) {
	h := testHandler(t)
	r := httptest.NewRequest("GET", "/shell", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusMovedPermanently {
		t.Errorf("GET /shell status = %d, want 301", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "/jslinux" {
		t.Errorf("GET /shell Location = %q, want /jslinux", loc)
	}
}

func TestPostsRSSRedirect(t *testing.T) {
	h := testHandler(t)
	r := httptest.NewRequest("GET", "/posts/rss.xml", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusMovedPermanently {
		t.Errorf("GET /posts/rss.xml status = %d, want 301", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "/rss.xml" {
		t.Errorf("GET /posts/rss.xml Location = %q, want /rss.xml", loc)
	}
}

func TestNotFound(t *testing.T) {
	h := testHandler(t)
	r := httptest.NewRequest("GET", "/this-does-not-exist", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusNotFound {
		t.Errorf("GET /unknown status = %d, want 404", w.Code)
	}
}

func TestStaticFiles(t *testing.T) {
	h := testHandler(t)
	paths := []string{"/static/style.css", "/static/normalize.css", "/static/robots.txt"}
	for _, path := range paths {
		r := httptest.NewRequest("GET", path, nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Errorf("GET %s status = %d, want 200", path, w.Code)
		}
	}
}

func TestTeapot(t *testing.T) {
	h := testHandler(t)
	r := httptest.NewRequest("GET", "/teapot", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusTeapot {
		t.Errorf("GET /teapot status = %d, want 418", w.Code)
	}
}

func TestRSSFeed(t *testing.T) {
	h := testHandler(t)
	r := httptest.NewRequest("GET", "/rss.xml", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /rss.xml status = %d, want 200", w.Code)
	}

	body := w.Body.String()

	if !strings.Contains(body, "Aaron Swartz") {
		t.Error("RSS feed missing Aaron Swartz comment")
	}

	var rss struct {
		XMLName xml.Name `xml:"rss"`
		Channel struct {
			Items []struct{} `xml:"item"`
		} `xml:"channel"`
	}
	if err := xml.Unmarshal([]byte(body), &rss); err != nil {
		t.Fatalf("RSS feed is not valid XML: %v", err)
	}
	if len(rss.Channel.Items) == 0 {
		t.Error("RSS feed has no items")
	}
}

func TestMarkdownContentNegotiation(t *testing.T) {
	h := testHandler(t)
	r := httptest.NewRequest("GET", "/posts/guestfish", nil)
	r.Header.Set("Accept", "text/markdown")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	ct := w.Header().Get("Content-Type")
	if ct != "text/markdown; charset=utf-8" {
		t.Errorf("Content-Type = %q, want text/markdown; charset=utf-8", ct)
	}
	if body := w.Body.String(); len(body) == 0 {
		t.Error("body is empty")
	}
}

func TestPostsSortedNewestFirst(t *testing.T) {
	h := testHandler(t)
	r := httptest.NewRequest("GET", "/rss.xml", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	body := w.Body.String()
	var dates []time.Time
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "<pubDate>") {
			s := strings.TrimPrefix(line, "<pubDate>")
			s = strings.TrimSuffix(s, "</pubDate>")
			d, err := time.Parse(time.RFC1123Z, s)
			if err == nil {
				dates = append(dates, d)
			}
		}
	}
	for i := 1; i < len(dates); i++ {
		if dates[i].After(dates[i-1]) {
			t.Errorf("RSS items not sorted newest-first: item %d (%v) is after item %d (%v)", i, dates[i], i-1, dates[i-1])
		}
	}
}
