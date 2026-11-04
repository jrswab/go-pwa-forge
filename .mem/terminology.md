# Terminology

- anchor - `// @anchor:<surface>:<role> [context]` comment that marks a boundary in source.
- surface - first part of an anchor, naming the boundary's area (e.g. `httpapp`).
- role - second part of an anchor, naming the kind of boundary (e.g. `composition`, `lifecycle`, `load`).
- context - one-clause description of the boundary. Copied from the anchor into `TAGS.md`.
- TAGS.md - index with one row per anchor. Not documentation.
- composed handler - the `http.Handler` returned by `httpapp.New`.
- file routes - `/static/*` and `/media/*`. Excluded from CORS.
- app routes - every route except file routes. CORS applies here.
- exact origin - `scheme://host[:port]` with no path, query, fragment, userinfo, or wildcard.
- working docs - `docs/plans/` and `docs/notes/`. Gitignored, local only.
- records review - check of `.mem/`, `TAGS.md`, `README.md`, and `CHANGELOG.md` against the code.
- records-stale-check - shared tool at `~/.agents/scripts/records-stale-check`. Lists records that mention names a diff removed or renamed.
- Datastar first - build behavior from `data-*` attributes and backend patches before writing a `<script>`.
- Templ - Go template language. `.templ` sources compile to tracked `_templ.go` files.
