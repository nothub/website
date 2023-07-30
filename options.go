package main

import "flag"

var optLoadDrafts bool

func flags() {
	flag.BoolVar(&optLoadDrafts, "drafts", false, "load drafts")
	flag.Parse()
}
