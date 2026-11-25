package github

import (
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestRateLimit(t *testing.T) {
	reset := strconv.FormatInt(time.Now().Add(time.Hour).Unix(), 10)
	tests := []struct {
		name    string
		status  int
		headers map[string]string
		want    bool
	}{
		{"429", 429, nil, true},
		{"403 primary limit", 403, map[string]string{"X-RateLimit-Remaining": "0", "X-RateLimit-Reset": reset}, true},
		{"403 secondary limit", 403, map[string]string{"Retry-After": "60"}, true},
		{"403 forbidden", 403, map[string]string{"X-RateLimit-Remaining": "59", "X-RateLimit-Reset": reset}, false},
		{"200", 200, map[string]string{"X-RateLimit-Remaining": "0"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := &http.Response{StatusCode: tt.status, Header: http.Header{}}
			for k, v := range tt.headers {
				res.Header.Set(k, v)
			}
			if _, got := rateLimit(res); got != tt.want {
				t.Errorf("rateLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}
