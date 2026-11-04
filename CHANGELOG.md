# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added

- `.mem/` in the repository root: tracked project memory, indexed by `.mem/mem-map.md`, seeded from the HTTP composition, CORS, configuration, frontend, and anchor-tag decisions in the code.

### Changed

- `Makefile` `gen` and the deploy jobs in `.github/workflows/pr-workflow.yml` and `prod-workflow.yml` generate templ code with the version pinned in `go.mod` (v0.3.819), via `go run github.com/a-h/templ/cmd/templ@<version> generate`. The `Install Templ` step (`go install ...@latest`) is removed, so CI and local `make gen` use the same templ version.
- `config/config.go` `Load` is now a package function, `Load(getenv func(string) string) (Config, error)`, replacing `(*Config).Load()`. `main.go` passes `os.Getenv`; tests pass a map lookup. It returns a value and reads only `CORS_ALLOWED_ORIGINS`. CORS rules and fail-closed startup are unchanged.
- `AGENTS.md` restructured: the Datastar and working-docs sections are kept; the anchor-tag and `TAGS.md` rules moved into `TAGS.md`; the `.mem/` rules sit under `## Project Memory`; the records-review step points to `~/.agents/scripts/records-stale-check`; `make check` is the completion bar.
- `TAGS.md` now opens with the anchor-tag format and `TAGS.md` rules; the anchor table is unchanged.
- `AGENTS.md` no longer references `LOG.md`; decision rationale is recorded in `.mem/`.

### Removed

- `config/config.go` `Config`: `RedisHost`, `RedisPass`, `PostgresUser`, `PostgresPass`, `PostgresHost`, `PostgresPort`, `PostgresDB`, and their env reads (`REDIS_HOST`, `REDIS_PASS`, `DB_USER`, `DB_PASS`, `DB_HOST`, `DB_PORT`, `DB_NAME`). Redis and Postgres config comes back when an integration needs it.
- `Makefile`: `SERVER_ENV` (`POSTGRES_*` entries). `make run` is now `go run main.go -l`.
- `scripts/appd`: blank `REDIS_*` and `DB_*` exports in `rc_start`.
- `.agents/skills/provenance-log/SKILL.md`, and the `AGENTS.md` provenance-log section that pointed to it.
- `LOG.md`. Its decisions are recorded in `.mem/http/composition.md` and `.mem/config/cors-origins.md`; `.mem/` replaces it (see `.mem/conventions/decision-records.md`).

### Fixed

- `config/config.go` `Load`: the blank-variable warnings logged the whole `*Config`, including `RedisPass` and `PostgresPass`, to `logs/site.log`. `Load` no longer logs. `TestLoad_NeverLogsOrReportsSecrets` in `config/config_test.go` covers this.
- `main.go`, `main_test.go`, `handlers/decoder.go`, and `internal/httpapp/app_test.go`: ignored `Close` errors are now handled (`_ =` in production code, `t.Fatal` in tests), clearing golangci-lint errcheck.
- `handlers/untils.go`: removed the unused `setPersistentCookie` helper, clearing golangci-lint unused. It had no callers.
- `.github/workflows/pr-workflow.yml` and `prod-workflow.yml`: the deploy job built with `make pipeline` and `go build cmd/main.go`, both gone. It now runs `make gen` and `go build -o server .`. The npm Tailwind install step is removed; `static/output.css` is committed.
- `main.go`: a config error is now printed to stderr as well as logged to `logs/site.log`, before the process exits 1.
- `Makefile` `coverage` target: it called `.scripts/check_coverage.sh`, which is not tracked, so `make coverage` and `make prep` failed. It now runs `go tool cover -func=coverage.out` over the profile that `make test` writes.
- `README.md` Project Structure matches the tracked tree. Development lists the Makefile targets that exist; `make pipeline` is removed.
- `scripts/appd`: the `daemon` comment describes the deployed server binary path, not `cmd/main.go`.
- `views/layout/base_templ.go` is regenerated from `base.templ` with the pinned templ v0.3.819, so the committed page loads Datastar, not htmx. The htmx→Datastar replacement changed the source but left the generated file stale.
