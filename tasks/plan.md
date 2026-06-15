# Rewrite Plan — hub.lol website

Source of truth: `SPEC.md`. Tasks are ordered by dependency; each must leave `go build .`
passing before the next begins.

---

## T01 — Housekeeping

**Files:** `Dockerfile`, `random.go`, `static/reads.js`, `options.go`, `github.go`

**What:**
- Delete `Dockerfile` (replaced by ocipack in T13)
- Delete `random.go` (dead — nothing calls `random.*`; Go 1.20+ auto-seeds the global source)
- Delete `static/reads.js` (Reads page is gone)
- Remove `--github-token` flag and `optGithubToken` var from `options.go`
- Remove the `Authorization` header branch from `github.go`

**Acceptance criteria:**
- None of the deleted symbols are referenced anywhere in `*.go` files
- `go build .` passes

**Verification:** `grep -r 'optGithubToken\|random\.' --include='*.go' .` returns nothing;
`go build .` exits 0.

**Deps:** none

---

## T02 — Middleware rewrite

**Files:** `internal/middleware/logger.go`, `realip.go`, `trailingslash.go`, `recovery.go`,
`clacks.go`

**What:**
- **logger.go** — rewrite: accept `*slog.Logger` constructor param; wrap `ResponseWriter` to
  capture status code; log structured JSON *after* calling `next` (so `r.RemoteAddr` reflects
  the real IP set by RealIP); fields: status, method, path, query, remote-addr, latency,
  user-agent; level Info <500, Error ≥500; never log Authorization or Cookie headers.
- **realip.go** — rewrite: accept `*slog.Logger` constructor param; parse `X-Forwarded-For`,
  take the leftmost non-empty value, set `r.RemoteAddr`; log a structured slog Warn if the
  header is absent.
- **trailingslash.go** — new: if `r.URL.Path != "/"` and ends with `"/"`, strip the slash
  and 301-redirect using `r.URL.String()` (preserves query string); else call next.
- **recovery.go** — new: accept `*slog.Logger`; defer/recover; log panic value + stack at
  Error level; write 500 if headers not yet sent.
- **clacks.go** — new: construct with `*rand.Rand` seeded at construction time (`rand.New(
  rand.NewSource(time.Now().UnixNano()))`); wrap ResponseWriter; after `next` returns,
  inject `X-Clacks-Overhead: GNU Terry Pratchett` with probability 1/42.

**Acceptance criteria:**
- `go build ./internal/middleware/` passes
- Each middleware function has signature `func X(args...) func(http.Handler) http.Handler`

**Verification:** `go build ./internal/middleware/` exits 0; unit tests in T11 will cover
behaviour.

**Deps:** T01 (clean state)

---

## T03 — Error template and helper

**Files:** `templates/error.gohtml`, `http.go`

**What:**
- **error.gohtml** — renders a full HTML page for error codes. Receives an `int` (the status
  code) as template data. Shows the matching gopher SVG for codes 404, 429, 500
  (`/static/gophers/{code}.svg`); falls back to a plain message for other codes.
- **http.go** — add `writeError(w http.ResponseWriter, r *http.Request, code int, tmpl
  *template.Template)`: use `elnormous/contenttype` to check if `text/html` is acceptable;
  if yes, `w.WriteHeader(code)` then `tmpl.ExecuteTemplate(w, "error.gohtml", code)`, log
  Warn on template failure; otherwise `http.Error(w, http.StatusText(code), code)`.

**Acceptance criteria:**
- `writeError` compiles and is callable with `(w, r, 404, tmpl)`
- `error.gohtml` parses without error when added to the template set

**Verification:** `go build .` passes; manually inspect that `error.gohtml` references correct
gopher paths.

**Deps:** none

---

## T04 — Post directory migration

**Files:** `posts/` (directory restructure), `posts.go`, `website.go`

**What:**
- Rename each `posts/{slug}.md` → `posts/{slug}/index.md`, creating the subdirectory.
  Existing slugs must not change.
- Update the `//go:embed` directive in `website.go`: change `posts/*` to `posts` (bare
  directory name — Go embed recurses into subdirs automatically; `posts/*` does not).
- Update `initPosts` in `posts.go` to:
  - `fs.ReadDir("posts")` to list slug directories
  - Skip entries that are not directories (log a warning)
  - Read `posts/{slug}/index.md` for the markdown source
  - Serve per-post asset files at `/posts/{slug}/{file}` by reading from the embedded FS
    subtree `posts/{slug}/`; skip `index.md` in this file list.

**Acceptance criteria:**
- All existing post slugs still reachable at the same URL
- `go build .` passes with the updated embed directive
- Entries lacking `index.md` are skipped with a log warning, not a fatal error

