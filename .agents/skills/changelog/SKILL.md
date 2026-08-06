---
name: changelog
description: Write a CHANGELOG.md entry for the current change under `## [Unreleased]`. Use before committing work to keep the changelog current.
license: MIT
metadata: https://keepachangelog.com
---

The changelog has two stages, owned by different skills:

- **This skill** — write one working-log entry per change under `## [Unreleased]`, before committing. The entry captures what landed, for the upgrader who will read it at release time.
- **The `release` skill** — distill the accumulated `## [Unreleased]` entries into a version section at release time (group by theme, collapse related entries, drop noise). That is not this skill's job; do not try to pre-distill.

## Steps

1. **Read the change you're about to commit.** Run both `git diff` (unstaged) and `git diff --cached` (staged) so neither half of the work is missed, and `git log` to see what landed since the last entry. The entry describes the *noteworthy difference* — what an upgrader or reviewer needs to know — not the commit message or the process of making the change.
2. **Pick the category.** Use the Keep a Changelog buckets:
   - `Added` — new features.
   - `Changed` — changes to existing functionality.
   - `Deprecated` — soon-to-be-removed features.
   - `Removed` — features removed now.
   - `Fixed` — bug fixes.
   - `Security` — vulnerabilities.
   If the `## [Unreleased]` section is missing a category sub-heading it needs, add it. Place categories in that order.
3. **Write one bullet.** Lead with the thing that changed (a file, a component, a behaviour), then the noteworthy difference. Cite file paths so the entry is navigable. One entry per logical change — if you made two unrelated changes, write two bullets. Match the voice of the existing entries (the repo's entries are factual, specific, and cite paths and symbols).
4. **Place it under `## [Unreleased]`.** Append to the appropriate category sub-heading. Never edit a released version section — those are frozen.
5. **Stop when the entry is written.** Do not reorganize, re-distill, or rewrite other entries. The release skill handles consolidation.

## What stays out of the working log

Not every change earns a bullet. Skip pure noise with no user-facing or developer-facing impact: whitespace, formatting, dead-code removal with no behaviour change. If in doubt, include it — the release skill will drop it during distillation if it truly doesn't belong.

## Reference

- Format: `## [X.Y.Z] - YYYY-MM-DD` for released sections; `## [Unreleased]` for the working log. Dates are ISO `YYYY-MM-DD` (largest-to-smallest units, unambiguous across locales).
- Comparison links live at the bottom of the file, one per version: `[X.Y.Z]: https://github.com/<owner>/<repo>/compare/vPREV...vX.Y.Z`. The release skill manages these.
- Full convention: [keepachangelog.com](https://keepachangelog.com/en/1.1.0/).
