package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
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

var repoCache = make(map[string]RepoMeta)

func githubRepoMeta(repo string) (*RepoMeta, error) {
	repo = strings.TrimSuffix(repo, ".git")
	split := strings.Split(repo, "/")
	if len(split) < 2 {
		return nil, errors.New("malformed repo string")
	}

	path := fmt.Sprintf("%s/%s", split[len(split)-2], split[len(split)-1])

	if meta, ok := repoCache[path]; ok {
		return &meta, nil
	}

	u := "https://api.github.com/repos/" + path
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "hub.lol")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", optGithubToken))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Close = true

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

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
	repoCache[path] = meta

	return &meta, nil
}