**Verification:** `go build .` exits 0; `go run . --drafts` loads posts without panic (manual
smoke test).

**Deps:** none

---

## T05 — Template startup wiring and init signatures

**Files:** `website.go`, `posts.go`, `projects.go`, `tags.go`, `version.go`

**What:**
- In `website.go` `main()`: parse templates with `template.Must(template.New("").
  ParseFS(fs, "templates/*.gohtml"))`; then validate that each required name exists in the
  parsed set (`error.gohtml`, `about.gohtml`, `posts.gohtml`, `post.gohtml`,
  `projects.gohtml`, `tags.gohtml`); `log.Fatalf` if any is missing.
- Update init function signatures:
  - `initPosts(mux *http.ServeMux, tmpl *html/template.Template) error`
  - `initProjects(mux *http.ServeMux, tmpl *html/template.Template) error`
  - `initTags(mux *http.ServeMux, tmpl *html/template.Template) error`
  - `initVersion(mux *http.ServeMux) error` (no template needed — returns plain text)
- Pass `tmpl` from `main()` into each call.

**Acceptance criteria:**
- All five files compile with the new signatures
- `go build .` passes
- Startup aborts with a clear fatal message if any template is missing

**Verification:** `go build .` exits 0; temporarily rename `error.gohtml` and confirm
`go run .` fatals with the expected message, then rename back.

**Deps:** T03 (error.gohtml must exist), T04 (embed glob updated)

---

## T06 — Core routes and server wiring

**Files:** `website.go`, `slog.go`, `http.go`

**What:**
- **slog.go** — ensure `slogger` is a package-level `*slog.Logger` (JSON handler to stderr,
  no flags). This is the instance passed into middleware constructors.
- **http.go** — add `setCacheHeader(w http.ResponseWriter)`: sets
  `Cache-Control: public, max-age=604800, immutable`.
- **website.go** — factor route wiring into `buildHandler(tmpl *html/template.Template,
  slogger *slog.Logger, clacksRand *rand.Rand) http.Handler` so tests can call it directly.
  The function reads from the package-level `fs embed.FS` (embed directives must be in
  package main). `main()` calls `buildHandler` with the parsed templates and slogger. Middleware chain (outermost → innermost):
  ```
  handler = middleware.Clacks(clacksRand)(mux)
  handler = middleware.TrailingSlash(handler)
  handler = middleware.Recovery(slogger)(handler)
  handler = middleware.RealIP(slogger)(handler)
  handler = middleware.Logger(slogger)(handler)
  ```
  Init call order inside `buildHandler` (tags depend on posts and projects):
  ```
  initPosts(mux, tmpl)
  initProjects(mux, tmpl)
  initTags(mux, tmpl)    // must come after initPosts and initProjects
  initVersion(mux)
  ```
  Routes registered on `mux`:
  - `"GET /{$}"` → `http.RedirectHandler("/about", 301)` (exact root; `{$}` prevents this
    from acting as a subtree catch-all for all GET requests)
  - `"GET /shell"` → `http.RedirectHandler("/jslinux", 301)`
  - `"GET /about"` → render `about.gohtml`
  - `"GET /assets/"` → `http.FileServerFS` subtree with cache header
  - `"GET /static/"` → same
  - `"GET /jslinux/"` → same
  - `"GET /robots.txt"` → serve `static/robots.txt` from embedded FS
  - `"GET /sitemap.xml"` → serve `static/sitemap.xml` from embedded FS
  - `"GET /teapot"` → cache header + 418 + `🫖`
  - `"/"` (no method) → `writeError(w, r, 404, tmpl)` catch-all

**Acceptance criteria:**
- Server starts and serves all listed routes
- `GET /unknown` returns 404 via `writeError`
- `GET /posts/` (trailing slash) returns 301 → `/posts`
- `GET /` returns 301 → `/about` (matched by `"GET /{$}"`, not by `"GET /"` subtree)

**Verification:** `go run . --drafts` starts without error; `curl -i localhost:8080/teapot`
returns 418; `curl -i localhost:8080/unknown` returns 404.

**Deps:** T02 (middleware), T05 (templates and init sigs)

---

## T07 — Posts handlers

**Files:** `posts.go`

**What:**
In `initPosts`, register on the mux (received as parameter):
- `"GET /posts"` — render `posts.gohtml` with posts sorted newest-first by
  `Meta.Date`. Content negotiation handled by T08 (RSS on this route added there).
