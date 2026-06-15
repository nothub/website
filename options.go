package main

import "flag"

var optLoadDrafts bool
var optTrustProxy bool

func initFlags() {
	flag.BoolVar(&optLoadDrafts, "drafts", false, "")
	flag.BoolVar(&optTrustProxy, "trust-proxy", false, "trust X-Forwarded-For header set by upstream proxy")
	flag.Parse()
}
