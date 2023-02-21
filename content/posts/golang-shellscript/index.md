---
title: "Golang Shebang Scripts"
description: "Polyglot Golang-Script with Shell Entrypoint"
date: "2023-01-03"
tags: ["golang"]
draft: true
---




Recently, I required some glue code.

single file uncompiled go src
bash header with polyglot trick
embedded go.mod for deps
cleanup go.mod
check for existing go.mod
formatter disabler

i made a really bad polyglot header for running go source files from shell scripts, check this out xD

bash and go do not share the same comment syntax
// is a comment in go but not in bash
// is simplified to / in a unix environment
we can abuse this fact and wrap shell commands in go comments

after working with go for some time, when writing shell scripts, i miss the clarity/simplicity/? of go...
how run go scripts?
solution -> polyglot

tldr

intellij formatter

external dependencies -> go.mod

create go.mod and go.sum on demand

```go
// @formatter:off 2>/dev/null
// 2>/dev/null || set -o errexit; set -o nounset
// 2>/dev/null || if [ -f "go.mod" ]; then echo "Unable to exec because a go module is present in the working directory!"; exit 1; fi
// 2>/dev/null || mod_path="hub.lol/foo/bar"; mod_gover="1.19"; mod_pkgs=('github.com/spf13/pflag v1.0.5' 'gopkg.in/yaml.v3 v3.0.1')
//usr/bin/env -S printf "module %s\n\ngo %s\n\nrequire (\n%s\n)" "${mod_path}" "${mod_gover}" "$(IFS=$'\n'; echo "${mod_pkgs[*]}")" > go.mod
//usr/bin/env -S go mod tidy; set +o errexit; go run "$0" "$@"; exit_status="$?"; rm -f go.mod go.sum; exit "${exit_status}"
// @formatter:on

package main

import (
    "fmt"
    "net/http"
    "sort"
    "time"

    "gopkg.in/yaml.v3"
    flag "github.com/spf13/pflag"
)

func main() {
    println("Hello, World!")
	// ...
}
```

not nice: this is not a shebang! whatever shell is invoking this script will also be used to interpret it!
this is intended to be interpreted by bash, not sh!

todo for sh compat: remove array usage

---

shebang.go
- [ ] valid shell
- [ ] valid bash
- [ ] valid go

simple.go
- [x] valid shell
- [x] valid bash
- [x] valid go
