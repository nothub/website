#!/usr/bin/env sh

set -eu
cd "$(dirname "$(realpath "$0")")/.."

GOOS=linux GOARCH=amd64 go build -tags osusergo,netgo,timetzdata -ldflags="-s -w" -o website
