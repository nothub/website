---
title: "Docker Compose inline Dockerfile"
description: "Include Dockerfiles inline as attribute in Compose files"
date: "2024-03-24"
tags: [ "docker" ]
draft: false
---

The [dockerfile_inline](https://docs.docker.com/compose/compose-file/build/#dockerfile_inline) attribute got introduced
to Docker Compose in [version 2.17.0](https://github.com/docker/compose/releases/tag/v2.17.0), about a year ago.

I just found out and I really like it!

Example `docker-compose.yaml` file:

```yaml
services:
  server:
    container_name: "hi"
    build:
      dockerfile_inline: |
        FROM debian:12-slim as builder
        RUN echo "f0VMRgEBAfCfkYsKAAAAAAIAAwABAAAAgIAECDQAAAAAuAQAAADNgOtYIAACACgABQAEAAEAAAAAAAAAAIAECACABAiiAAAAogAAAAUAAAAAEAAAAQAAAKQAAACkkAQIpJAECAkAAAAJAAAAugkAAAC5B5AECLsBAAAA66QAAADr6rsAAAAAuAEAAADNgA==" \
          | base64 -d > /bin/hi && chmod +x /bin/hi
        FROM scratch
        COPY --from=builder /bin/hi /bin/hi
    command: [ /bin/hi ]
```
