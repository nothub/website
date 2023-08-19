package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"net/url"
	"time"
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

func initProjects(router *gin.Engine) (err error) {
	log.Println("loading projects")

	// load projects from data
	bytes, err := fs.ReadFile("data/projects.yaml")
	if err != nil {
		return err
	}
	var projects []Project
	err = yaml.Unmarshal(bytes, &projects)
	if err != nil {
		return err
	}

	// register tags
	for _, project := range projects {
		for _, tag := range project.Tags {
			linkTag(tag, "Project: "+project.Title, project.Url)
		}
		for _, lang := range project.Langs {
			linkTag(lang, "Project: "+project.Title, project.Url)
		}
	}

	// fetch stars on startup
	go fetchStars(&projects)

	// fetch stars again every day
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		// when the server is shut down, this go routine will not be shutdown
		// gracefully (it will just be killed), so do not do fancy stuff in here!
		for {
			select {
			case <-ticker.C:
				fetchStars(&projects)
			}
		}
	}()

	router.GET("/projects", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "projects.gohtml", projects)
	})

	return nil
}

func fetchStars(projects *[]Project) {
	log.Printf("fetching github stars for %v projects\n", len(*projects))

	for i, proj := range *projects {

		u, err := url.Parse(proj.Url)
		if err != nil {
			log.Fatalf("invalid project url %s caused %s\n", proj.Url, err.Error())
		}

		if u.Host != "github.com" {
			continue
		}

		meta, err := githubRepoMeta(u.Path)
		if err != nil {
			log.Printf("stargazer lookup for %s caused %s\n", proj.Url, err.Error())
			continue
		}

		(*projects)[i].Stars = meta.StargazersCount
	}
}
