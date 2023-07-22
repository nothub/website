package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
)

var projectsAuthor []string
var projectsContrib []string

func initProjects(router *gin.Engine) (err error) {
	log.Println("loading projects")
	file, err := fs.ReadFile("data/projects.yaml")
	if err != nil {
		return err
	}

	data := make(map[string][]string)
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		return err
	}

	projectsAuthor = data["authored"]
	projectsContrib = data["contributed"]

	router.GET("/projects/*path", func(c *gin.Context) {
		path := c.Param("path")
		log.Printf("path=%q\n", path)
		// TODO
		c.AbortWithStatus(http.StatusTeapot)
	})

	return nil
}
