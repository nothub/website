package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
)

// TODO: fetch infos from github api (https://github.com/orgs/community/discussions/24350)

type project struct {
	name  string
	desc  string
	url   string
	langs []string
	tags  []string
	stars int
}

var langColors = make(map[string]string)

func init() {
	langColors["dockerfile"] = "#384D54"
	langColors["go"] = "#00ADD8"
	langColors["java"] = "#B07219"
	langColors["lua"] = "#000080"
	langColors["perl"] = "#0298C3"
	langColors["python"] = "#3572A5"
	langColors["shell"] = "#89E051"
}

var projectsAuthor []project
var projectsContrib []project

func initProjects(router *gin.Engine) (err error) {
	log.Println("loading projects")
	file, err := fs.ReadFile("data/projects.yaml")
	if err != nil {
		return err
	}

	data := make(map[string][]project)
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		return err
	}

	projectsAuthor = data["author"]
	projectsContrib = data["contrib"]

	router.GET("/projects/*path", func(c *gin.Context) {
		path := c.Param("path")
		log.Printf("path=%q\n", path)
		// TODO
		c.AbortWithStatus(http.StatusTeapot)
	})

	return nil
}
