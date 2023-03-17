---
title: "Docker image hash-tags"
description: "Pin Docker images with a hashsum as tag"
date: "2023-03-17"
tags: ["docker", "TIL"]
draft: false
---

Today I learned:

Docker images can be referenced using the [content-digest](https://docs.docker.com/registry/spec/api/#content-digests)
hash in place of a tag.

```console
$ docker inspect alpine:latest | jq .[0].RepoDigests
[
  "alpine@sha256:69665d02cb32192e52e07644d76bc6f25abeb5410edc1c7a81a10ba3f0efb90a"
]

$ docker run -it alpine@sha256:69665d02cb32192e52e07644d76bc6f25abeb5410edc1c7a81a10ba3f0efb90a
/ # exit
```
