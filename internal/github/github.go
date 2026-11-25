package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var client = &http.Client{
	Timeout: time.Second * 10,
}

var apiURL = "https://api.github.com"

type RateLimitError struct {
	RetryAfter time.Duration
}

func (e RateLimitError) Error() string {
	return fmt.Sprintf("github rate limit exceeded; retry after %s", e.RetryAfter)
}

type RepoMeta struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
	Owner    struct {
		Login string `json:"login"`
		Type  string `json:"type"`
	} `json:"owner"`
	Description     string    `json:"description"`
	Fork            bool      `json:"fork"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	PushedAt        time.Time `json:"pushed_at"`
	Homepage        string    `json:"homepage"`
	Size            int       `json:"size"`
	StargazersCount int       `json:"stargazers_count"`
	WatchersCount   int       `json:"watchers_count"`
	Language        string    `json:"language"`
	ForksCount      int       `json:"forks_count"`
	Archived        bool      `json:"archived"`
	License         struct {
		Key    string `json:"key"`
		Name   string `json:"name"`
		SpdxId string `json:"spdx_id"`
		Url    string `json:"url"`
		NodeId string `json:"node_id"`
	} `json:"license"`
	AllowForking  bool     `json:"allow_forking"`
	IsTemplate    bool     `json:"is_template"`
	Topics        []string `json:"topics"`
	Forks         int      `json:"forks"`
	Watchers      int      `json:"watchers"`
	DefaultBranch string   `json:"default_branch"`
	Organization  struct {
		Login string `json:"login"`
		Type  string `json:"type"`
	} `json:"organization"`
}

func FetchRepoMeta(repo string) (*RepoMeta, error) {
	repo = strings.TrimSuffix(repo, ".git")
	split := strings.Split(repo, "/")
	if len(split) < 2 {
		return nil, errors.New("malformed repo string")
	}

	path := fmt.Sprintf("%s/%s", split[len(split)-2], split[len(split)-1])

	var meta RepoMeta
	if _, err := getJSON(apiURL+"/repos/"+path, &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}

// FetchUserRepos returns all public repos owned by user, following the
// pagination links. The API returns at most 100 repos per page.
func FetchUserRepos(user string) ([]RepoMeta, error) {
	var all []RepoMeta
	next := apiURL + "/users/" + url.PathEscape(user) + "/repos?per_page=100&type=owner"
	for next != "" {
		var page []RepoMeta
		var err error
		next, err = getJSON(next, &page)
		if err != nil {
			return nil, err
		}
		all = append(all, page...)
	}
	return all, nil
}

// getJSON decodes the response body of a GET to u into v and returns the
// URL of the next page from the Link header, or "" on the last page.
func getJSON(u string, v any) (next string, err error) {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "hub.lol")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Close = true

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if rle, ok := rateLimit(res); ok {
		return "", rle
	}

	if res.StatusCode < 200 || res.StatusCode >= 400 {
		return "", errors.New("http status " + res.Status)
	}

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(buf, v); err != nil {
		return "", err
	}

	return nextLink(res.Header.Get("Link")), nil
}

// nextLink extracts the rel="next" target from a Link header such as
// `<https://...&page=2>; rel="next", <https://...&page=3>; rel="last"`.
// https://docs.github.com/en/rest/using-the-rest-api/using-pagination-in-the-rest-api
func nextLink(header string) string {
	for _, part := range strings.Split(header, ",") {
		target, params, ok := strings.Cut(part, ";")
		if ok && strings.Contains(params, `rel="next"`) {
			return strings.Trim(strings.TrimSpace(target), "<>")
		}
	}
	return ""
}

// rateLimit reports whether res is a rate-limit response. GitHub also sends 403
// for plain permission errors, so a 403 only counts when it carries a
// Retry-After or an exhausted X-RateLimit-Remaining.
// https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api
func rateLimit(res *http.Response) (RateLimitError, bool) {
	retryAfter := res.Header.Get("Retry-After")
	switch {
	case res.StatusCode == 429:
	case res.StatusCode == 403 && (retryAfter != "" || res.Header.Get("X-RateLimit-Remaining") == "0"):
	default:
		return RateLimitError{}, false
	}

	var dur time.Duration
	if retryAfter != "" {
		secs, err := strconv.Atoi(retryAfter)
		if err == nil {
			dur = time.Duration(secs) * time.Second
		}
	} else if v := res.Header.Get("X-RateLimit-Reset"); v != "" {
		unix, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			dur = time.Until(time.Unix(unix, 0))
		}
	}
	return RateLimitError{RetryAfter: dur}, true
}
