# CORS on app routes

CORS is an app-route policy. It is mounted only when `CORS_ALLOWED_ORIGINS` yields at least one origin. Origin rules are in [../config/cors-origins.md](../config/cors-origins.md).

## Behavior

- Mounted at mux level with `r.Use`, not on a route group. Group-scoped CORS misses `OPTIONS` preflight for method-not-allowed routes. The comment in `newApp` states this.
- Skipped for `/static/*` and `/media/*` by `isFileRoute`.
- Allowed methods: `GET`, `OPTIONS`. Allowed headers: `Accept`, `Content-Type`, `Origin`, `hx-request`.
- `AllowOriginFunc` is nil. The origin slice is copied at construction, so callers cannot change the policy afterward.
- chi's default `*` origin is never used.

```mermaid
flowchart LR
  req["Request"] --> isFile{"file route?"}
  isFile -- yes --> next["next handler"]
  isFile -- no --> cors["chi cors handler"]
  cors --> next
```

Related: [composition.md](composition.md).
