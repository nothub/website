package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"net/url"
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
	bytes, err := fs.ReadFile("data/projects.yaml")
	if err != nil {
		return err
	}

	var projects []Project
	err = yaml.Unmarshal(bytes, &projects)
	if err != nil {
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

	for i, proj := range projects {
		u, err := url.Parse(proj.Url)
		if err != nil {
			log.Fatalln(err.Error())
		}

		if u.Host != "github.com" {
			continue
		}

		meta, err := githubRepoMeta(u.Path)
		if err != nil {
			log.Printf("stargazer lookup for %s caused http status %s\n", proj.Url, err.Error())
			continue
		}

		projects[i].Stars = meta.StargazersCount
	}

	router.GET("/projects", func(c *gin.Context) {
		c.HTML(http.StatusOK, "projects.gohtml", projects)
	})

	return nil
}
