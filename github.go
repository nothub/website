package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// TODO: use some ready to go api lib?

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

func githubRepoMeta(repo string) (*RepoMeta, error) {
	repo = strings.TrimSuffix(repo, ".git")
	split := strings.Split(repo, "/")
	if len(split) < 2 {
		return nil, errors.New("malformed repo string")
	}

	path := fmt.Sprintf("%s/%s", split[len(split)-2], split[len(split)-1])

	u := "https://api.github.com/repos/" + path
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "hub.lol")
	if optGithubToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", optGithubToken))
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Close = true

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// maybe rate-limited
	// https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api?apiVersion=2022-11-28
	if res.StatusCode == 403 || res.StatusCode == 429 {

		if res.Header.Get("retry-after") != "" {
			log.Printf("github api rate limit exceeded; header detected: retry-after = %s\n", res.Header.Get("retry-after"))

			v, err := strconv.Atoi(res.Header.Get("retry-after"))
			if err != nil {
				return nil, fmt.Errorf("error awaiting rate-limit: %w", err)
			}

			dur := time.Duration(v) * time.Second
			log.Printf("coplying with retry-after header; sleeping for %s\n", dur)
			time.Sleep(dur)
			return nil, errors.New("rate-limit awaited")
		}

		if res.Header.Get("x-ratelimit-remaining") == "0" {
			log.Printf("github api rate limit exceeded; header detected: x-ratelimit-remaining = %s\n", res.Header.Get("x-ratelimit-remaining"))

			v, err := strconv.ParseInt(res.Header.Get("x-ratelimit-reset"), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error awaiting rate-limit: %w", err)
			}

			t := time.Unix(v, 0)
			dur := t.Sub(time.Now())

			log.Printf("coplying with x-ratelimit-* headers; sleeping for %s\n", dur)
			time.Sleep(dur)
			return nil, errors.New("rate-limit awaited")
		}

	}

	if res.StatusCode < 200 || res.StatusCode >= 400 {
		return nil, errors.New("http status " + res.Status)
	}

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var meta RepoMeta
	err = json.Unmarshal(buf, &meta)
	if err != nil {
		return nil, err
	}

	log.Printf("project %s has %v stargazers\n", meta.FullName, meta.StargazersCount)

	return &meta, nil
}
