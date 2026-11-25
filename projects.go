package main

import (
	"errors"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync/atomic"
	"time"

	"github.com/nothub/website/internal/github"
	"gopkg.in/yaml.v3"
)

type Project struct {
	Title string   `yaml:"title"`
	Url   string   `yaml:"url"`
	Desc  string   `yaml:"desc"`
	Role  string   `yaml:"role"`
	Tags  []string `yaml:"tags"`
	Langs []string `yaml:"langs"`
	Stars int      `yaml:"stars"`
}

func initProjects(mux *http.ServeMux, tmpl *template.Template) (err error) {
	log.Println("loading projects")

	data, err := fs.ReadFile("data/projects.yaml")
	if err != nil {
		return err
	}

	var projects []Project
	if err := yaml.Unmarshal(data, &projects); err != nil {
		return err
	}

	for _, project := range projects {
		for _, tag := range project.Tags {
			linkTag(tag, "Project: "+project.Title, project.Url)
		}
		for _, lang := range project.Langs {
			linkTag(lang, "Project: "+project.Title, project.Url)
		}
	}

	// Requests read the list while a refresh runs, so each refresh builds a
	// new slice and swaps it in whole instead of writing into the shared one.
	var current atomic.Pointer[[]Project]
	current.Store(&projects)
	refresh := func() {
		updated := fetchStars(*current.Load())
		current.Store(&updated)
	}

	go refresh()

	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for range ticker.C {
			refresh()
		}
	}()

	mux.HandleFunc("GET /projects", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.ExecuteTemplate(w, "projects.gohtml", *current.Load()); err != nil {
			slog.Warn("template error", "err", err)
		}
	})

	return nil
}

// githubUser's repos come from one paginated listing instead of a request
// per repo.
const githubUser = "nothub"

func fetchStars(projects []Project) []Project {
	log.Printf("fetching github stars for %v projects\n", len(projects))

	stars := make(map[string]int)
	repos, err := withRetry(githubUser+"'s repos", func() ([]github.RepoMeta, error) {
		return github.FetchUserRepos(githubUser)
	})
	if err != nil {
		slog.Warn("failed to list repos, fetching one by one", "user", githubUser, "err", err)
	}
	for _, r := range repos {
		stars[strings.ToLower(r.FullName)] = r.StargazersCount
	}

	updated := slices.Clone(projects)
	fetched := false
	for i, proj := range updated {
		u, err := url.Parse(proj.Url)
		if err != nil {
			log.Printf("invalid project url %s: %s\n", proj.Url, err)
			continue
		}

		if u.Host != "github.com" {
			continue
		}

		if n, ok := stars[strings.ToLower(strings.Trim(u.Path, "/"))]; ok {
			updated[i].Stars = n
			continue
		}

		if fetched {
			time.Sleep(30 * time.Second)
		}
		fetched = true

		meta, err := withRetry(proj.Url, func() (*github.RepoMeta, error) {
			return github.FetchRepoMeta(u.Path)
		})
		if err != nil {
			slog.Warn("failed to fetch stars", "url", proj.Url, "err", err)
			continue
		}
		updated[i].Stars = meta.StargazersCount
	}
	return updated
}

func withRetry[T any](what string, fetch func() (T, error)) (T, error) {
	var res T
	var err error
	backoff := []time.Duration{0, 5 * time.Second, 10 * time.Second, 20 * time.Second}
	for attempt, wait := range backoff {
		if wait > 0 {
			time.Sleep(wait)
		}
		res, err = fetch()
		if err == nil {
			return res, nil
		}
		var rle github.RateLimitError
		if errors.As(err, &rle) {
			dur := min(rle.RetryAfter, 90*time.Minute)
			log.Printf("rate limited fetching %s; waiting %s\n", what, dur)
			time.Sleep(dur)
			return fetch()
		}
		log.Printf("attempt %d for %s: %s\n", attempt+1, what, err)
	}
	return res, err
}
