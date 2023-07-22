package main

import (
	"embed"
	_ "embed"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

//go:embed data/*
var dataFs embed.FS

func Init() {

	log.Println("generating projects")
	file, err := ioutil.ReadFile("data/projects.yaml")
	if err != nil {
		log.Fatalln(err.Error())
	}

	data := make(map[string][]string)
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		log.Fatalln(err.Error())
	}

	for _, project := range data["authored"] {
		log.Printf("authored: %s\n", project)
	}

	for _, project := range data["contributed"] {
		log.Printf("contributed: %s\n", project)
	}
}
