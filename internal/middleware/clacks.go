package middleware

import (
	"math/rand/v2"
	"net/http"
)

func Clacks(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rand.IntN(42) == 0 {
			w.Header().Set("X-Clacks-Overhead", "GNU Terry Pratchett")
		}
		next.ServeHTTP(w, r)
	})
}
