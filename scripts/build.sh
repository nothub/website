#!/usr/bin/env sh

set -eu
cd "$(dirname "$(realpath "$0")")/.."

TAG="${TAG:-dev}"
CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o website .
go tool ocipack \
    -tag "ghcr.io/nothub/website:${TAG}" \
    website image.tar.gz
rm -f website
