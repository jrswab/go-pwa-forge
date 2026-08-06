# Anchor tags

Query with `grep -r '@anchor:'`.

| Surface | Role | Location | Context |
|---------|------|----------|---------|
| `main` | `lifecycle` | `main.go` | flags, config load, logging, FS rooting, listener |
| `config` | `load` | `config/config.go` | environment → process config including CORS origins |
| `httpapp` | `composition` | `internal/httpapp/app.go` | request-serving router, middleware policy, routes |
