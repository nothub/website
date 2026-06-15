#!/usr/bin/env sh

set -eu
cd "$(dirname "$(realpath "$0")")/.."

# shellcheck disable=SC2046
go test -v -race -vet=all -count=1 $(go list ./... | grep -v '/posts/')

./scripts/e2e.sh
