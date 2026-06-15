package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggerSensitiveHeaders(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	handler := Logger(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer secret-token")
	r.Header.Set("Cookie", "session=abc123")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	logged := buf.String()
	if strings.Contains(logged, "secret-token") {
		t.Error("Logger logged Authorization header value")
	}
	if strings.Contains(logged, "abc123") {
		t.Error("Logger logged Cookie header value")
	}
}
