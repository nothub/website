#!/usr/bin/env sh

set -eu
cd "$(dirname "$(realpath "$0")")/.."

CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o website .
go tool ocipack \
    -tag "docker.io/n0thub/website:dev" \
    -label "org.opencontainers.image.source=https://github.com/nothub/website" \
    -label "org.opencontainers.image.revision=$(git rev-parse HEAD)" \
    website image.tar.gz
rm -f website
