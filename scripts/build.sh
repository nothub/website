#!/usr/bin/env sh

set -eu
cd "$(dirname "$(realpath "$0")")/.."

TAG="${TAG:-dev}"
CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o website .
go tool ocipack \
    -tag "docker.io/n0thub/website:${TAG}" \
    -label "org.opencontainers.image.source=https://github.com/nothub/website" \
    -label "org.opencontainers.image.revision=$(git rev-parse HEAD)" \
    website image.tar.gz
rm -f website
