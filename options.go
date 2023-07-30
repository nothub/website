package main

import (
	"flag"
	"os"
)

var optLoadDrafts bool
var optGithubToken string = os.Getenv("GITHUB_TOKEN")

func init() {
	flag.BoolVar(&optLoadDrafts, "drafts", false, "load drafts")
	flag.StringVar(&optGithubToken, "github-token", "", "github token")
	flag.Parse()
}
