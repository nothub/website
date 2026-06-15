#!/usr/bin/env sh

set -eu
cd "$(dirname "$(realpath "$0")")/.."

TAG="${TAG:-e2e}"
CONTAINER="${CONTAINER:-website-e2e}"
PORT=$(python3 -c "import socket; s=socket.socket(); s.bind(('',0)); print(s.getsockname()[1]); s.close()")

# shellcheck disable=SC2317
cleanup() {
    podman stop "${CONTAINER}" 2> /dev/null || true
    rm -f image.tar.gz
}
trap cleanup EXIT

scripts/build.sh

podman load -i image.tar.gz

podman rm -f "${CONTAINER}" 2> /dev/null || true
podman run -d --rm -p "${PORT}:8080" --name "${CONTAINER}" "docker.io/n0thub/website:${TAG}"

i=0
while [ "${i}" -lt 10 ]; do
    STATUS=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:${PORT}/teapot" 2> /dev/null || true)
    if [ "${STATUS}" = "418" ]; then
        echo "e2e: healthcheck passed (418 Teapot)"
        exit 0
    fi
    i=$((i + 1))
    echo "e2e: attempt ${i}/10: got ${STATUS}, waiting..."
    sleep 1
done

echo "e2e: healthcheck failed after 10 attempts"
exit 1
