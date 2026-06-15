package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTrailingSlash(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := TrailingSlash(next)

	tests := []struct {
		path     string
		wantCode int
		wantLoc  string
	}{
		{"/", 200, ""},
		{"/about", 200, ""},
		{"/about/", 301, "/about"},
		{"/posts/foo/", 301, "/posts/foo"},
		{"/posts/foo/bar.png", 200, ""},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			r := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)
			if w.Code != tt.wantCode {
				t.Errorf("path %q: status = %d, want %d", tt.path, w.Code, tt.wantCode)
			}
			if tt.wantLoc != "" && w.Header().Get("Location") != tt.wantLoc {
				t.Errorf("path %q: Location = %q, want %q", tt.path, w.Header().Get("Location"), tt.wantLoc)
			}
		})
	}
}
