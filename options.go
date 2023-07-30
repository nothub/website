package main

import (
	"flag"
	"os"
)

var optLoadDrafts bool
var optGithubToken string

func init() {
	flag.BoolVar(&optLoadDrafts, "drafts", false, "")
	flag.StringVar(&optGithubToken, "github-token", "", "")
	flag.Parse()

	if optGithubToken == "" {
		optGithubToken = os.Getenv("GITHUB_TOKEN")
	}
}
