#!/usr/bin/env sh

set -eu
cd "$(dirname "$(realpath "$0")")/.."

# shellcheck disable=SC2046
go test -v -race -vet=all $(go list ./... | grep -v '/assets$')
