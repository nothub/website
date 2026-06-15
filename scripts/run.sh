#!/usr/bin/env sh

set -eu
cd "$(dirname "$(realpath "$0")")/.."

CONTAINER="website-run"
PORT=$(python3 -c "import socket; s=socket.socket(); s.bind(('',0)); print(s.getsockname()[1]); s.close()")

# shellcheck disable=SC2317
cleanup() {
    podman stop -i "${CONTAINER}" 2> /dev/null
    podman rm -f "${CONTAINER}" 2> /dev/null
    rm -f image.tar.gz
}
trap cleanup EXIT

scripts/build.sh

podman stop -i "${CONTAINER}" 2> /dev/null
podman rm -f "${CONTAINER}" 2> /dev/null
podman load -i image.tar.gz
printf "Visit: https://127.0.0.1:%s\n" "$PORT"
podman run -it --rm -p "${PORT}:8080" --name "${CONTAINER}" "docker.io/n0thub/website:dev"