- `"GET /posts/{slug}"` — look up slug via `r.PathValue("slug")`; 404 if not found; check
  Accept header via `elnormous/contenttype`:
  - `text/markdown` → `w.Header().Set("Content-Type", "text/markdown; charset=utf-8")`;
    serve raw `index.md` bytes from embedded FS.
  - default → render `post.gohtml`.
- `"GET /posts/rss.xml"` — literal pattern; 301 → `/rss.xml`. Specificity rules mean no
  registration-order dependency with `{slug}`.
- `"GET /posts/{slug}/"` — subtree handler serving per-post asset files from
  `posts/{slug}/` in embedded FS, skipping `index.md`. TrailingSlash middleware redirects
  `/posts/my-post/` → `/posts/my-post` before it reaches this handler; requests like
  `/posts/my-post/foo.png` (no trailing slash) pass TrailingSlash and hit this subtree.
  **Note:** asset paths in markdown must be absolute (`/posts/{slug}/filename`) since the
  canonical post URL has no trailing slash and relative `./foo.png` paths resolve to the
  wrong location.

**Note:** Go 1.22+ mux resolves `"GET /posts/rss.xml"` before `"GET /posts/{slug}"` by
specificity; no registration-order dependency.

**Acceptance criteria:**
- `GET /posts` returns 200 with HTML listing sorted newest-first
- `GET /posts/my-post` returns 200 HTML
- `GET /posts/my-post` with `Accept: text/markdown` returns raw markdown
- `GET /posts/nonexistent` returns 404
- `GET /posts/rss.xml` returns 301 → `/rss.xml`
- Heading anchors (`¶`) appear on h1/h2 in rendered posts

**Verification:** `curl -i localhost:8080/posts`; `curl -H 'Accept: text/markdown'
localhost:8080/posts/golang-shellscript`; confirm `¶` anchors in rendered HTML.

**Deps:** T06 (mux wired), T04 (post dirs and loading)

---

## T08 — RSS feed

**Files:** `posts.go`, `website.go`

**What:**
- Add RSS 2.0 handler function (can live in `posts.go`). Generates feed with:
  - XML comment as first child of `<rss>`: `<!-- In memory of Aaron Swartz (1986–2013) —
    https://www.aaronsw.com -->`
  - Channel `<title>`: `hub.lol`
  - Channel `<link>`: `https://hub.lol`
  - One `<item>` per non-draft post, sorted newest-first: `<title>`, `<link>`,
    `<description>` (post Desc), `<pubDate>` (RFC 1123Z)
  - `Content-Type: application/rss+xml; charset=utf-8`
- Register `"GET /rss.xml"` on the mux in `website.go`.
- In the `"GET /posts"` handler (T07), add content negotiation: if Accept includes
  `application/rss+xml` or `application/atom+xml`, serve the RSS feed directly.

**Acceptance criteria:**
- `GET /rss.xml` returns valid RSS 2.0 XML with Aaron Swartz comment present
- `GET /posts` with `Accept: application/rss+xml` returns same feed
- Feed items are sorted newest-first
- `xml.Unmarshal` (or `xmllint`) parses the output without error

**Verification:** `curl -s localhost:8080/rss.xml | xmllint --noout -`; grep for Aaron
Swartz comment; count `<item>` elements matches non-draft post count.

**Deps:** T07 (posts loaded and sorted)

---

## T09 — Projects handler

**Files:** `projects.go`

**What:**
- Restore yaml loading in `initProjects`: read `data/projects.yaml` using `gopkg.in/yaml.v3`
  (re-add import); unmarshal into `[]Project`; fatal on error.
- Register `"GET /projects"` on the mux; render `projects.gohtml` with project list.
- Keep the `fetchStars` goroutine and 24h ticker (from T02 GitHub API refactor).

**Acceptance criteria:**
- `GET /projects` returns 200 with rendered project list
- Star counts are fetched on startup (may be 0 until goroutine completes)
- `gopkg.in/yaml.v3` is back in `go.mod` as a direct dependency

**Verification:** `curl -i localhost:8080/projects` returns 200; check rendered HTML shows
project titles.

**Deps:** T06 (mux wired), T12 (github API backoff wired into fetchStars)

---

## T10 — Tags and Version handlers

**Files:** `tags.go`, `version.go`

**What:**
- **tags.go** `initTags`: register `"GET /tags"` → render `tags.gohtml` with full tag map;
  register `"GET /tags/{tag}"` → look up tag via `r.PathValue("tag")`; 404 if not found;
  render `tags.gohtml` with single-tag map.
- **version.go** `initVersion`: register `"GET /version"` → `fmt.Fprint(w, version)` with
  `Content-Type: text/plain; charset=utf-8`. Remove the `_ = version` blank identifier added
  in the gin-removal pass — the handler now uses `version` directly.

