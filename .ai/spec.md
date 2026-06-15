# Spec: scripts/e2e.sh

## Objective

A single shell script that exercises the full production path end-to-end:
build binary → pack OCI image → import into podman → run container →
verify the `/teapot` route returns HTTP 418 → clean up on success.

Target user: developer running locally or CI pipeline verifying a release.

---

## Commands

```sh
# Run with default tag "e2e"
bash scripts/e2e.sh

# Override tag
TAG=v1.2.3 bash scripts/e2e.sh
```

---

## Steps (acceptance criteria per step)

### 1. Build
- Call `scripts/build.sh` (inherits `TAG` env var).
- Produces `image.tar.gz` in repo root.

### 2. Load into podman
- `podman load -i image.tar.gz`
- Image is now available as `docker.io/n0thub/website:${TAG}`.

### 3. Run container
- Pick a random free host port: bind a socket to port 0 and read back the assigned port,
  then close the socket before podman uses it. Use Python one-liner for portability:
  `PORT=$(python3 -c "import socket; s=socket.socket(); s.bind(('',0)); print(s.getsockname()[1]); s.close()")`
- Remove any stale container with the same name first:
  `podman rm -f "${CONTAINER}" 2>/dev/null || true`
- `podman run -d --rm -p "${PORT}:8080" --name "${CONTAINER}" "docker.io/n0thub/website:${TAG}"`
- `CONTAINER` defaults to `website-e2e`.

### 4. Healthcheck loop
- Poll `http://localhost:${PORT}/teapot` with `curl -s -o /dev/null -w "%{http_code}"`.
- Retry up to 10 times with 1 s sleep between attempts.
- Accept HTTP 418 as passing.
- Any other status (including curl failure) counts as a failed attempt.
- If all retries exhausted → print failure message, exit 1.

### 5. Cleanup (always)
- Always: `podman stop "${CONTAINER}" 2>/dev/null || true` (container auto-removes due to `--rm`),
  then `rm -f image.tar.gz`.
- Use a `trap cleanup EXIT` — runs unconditionally on success and failure.

---

## Project structure

Only one file is added:

```
scripts/e2e.sh
```

No other files are touched.

---

## Code style

- Shebang: `#!/usr/bin/env sh`
- `set -eu` (no `-o pipefail` — POSIX sh, matching existing scripts)
- 4-space indent
- `cd "$(dirname "$(realpath "$0")")/.."` at top (matching build.sh / test.sh)
- `TAG` defaults to `e2e` (not `dev`) to distinguish from interactive builds
- `CONTAINER` defaults to `website-e2e`
- `PORT` is selected randomly at runtime via Python's `socket` module

---

## Testing strategy

The script itself is the test. Acceptance is verified by running it once:

```sh
bash scripts/e2e.sh
```

Must exit 0 and print `e2e: healthcheck passed (418 Teapot)`.

Also run `shellcheck scripts/e2e.sh` (no warnings without inline disables).

---

## Boundaries

| Always do | Ask first | Never do |
|---|---|---|
| Clean up on success | — | Force-push or publish images |
| `set -eu` | — | Add dependencies beyond curl + podman |
| Exit non-zero on any failure | — | Use bash-specific syntax |
| Print attempt progress to stdout | — | Sleep before first curl attempt |
