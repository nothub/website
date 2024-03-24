#!/usr/bin/env sh

set -eu
cd "$(dirname "$(realpath "$0")")/.."

go test -vet=all .
