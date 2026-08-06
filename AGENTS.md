# AGENTS.md

## Working Docs

- Use `./docs/plans/` for ephemeral, non-code planning and implementation working documents.
- Use `./docs/notes/` for durable project context that should remain useful across agent sessions and context-window resets.
- Files under `./docs/plans/` and `./docs/notes/` are local and gitignored. They are not committed or shared through Git and may be removed when the local workspace is cleaned up.
- Keep canonical product and architecture documents directly under `./docs/` unless they are explicitly plans or notes.

## Anchor Tags

Anchor tags are structured metadata comments that make the repository queryable:

```go
// @anchor:<surface>:<role> [context]
```

- Use native plain-comment syntax.
- Reuse established surfaces and roles when they describe the same boundary.
- Tag meaningful boundaries—handlers, service entry points, models,
  configuration, and integrations—not every helper.
- Keep context brief.
- Maintain root `TAGS.md` when adding, removing, or changing boundaries.

When planning work, start with `TAGS.md`, identify affected surfaces, and query
their anchors (for example, `grep -r "@anchor:billing"`).

## Provenance Log

`LOG.md` is the canonical, tracked record of non-obvious project-level decisions
and rationale; it is not a changelog. Read relevant entries before changing code
and treat them as constraints.

When work establishes, changes, or rejects a decision that a future contributor
could plausibly misunderstand, record the rationale in `LOG.md` before completion.
Explain why the decision was made and, when useful, which serious alternative was
rejected. Do not substitute `docs/plans/`, `docs/notes/`, commit messages, or code
comments for project-level rationale that belongs in `LOG.md`.

Before completing or reviewing work, use `.agents/skills/provenance-log/SKILL.md` when
rationale may need recording in `LOG.md`. If the skill determines that no entry is
needed, leave `LOG.md` unchanged; routine progress and ordinary implementation
details do not belong there.

## Completion

Before completion:

- Run the relevant checks and review the final diff.
- Do not expose secrets or real PII.
- Update `TAGS.md` and `LOG.md` when their rules require it.
- Report checks run and any checks that could not be run.
