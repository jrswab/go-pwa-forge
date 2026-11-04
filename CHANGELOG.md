# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added

- `.mem/` in the repository root: tracked project memory, indexed by `.mem/mem-map.md`, seeded from the HTTP composition, CORS, configuration, frontend, and anchor-tag decisions in the code.

### Changed

- `AGENTS.md` restructured: the Datastar and working-docs sections are kept; the anchor-tag and `TAGS.md` rules moved into `TAGS.md`; the `.mem/` rules sit under `## Project Memory`; the records-review step points to `~/.agents/scripts/records-stale-check`; `make check` is the completion bar.
- `TAGS.md` now opens with the anchor-tag format and `TAGS.md` rules; the anchor table is unchanged.
- `AGENTS.md` no longer references `LOG.md`; decision rationale is recorded in `.mem/`.

### Removed

- `.agents/skills/provenance-log/SKILL.md`, and the `AGENTS.md` provenance-log section that pointed to it.
