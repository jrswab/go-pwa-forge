# Mem map

Index of `.mem/`. Read this first, then `terminology.md` and `summary.md`, then the files for the surface you change.

- [summary.md](summary.md) - project snapshot
- [terminology.md](terminology.md) - domain terms
- [practices.md](practices.md) - patterns that recur in the code
- [http/](http/summary.md) - HTTP request-serving surface
  - [composition.md](http/composition.md) - `httpapp` seam and the `main.go` boundary
  - [cors.md](http/cors.md) - CORS scope on app routes
- [config/](config/summary.md) - environment configuration
  - [cors-origins.md](config/cors-origins.md) - `CORS_ALLOWED_ORIGINS` parsing and validation
- [frontend/summary.md](frontend/summary.md) - Datastar-first UI
- [conventions/](conventions/summary.md) - repo rules for records and working docs
  - [anchor-tags.md](conventions/anchor-tags.md) - `@anchor:` tags and the `TAGS.md` index
  - [working-docs.md](conventions/working-docs.md) - gitignored `docs/plans/` and `docs/notes/`
  - [records-review.md](conventions/records-review.md) - records review and the stale-reference check
