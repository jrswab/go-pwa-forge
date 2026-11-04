# AGENTS.md

## Frontend: Datastar

Load the `datastar` skill before changing UI, templates, or handlers that render or stream to the page.
If the skill is unavailable, read https://data-star.dev/reference/attributes first.
Datastar first: build each behavior from `data-*` attributes and backend patches.
Write a `<script>` only for behavior that neither can implement.

## Working Docs

- Files under `./docs/plans/` and `./docs/notes/` are local and gitignored: not
committed, not shared through Git, and may be removed when the local workspace
is cleaned up.
- Use `./docs/plans/` for ephemeral, non-code planning and implementation working documents.
- Use `./docs/notes/` for local context across agent sessions.
- Keep canonical product and architecture documents directly under `./docs/` unless they are explicitly plans or notes.

## Before Committing

- Use `.agents/skills/changelog/SKILL.md` to write the `CHANGELOG.md` entry.
- Update `README.md` when setup, commands, configuration, or documented behavior changes.
- Update `.mem/` when the change establishes or changes a decision.
- Run `~/.agents/scripts/records-stale-check` and assess each hit against current code.
  A records review checks `.mem/`, `TAGS.md`, `README.md`, and `CHANGELOG.md` against the code.

## Anchor Tags

Anchor tags (`@anchor:`): read `TAGS.md` before planning work or adding/changing an anchor.

## Project Memory

`.mem/` is the canonical, tracked record of project decisions and rationale.
Read `.mem/mem-map.md` first, then the relevant domain files, and treat them as constraints before changing a surface.
Use `.agents/skills/mem/SKILL.md` for the workflow. Describe current state, not history: Git keeps superseded detail.

## Completion

Before completion:

- Run `make check`; it passes.
- Review the final diff.
- Do not expose secrets or real PII.
