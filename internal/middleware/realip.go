package middleware

import (
	"log/slog"
	"net/http"
	"strings"
)

func RealIP(logger *slog.Logger, trust bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if trust {
				xff := r.Header.Get("X-Forwarded-For")
				if xff == "" {
					logger.Warn("missing X-Forwarded-For header", "remote", r.RemoteAddr)
				} else {
					// leftmost value is the real client IP (Envoy Gateway strips client-supplied values)
					r.RemoteAddr = strings.TrimSpace(strings.SplitN(xff, ",", 2)[0])
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
