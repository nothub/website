package middleware

import (
	"log"
	"net/http"
)

func RealIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Real-IP") != "" {
			r.RemoteAddr = r.Header.Get("X-Real-IP")
			r.Header.Del("X-Real-IP")
		} else {
			log.Printf("Request from %s has no X-Real-IP header!\n", r.RemoteAddr)
		}
		next.ServeHTTP(w, r)
	})
}
