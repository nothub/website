# Rewrite Spec — hub.lol website

## Objective

Finish the stdlib `net/http` migration (gin is already removed), restore all routes and
middleware, and deliver the polish items from the TODO list. No new pages; same URL surface.
Target: a clean, minimal Go web app that matches the `go-web-project` skill description.

Deployment: OCI container (`scratch`-based), currently on a home server behind a Cloudflare
tunnel, moving to a k8s cluster on Hetzner. The image is built with
[ocipack](https://codeberg.org/fhuebner/ocipack) — no Dockerfile or Docker daemon required.

---

## Routes

| Method | Path | Behaviour |
|--------|------|-----------|
| `*` | `/*` | 404 via `writeError` (catch-all; registered as `"/"` with no method) |
| GET | `/` | 301 → `/about` (registered as `"GET /"`, takes precedence over catch-all) |
| GET | `/shell` | 301 → `/jslinux` |
| GET | `/about` | render `about.gohtml` |
| GET | `/posts` | render `posts.gohtml`; RSS via content negotiation |
| GET | `/posts/{slug}` | render `post.gohtml`; markdown via content negotiation |
| GET | `/posts/{slug}/{file}` | serve per-post asset file |
| GET | `/posts/rss.xml` | 301 → `/rss.xml` (registered as literal `"GET /posts/rss.xml"`, takes precedence over `"GET /posts/{slug}"` by Go 1.22+ specificity rules) |
| GET | `/rss.xml` | RSS 2.0 feed |
| GET | `/projects` | render `projects.gohtml` |
| GET | `/tags` | render `tags.gohtml` (all tags) |
| GET | `/tags/{tag}` | render `tags.gohtml` (single tag) |
| GET | `/version` | plain text version string |
| GET | `/assets/{file}` | embedded FS, `Cache-Control: public, max-age=604800, immutable` |
| GET | `/static/{file}` | embedded FS, same cache header |
| GET | `/jslinux/{file}` | embedded FS, same cache header |
| GET | `/robots.txt` | serve `static/robots.txt` |
| GET | `/sitemap.xml` | serve `static/sitemap.xml` |
| GET | `/teapot` | 418, cache header, 🫖 |

---

## Content negotiation

Use `github.com/elnormous/contenttype` to parse the `Accept` header.

**Endpoints that negotiate:**

- `GET /posts` — returns `text/html` (default), or the RSS feed (same as `/rss.xml`) when
  `Accept: application/rss+xml` or `application/atom+xml`.
- `GET /posts/{slug}` — returns `text/html` (default), or raw `index.md` source (no
  transformation) when `Accept: text/markdown`.

---

## Per-post asset directories

Posts move from a flat `posts/*.md` layout to per-post directories:

```
posts/
  my-post/
    index.md          ← the post itself
    diagram.png       ← any asset for this post
  other-post/
    index.md
```

- URL slug is the directory name (unchanged from the old filename stem).
- Per-post assets served at `/posts/{slug}/{file}`.
- Slugs that do not have an `index.md` are skipped with a log warning.
- The embed glob changes from `posts/*.md` to `posts/**`.

---

## Middleware stack

Applied in request-arrival order (first to receive the request → last before the handler):

1. **Logger** (`internal/middleware.Logger`) — outermost so it captures the final response
   status from everything below, including 500s from Recovery and 301s from TrailingSlash.
   Reads `r.RemoteAddr` *after* calling `next` (not before), so it gets the real IP that
   RealIP has already set. Structured slog JSON: status, method, path, query, remote IP,
   latency, user-agent. Level Info for <500, Error for ≥500. Must NOT log Authorization,
   Cookie, or other sensitive headers. Receives a `*slog.Logger` as a constructor parameter.
2. **RealIP** (`internal/middleware.RealIP`) — parses `X-Forwarded-For`, takes the leftmost
   value, sets `r.RemoteAddr`. Logs a warning if the header is absent. Envoy Gateway appends
   the downstream address to `X-Forwarded-For` and strips any client-supplied values
   (`xff_num_trusted_hops`), so the leftmost entry is the real client IP.
   **Security assumption**: the middleware trusts this header unconditionally. This is safe
   because a Kubernetes `NetworkPolicy` restricts inbound traffic to the pod to Envoy Gateway
   pods only — no client can reach the pod directly to forge the header. If the NetworkPolicy
   is ever removed, this trust is broken.
3. **Recovery** (`internal/middleware.Recovery`) — inside Logger so that Logger sees the 500
   it writes, rather than an unlogged panic. Catches panics from everything further in,
   logs them with slog, returns 500. Receives a `*slog.Logger` as a constructor parameter.
   Note: if response headers were already flushed when the panic occurs, the status code
   cannot be changed; the panic is still logged.
4. **TrailingSlash** (`internal/middleware.TrailingSlash`) — inside Recovery so redirect
   logic bugs do not escape as unhandled panics.
5. **Clacks** (`internal/middleware.Clacks`) — innermost, wraps the `http.ResponseWriter` to
   probabilistically inject `X-Clacks-Overhead: GNU Terry Pratchett` on ~1 in 42 requests.
   Uses its own `*rand.Rand` seeded at construction time; no dependency on `main` package.
   No security relevance.

---

## Project structure

```
website/
├── assets/                  # tool scripts (formatter.go, module.go, etc.) + static images
├── data/
│   └── projects.yaml
├── internal/
│   └── middleware/
│       ├── logger.go           # Logger middleware (slog)
│       ├── realip.go           # RealIP middleware
│       ├── trailingslash.go    # TrailingSlash middleware (new)
│       ├── recovery.go         # Recovery middleware (new)
│       └── clacks.go           # Clacks easter-egg middleware (new)
├── jslinux/                 # embedded as-is
├── posts/
│   └── {slug}/
│       ├── index.md
│       └── (assets)
├── scripts/
│   ├── build.sh
│   └── test.sh
├── static/                  # CSS, JS, SVGs, robots.txt, sitemap.xml (reads.js deleted)
├── templates/               # *.gohtml (error.gohtml added, rest unchanged)
├── github.go                # GitHub API client (no token, backoff strategy)
├── http.go                  # shared HTTP helpers (cache header, httpClient)
├── options.go               # flag/env parsing (--drafts only; GITHUB_TOKEN removed)
├── posts.go                 # post loading + handlers
├── projects.go              # project loading + handlers
├── slog.go                  # slog setup
├── tags.go                  # tag index + handlers
├── version.go               # version handler
├── website.go               # main(), mux setup, server lifecycle
└── go.mod
```

---

## Container image

The Dockerfile is removed. Image packaging uses `ocipack` declared as a Go tool dependency:

```
# go.mod
tool codeberg.org/fhuebner/ocipack/cmd/ocipack
```

`scripts/build.sh` runs:

```sh
CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o website .
go tool ocipack \
  -tag "ghcr.io/nothub/website:${TAG}" \
  website image.tar.gz
```

- `AddBinary` auto-detects arch from the ELF header.
- CA certificates are bundled automatically (needed for GitHub API calls).
- The resulting `image.tar.gz` is a standard OCI layout loadable by Docker, Podman, containerd, or k3s.

---

## Template setup

Templates are parsed once at startup into a `*html/template.Template` and passed as a
parameter to each `init*` function (`initPosts`, `initProjects`, `initTags`, etc.) rather
than stored in a package-level global. The parse glob is `templates/*.gohtml`.

After parsing, validate that every template name the code will call exists in the parsed set
(including `error.gohtml`); fail fast at startup if any are missing.

Template execution errors are logged at **Warn** level. The status code cannot be changed
once writing has started, so the log is the only recovery action. The one exception is a
missing template name, which is caught before writing begins and can return a proper 500 —
but startup validation should make this unreachable in practice.

---

## Error responses

A single shared helper handles all error responses:

```go
func writeError(w http.ResponseWriter, r *http.Request, code int, tmpl *template.Template)
```

- If the request `Accept` header includes `text/html` (parsed via `elnormous/contenttype`):
  write the status code, then render `error.gohtml` with the code as data. Template failures
  are logged at Warn level; no further action is possible.
- Otherwise: call `http.Error(w, http.StatusText(code), code)`, which sets
  `Content-Type: text/plain; charset=utf-8` and `X-Content-Type-Options: nosniff`.

`error.gohtml` receives the integer status code and renders an appropriate gopher image from
`static/gophers/` (404.svg, 429.svg, 500.svg). For unmapped codes it falls back to a plain
message with no image.

This helper is **not** used for mid-stream template failures — by that point the status code
is locked and the response is partially written. Those are Warn-logged and left as-is.

---

## GitHub API

Only `GET /repos/{owner}/{repo}` is used — a public endpoint requiring no token for public
repositories. `GITHUB_TOKEN` and the `--github-token` flag are removed entirely.

**Rate limit**: 60 requests/hour per source IP (unauthenticated). For a personal site with
~10–20 GitHub projects refreshed once per day this is trivially within budget.

**Request timeout**: 10 seconds (existing `httpClient` timeout, unchanged).

**Inter-request delay**: 30 seconds between successive project lookups in `fetchStars`.
This is a courtesy to GitHub, not a hard requirement given the daily refresh cadence.

**Backoff**: `githubRepoMeta` returns a typed `RateLimitError` (carrying the suggested
retry-after duration from the response headers) and never sleeps itself. `fetchStars` owns
all sleeping and retry decisions:

- On `RateLimitError`: wait the suggested duration (capped at 90 minutes), then retry once.
- On any other error: exponential backoff — 5s, 10s, 20s — up to 3 retries.
- After all retries exhausted: log a warning and skip the project for this refresh cycle;
  the previous star count is kept.

---

## RSS feed

Standard RSS 2.0 at `/rss.xml`. One `<item>` per non-draft post, sorted newest-first.
Fields: `<title>`, `<link>`, `<description>` (post `Desc` field), `<pubDate>` (RFC 1123Z).
Channel `<title>` and `<link>` are hardcoded to the site name and URL.

Posts are also sorted newest-first on `GET /posts` (this is a pre-existing TODO).

The feed includes a hidden XML comment paying respect to Aaron Swartz:

```xml
<!-- In memory of Aaron Swartz (1986–2013) — https://www.aaronsw.com -->
```

---

## Heading anchors

`goldmark-anchor` and `anchorTexter` (renders `¶` on h1/h2) are already implemented in
`posts.go`. Acceptance criterion: verify the anchors render and the `¶` links work once the
post handler is wired to a real route. No rework needed unless manual testing reveals a bug.

---

## Code style

- Standard Go: `gofmt`, tabs, no line-length limit.
- No comments unless the WHY is non-obvious.
- All handlers take `(w http.ResponseWriter, r *http.Request)` — no custom context types.
- Middleware uses the standard `func(http.Handler) http.Handler` signature.
- Error responses use `writeError` — HTML with gopher image for `text/html` clients,
  `http.Error` (plain text) for everything else. See the Error responses section.
- `log/slog` everywhere. Middleware that needs a logger receives `*slog.Logger` as a
  constructor parameter. The package-level `slogger` in `slog.go` is the instance passed in
  from `main`.

---

## Testing strategy

- Unit tests in `*_test.go` alongside each source file.
- Use `net/http/httptest` for handler tests — no external server needed.
- Priority test targets:
  - `parseMeta` (already complex, no tests exist)
  - Content negotiation logic
  - TrailingSlash middleware (various paths with and without trailing slash)
  - Route redirects (301 chains)
  - RSS feed output (well-formed XML, correct item count)
- No mocking of the filesystem; use `testing/fstest.MapFS` for post-loading tests.
- `scripts/test.sh` runs `go test ./...`.

---

## Boundaries

**Always do:**
- Keep the URL surface identical (no slug changes, no removed routes).
- Serve the site from a single statically-linked binary embedded with all assets.
- Produce a scratch-based OCI image via `ocipack`; the binary must be `CGO_ENABLED=0`.

**Ask first before:**
- Changing the HTML template structure or visual design beyond minor CSS fixes.
- Adding new dependencies to `go.mod`.
- Changing the `posts/` directory layout in a way that would break existing post slugs.

**Never do:**
- Reintroduce a third-party router (gin, chi, echo, etc.).
- Add a database or external state store.
- Break the existing `--drafts` flag.

---

## Out of scope

- Image map easter egg (explore after everything above is done).
- k8s / Hetzner deployment manifests (separate concern).
- New pages or routes beyond what is listed above.
