package main

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"net/url"
)

// TODO: fetch infos from github api (https://github.com/orgs/community/discussions/24350)

type role string

const (
	author  role = "author"
	contrib role = "contrib"
)

type project struct {
	Name  string   `yaml:"name"`
	Desc  string   `yaml:"desc"`
	Tags  []string `yaml:"tags"`
	Role  role     `yaml:"role"`
	Url   string   `yaml:"url"`
	Langs []lang   `yaml:"langs"`
}

func (pr project) Colors() (colors []string) {
	return lo.Map(pr.Langs, func(lang lang, index int) string {
		return langColors[lang]
	})
}

func (pr project) Stars() int {
	u, err := url.Parse(pr.Url)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if u.Host != "github.com" {
		return 0
	}
	// TODO: fetch star count from api and cache for some time
	return 0
}

func initProjects(router *gin.Engine) (err error) {
	log.Println("loading projects")
	bytes, err := fs.ReadFile("data/projects.yaml")
	if err != nil {
		return err
	}

	var projects []project
	err = yaml.Unmarshal(bytes, &projects)
	if err != nil {
		return err
	}

	log.Printf("%++q\n", projects)

	router.GET("/projects/*path", func(c *gin.Context) {
		path := c.Param("path")
		log.Printf("path=%q\n", path)
		// TODO
		c.HTML(http.StatusOK, "projects.tmpl", projects)
	})

	return nil
}
