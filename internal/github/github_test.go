package github

import (
	"net/http"
	"net/http/httptest"
	"slices"
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

func TestFetchUserRepos(t *testing.T) {
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/someone/repos" {
			t.Errorf("unexpected path %q", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		switch r.URL.Query().Get("page") {
		case "":
			w.Header().Set("Link", `<`+srv.URL+`/users/someone/repos?per_page=100&page=2>; rel="next", <`+srv.URL+`/users/someone/repos?per_page=100&page=2>; rel="last"`)
			w.Write([]byte(`[{"full_name":"someone/a","stargazers_count":1},{"full_name":"someone/b","stargazers_count":2}]`))
		case "2":
			w.Header().Set("Link", `<`+srv.URL+`/users/someone/repos?per_page=100&page=1>; rel="prev", <`+srv.URL+`/users/someone/repos?per_page=100&page=1>; rel="first"`)
			w.Write([]byte(`[{"full_name":"someone/c","stargazers_count":3}]`))
		default:
			t.Errorf("unexpected page %q", r.URL.Query().Get("page"))
		}
	}))
	defer srv.Close()

	orig := apiURL
	apiURL = srv.URL
	defer func() { apiURL = orig }()

	repos, err := FetchUserRepos("someone")
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, r := range repos {
		got = append(got, r.FullName+"="+strconv.Itoa(r.StargazersCount))
	}
	want := []string{"someone/a=1", "someone/b=2", "someone/c=3"}
	if !slices.Equal(got, want) {
		t.Errorf("repos = %v, want %v", got, want)
	}
}