**Acceptance criteria:**
- `GET /tags` returns 200
- `GET /tags/go` returns 200 (assuming posts/projects with that tag exist)
- `GET /tags/nonexistent` returns 404
- `GET /version` returns plain text version string

**Verification:** `curl -i localhost:8080/tags`; `curl -i localhost:8080/version`.

**Deps:** T07 (posts register tags), T09 (projects register tags), T06 (mux wired)

---

## T11 — Tests

**Files:** `posts_test.go`, `internal/middleware/trailingslash_test.go`,
`internal/middleware/logger_test.go`, `website_test.go`

**What:**
- **posts_test.go** — `parseMeta`: table-driven tests covering valid input, missing fields,
  wrong types, bad date format; use `testing/fstest.MapFS` for post-loading tests that
  verify slug extraction, draft filtering, and newest-first sort order.
- **trailingslash_test.go** — table-driven: paths with trailing slash redirect, paths
  without pass through, root `/` passes through.
- **logger_test.go** — verify Logger does not log Authorization or Cookie header values.
- **website_test.go** — call `buildHandler` (from T06) with the real embedded FS and
  parsed templates to get an `http.Handler`; use `httptest.NewRecorder` to test: 301
  redirects (`/` → `/about`, `/shell` → `/jslinux`, `/posts/rss.xml` → `/rss.xml`); 404
  for unknown path; 418 for `/teapot`; RSS feed is well-formed XML with Aaron Swartz
  comment; item count matches non-draft posts.

**Acceptance criteria:**
- `go test ./...` exits 0
- No test uses `time.Sleep` or makes real network calls
- RSS test parses XML and asserts comment and item count

**Verification:** `bash scripts/test.sh` exits 0.

**Deps:** T08, T10 (all handlers complete)

---

## T12 — GitHub API refactor

**Files:** `github.go`, `projects.go`

**What:**
- **github.go** — define `RateLimitError` struct carrying a `RetryAfter time.Duration`; in
  `githubRepoMeta`, return `RateLimitError` on 403/429 (parse `Retry-After` header first,
  fall back to `X-RateLimit-Reset` epoch); never call `time.Sleep` inside this function.
- **projects.go** `fetchStars` — replace existing retry loop:
  - Inter-request delay: `time.Sleep(30 * time.Second)` between projects.
  - Per-project retry with exponential backoff: attempts at 0s, 5s, 10s, 20s (3 retries
    after first attempt).
  - On `RateLimitError`: sleep `min(err.RetryAfter, 90*time.Minute)`, then retry once.
  - After all retries exhausted: `slog.Warn(...)`, keep previous star count, move to next
    project.

**Acceptance criteria:**
- `githubRepoMeta` never calls `time.Sleep`
- `fetchStars` handles `RateLimitError` and generic errors with separate strategies
- `go build .` passes

**Verification:** `go build .` exits 0; code review that no `time.Sleep` appears in
`github.go`.

**Deps:** T01 (token removed)

**Note:** This task is listed last among infrastructure tasks because it has no downstream
blockers until T09 (projects handler). Can be done after T01 in practice.

---

## T13 — Container image

**Files:** `go.mod`, `scripts/build.sh`

**What:**
- Add ocipack as a Go tool dependency:
  `go get -tool codeberg.org/fhuebner/ocipack/cmd/ocipack`
- Rewrite `scripts/build.sh`:
  ```sh
  #!/usr/bin/env sh
  set -eu
  cd "$(dirname "$(realpath "$0")")/.."
  TAG="${TAG:-dev}"
  CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o website .
  go tool ocipack \
    -tag "ghcr.io/nothub/website:${TAG}" \
    website image.tar.gz
  rm -f website
  ```

**Acceptance criteria:**
- `go tool ocipack --help` exits 0 after `go mod download`
- `TAG=test bash scripts/build.sh` produces `image.tar.gz`
- `image.tar.gz` is loadable by `docker load` or `podman load`

**Verification:** Run `bash scripts/build.sh`; inspect output with
`go tool ocipack` or `tar tzf image.tar.gz`.

**Deps:** T01 (Dockerfile deleted)

---

## Dependency graph

```
T01 ──► T02 ──► T06 ──► T07 ──► T08 ──► T11
         │       │       │
T03 ─────┤       │       └──────────────► T10 ──► T11
         │       │
T04 ─────┴──► T05 ──► T06
                        │
T12 ──────────────────► T09 ──────────► T10

T13 (after T01, otherwise independent)
```

**Safe parallel groups:**
- T01, T03, T04 have no mutual deps — can start simultaneously
- T02 and T12 can both start after T01
- T07, T09, T10 (version only) can start simultaneously after T06
