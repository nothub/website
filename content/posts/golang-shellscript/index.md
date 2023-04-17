---
title: "Golang Scripts"
description: "Polyglot Golang-Script with Shell Entrypoint"
date: "2023-01-03"
tags: ["go", "shell"]
draft: true
---

For a recent project, I needed a portable tool for the recurring task of reading YAML data from an HTTP endpoint.  
I was dreading the idea of implementing this in bash, so I started writing the code in Go.

Go however is a compiled language and I do not like to check binary files into version control.  
Instead, the tool should be implemented in a single text-file that contains common (uncompiled) go code.

## Shebang

The first idea coming to my mind was to wrap the [`go run`](https://pkg.go.dev/cmd/go#hdr-Compile_and_run_Go_program) command in a shebang header.

```go
#!/usr/bin/env -S go run
package main
func main() {
println("Hello, World!")
}
```

However, the syntax of Go does not allow this kind of header.

The prefix of a shebang header is `#` but that is, unlike in shell, not a valid token in Go.

```sh
$ ./shebang.go
shebang.go:1:1: illegal character U+0023 '#'
```

## Double-Path Trick

I found [this cool trick](https://gist.github.com/msoap/a9ee054f80a58b16867c) by [Serhii Mudryk](https://github.com/msoap).  
It utilizes the fact that, `//` is a valid Go token (a comment) and, [as defined](https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap03.html#tag_03_271) in the POSIX spec, all subpathes `//` are simplified to `/`.

This has the implication of being syntactically valid shell and Go code at the same time:

```go
//usr/bin/env -S go run "$0" "$@" ; exit
package main

func main() {
    println("Hello, World!")
}
```

To run the polyglot:

```go
$./simple.go
Hello, World!
```

Looking good so far 👌

## Handling Dependencies

Okay but the original problem also requires me to handle YAML, for that I want to use [gopkg.in/yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3).

In modern Go projects, dependencies are usually declared in the projects `go.mod` file and managed by the `go` cli tool.  
Ideally we want to leverage the same tooling for the script.

Instead of directly calling `go run` in the script-header, we can add some additional logic to create and use a `go.mod` file on the fly.

### Piggybacking

The way of calling `env` to execute stuff has a big issue; `env` is a program to "run (another) program in a modified environment".  
There is no way to use shell builtins or anything that is not the reference to another program.

We want however to use the full power of the shell language. To do so, we again start the line with our double path token but this time instead of invoking `//usr/bin/env`, we are going to invoke just `//`.

`//` has the correct syntax for a shell command but semantically it does not make much sense (and returns exit code 126).  
The `sh` shell will tell us that we do not have the correct permissions (`//: Permission denied`).  
The *bash* shell will tell us that: `//: Is a directory`

However, we do not care about this because, all we want is a command to piggyback on.  
We can just use the *separation operator* `;` to append another command:

```sh
$ //; echo "test"
sh: 1: //: Permission denied
test
```

To get rid of the permission error, we can redirect the first commands stderr output to `/dev/null`:

```sh
$ // 2>/dev/null; echo "test"
test
```

## Generate go.mod

With this method, we can finally build our `go.mod` file!

We are going to use some variables to comfortably hold some infos:

- `mod_path` will hold the script's own module path
- `mod_gover` Go version the script will target
- `mod_pkgs` will hold a comma separated list of dependencies

```sh
// 2>/dev/null; mod_path="hub.lol/foo/bar"
// 2>/dev/null; mod_gover="1.19"
// 2>/dev/null; mod_pkgs="github.com/spf13/pflag v1.0.5,gopkg.in/yaml.v3 v3.0.1"
```

To generate the `go.mod` file from this information, we can use `tr` and `printf`:

```sh
//usr/bin/env -S printf "module %s\n\ngo %s\n\nrequire (\n%s\n)" "${mod_path}" "${mod_gover}" "$(echo "${mod_pkgs}" | tr "," "\n")" > go.mod
```

## Autoformat

## Cleanup

## IntelliJ Formatter

## Drawbacks



INFO BOX
Sadly, by specification, there are some exceptions to this behaviour 😔  
A list of platforms where this trick does not work is found [here](https://unix.stackexchange.com/questions/256497/on-what-systems-is-foo-bar-different-from-foo-bar) (most notably: **Cygwin** and **Bazel**).

embedded go.mod for deps
cleanup go.mod
check for existing go.mod
formatter disabler

intellij formatter

external dependencies -> go.mod

create go.mod and go.sum on demand

not nice: this is not a shebang! whatever shell is invoking this script will also be used to interpret it!
this is intended to be interpreted by bash, not sh!

todo for sh compat: remove array usage

---

lint header with:

`# shellcheck disable=SC1127,SC2317`

---

shebang.go

- [ ] valid shell
- [ ] valid bash
- [ ] valid go

simple.go

- [x] valid shell
- [x] valid bash
- [x] valid go

👋 🌳 🪐

---

drawback: self extracting module needs write permissions in workdir
