## HTTP composition owned by internal/httpapp

**Decision:** Put complete request-serving composition in `internal/httpapp` (router, middleware policy/order, app/static/media routes) and test at that composed `http.Handler` seam. Keep `main.go` on process lifecycle only (flags, config load, logging, listener) and have it choose/root embedded static and host media filesystems before supplying them.

**Why:** Serving policy is one module’s contract, not bootstrap glue. Callers and tests should exercise the real composed handler rather than reassemble middleware and routes. Process concerns (how the binary starts and which FS roots exist) stay outside HTTP composition so lifecycle changes do not rewrite serving policy, and serving-policy changes do not touch process bootstrap.

**HTTP policy locked with this split:**
- Middleware order and scope are owned by the module.
- CORS allowed origins come from the environment through the existing config module; empty config means same-origin / no CORS; CORS applies to app routes only.
- Use standard Chi compression; do not add Brotli.
- Preserve existing observable error/status and `FileServer` behavior.
- Handlers stay stateless; `handlers.HomeHandler` may be removed.

**Rejected:** Leaving router/middleware assembly in `main` (untestable without a process harness; policy drifts with bootstrap). Pushing FS root selection into `httpapp` (couples composition to embed/host layout and blurs process vs request boundaries).

## CORS origins fail closed at config and composition

**Decision:** `CORS_ALLOWED_ORIGINS` is a strict comma-separated exact-origin list. Whitespace-only / unset means no CORS. When the whole value is non-empty, empty elements (leading, trailing, or repeated commas) are rejected. Any `*` anywhere in an origin is rejected (bare wildcard and host patterns). The same validation runs at the `httpapp.New` seam so caller-supplied wildcards cannot become chi’s allow-all policy; invalid input fails construction rather than degrading to a permissive default.

**Why:** Chi treats `*` in `AllowedOrigins` as match-all and also supports embedded wildcards. Skipping empty tokens or trusting caller slices would silently widen cross-origin access. Fail-closed at both env parse and composition matches “malformed never becomes permissive.”

**Rejected:** Skipping empty comma tokens as convenience; relying on chi defaults when the origin list is empty or tainted; normalizing garbage into “no CORS” after a non-empty malformed list.
