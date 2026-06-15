package main

import (
	"errors"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"net/url"
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
		log.Fatalln(err.Error())
	}

	var projects []Project
	if err := yaml.Unmarshal(data, &projects); err != nil {
		log.Fatalln(err.Error())
	}

	for _, project := range projects {
		for _, tag := range project.Tags {
			linkTag(tag, "Project: "+project.Title, project.Url)
		}
		for _, lang := range project.Langs {
			linkTag(lang, "Project: "+project.Title, project.Url)
		}
	}

	go fetchStars(&projects)

	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for range ticker.C {
			fetchStars(&projects)
		}
	}()

	mux.HandleFunc("GET /projects", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.ExecuteTemplate(w, "projects.gohtml", projects); err != nil {
			slog.Warn("template error", "err", err)
		}
	})

	return nil
}

func fetchStars(projects *[]Project) {
	log.Printf("fetching github stars for %v projects\n", len(*projects))

	for i, proj := range *projects {
		u, err := url.Parse(proj.Url)
		if err != nil {
			log.Printf("invalid project url %s: %s\n", proj.Url, err)
			continue
		}

		if u.Host != "github.com" {
			continue
		}

		var meta *github.RepoMeta
		backoff := []time.Duration{0, 5 * time.Second, 10 * time.Second, 20 * time.Second}
		for attempt, wait := range backoff {
			if wait > 0 {
				time.Sleep(wait)
			}
			meta, err = github.FetchRepoMeta(u.Path)
			if err == nil {
				break
			}
			var rle github.RateLimitError
			if errors.As(err, &rle) {
				cap := 90 * time.Minute
				dur := rle.RetryAfter
				if dur > cap {
					dur = cap
				}
				log.Printf("rate limited fetching %s; waiting %s\n", proj.Url, dur)
				time.Sleep(dur)
				meta, err = github.FetchRepoMeta(u.Path)
				break
			}
			log.Printf("attempt %d for %s: %s\n", attempt+1, proj.Url, err)
		}

		if meta != nil {
			(*projects)[i].Stars = meta.StargazersCount
		} else {
			slog.Warn("failed to fetch stars", "url", proj.Url)
		}

		if i < len(*projects)-1 {
			time.Sleep(30 * time.Second)
		}
	}
}
