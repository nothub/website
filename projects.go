package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"log"
	"math/rand"
	"net/http"
	"net/url"
)

// TODO: fetch infos from github api (https://github.com/orgs/community/discussions/24350)

type project struct {
	Title string   `yaml:"title"`
	Url   string   `yaml:"url"`
	Desc  string   `yaml:"desc"`
	Role  string   `yaml:"role"`
	Tags  []string `yaml:"tags"`
	Langs []string `yaml:"langs"`
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
	return rand.Intn(3)
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