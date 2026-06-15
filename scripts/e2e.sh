#!/usr/bin/env sh

set -eu
cd "$(dirname "$(realpath "$0")")/.."

CONTAINER="website-e2e"
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
podman run -d --rm -p "${PORT}:8080" --name "${CONTAINER}" "docker.io/n0thub/website:dev"

i=0
while test "${i}" -lt 10; do
    STATUS=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:${PORT}/teapot" 2> /dev/null || true)
    if test "${STATUS}" = "418"; then
        echo "e2e: healthcheck passed (418 Teapot)"
        exit 0
    fi
    i=$((i + 1))
    echo "e2e: attempt ${i}/10: got ${STATUS}, waiting..."
    sleep 1
done

echo "e2e: healthcheck failed after 10 attempts"
exit 1
