# Anchor tags

Anchor tags are structured comments that make the repository queryable:

```go
// @anchor:<surface>:<role> [context]
```

- Use native plain-comment syntax.
- Reuse established surfaces and roles when they describe the same boundary.
- Tag meaningful boundaries (handlers, service entry points, models,
  configuration, and integrations), not every helper.
- Context names the boundary in present tense, describing current state, in
  one clause of 80 characters or fewer.
- Put detail where it belongs: behavior in `.mem/`, history in Git, and
  release notes in `CHANGELOG.md`.

`TAGS.md` is an index, not documentation. It has one row per anchor: surface,
role, file, and context.

- Copy the context from the anchor itself. Add nothing else.
- When an anchor changes, rewrite its row to describe the current state.
- Update this file when you add, remove, or change an anchor.

When planning work, start here: identify affected surfaces, then query their
anchors.

Query with `grep -r '@anchor:'`.

| Surface | Role | Location | Context |
|---------|------|----------|---------|
| `main` | `lifecycle` | `main.go` | flags, config load, logging, FS rooting, listener |
| `config` | `load` | `config/config.go` | environment → process config including CORS origins |
| `httpapp` | `composition` | `internal/httpapp/app.go` | request-serving router, middleware policy, routes |
