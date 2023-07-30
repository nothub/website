package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

type Read struct {
	Title string   `yaml:"title"`
	Url   string   `yaml:"url"`
	Tags  []string `yaml:"tags"`
}

func initReads(router *gin.Engine) (err error) {
	log.Println("organizing reads")
	bytes, err := fs.ReadFile("data/reads.yaml")
	if err != nil {
		return err
	}

	var reads []Read
	err = yaml.Unmarshal(bytes, &reads)
	if err != nil {
		return err
	}

	for _, read := range reads {
		for _, tag := range read.Tags {
			linkTag(tag, "Read: "+read.Title, read.Url)
		}
	}

	router.GET("/reads", func(c *gin.Context) {
		c.HTML(http.StatusOK, "reads.gohtml", reads)
	})

	return nil
}
