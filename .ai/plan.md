# Plan: scripts/e2e.sh

## Task 1 — Write scripts/e2e.sh

**Files touched:** `scripts/e2e.sh`

**Steps:**
1. RED: `shellcheck scripts/e2e.sh` fails (file does not exist)
2. Write the script per spec
3. GREEN: `shellcheck scripts/e2e.sh` exits 0 with no warnings
4. Commit: `ci: add e2e end-to-end test script`

**Acceptance criteria:**
- `#!/usr/bin/env sh` + `set -eu`
- `cd "$(dirname "$(realpath "$0")")/.."` at top
- `TAG` defaults to `e2e`, `CONTAINER` defaults to `website-e2e`
- `PORT` chosen via Python socket one-liner
- Calls `scripts/build.sh`
- `podman load -i image.tar.gz`
- `podman rm -f "${CONTAINER}" 2>/dev/null || true` before run
- `podman run -d --rm -p "${PORT}:8080" --name "${CONTAINER}" ...`
- `trap cleanup EXIT` — always stops container + removes `image.tar.gz`
- Healthcheck polls `http://localhost:${PORT}/teapot`, up to 10 retries × 1 s
- Exits 0 on 418, exits 1 on exhaustion
- No sleep before first curl attempt
- `shellcheck` clean — no warnings without inline disables

**Verification:** `shellcheck scripts/e2e.sh` exits 0
