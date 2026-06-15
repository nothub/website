package middleware

import (
	"math/rand"
	"net/http"
)

func Clacks(rng *rand.Rand) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			if rng.Intn(42) == 0 {
				w.Header().Set("X-Clacks-Overhead", "GNU Terry Pratchett")
			}
		})
	}
}
