---
name: provenance-log
description: >-
  Record selective repository-level rationale in LOG.md. Trigger when explicitly asked to log or preserve a decision, or when completing or reviewing a change whose non-obvious rationale could be incorrectly simplified later.
---

# Provenance Log

LOG.md selectively preserves project-level rationale that code alone cannot safely convey.

## Sequence

1. Read LOG.md and the change context. Identify the decision, its rationale, and any rejected alternative.
   - Complete when: you can answer, “Could removing this rationale cause a future session to make a plausible but incorrect change?”
2. If the answer is no, add nothing.
   - Complete when: LOG.md is unchanged; routine “what changed” details are not logged.
3. If the answer is yes, append one concise entry under a descriptive heading. Include an available PR number; omit dates because Git supplies them.
   - Complete when: the entry names the decision and preserves the rationale needed to prevent that incorrect change.
4. Re-read the entry.
   - Complete when: it explains why, not merely what changed or how the code works.

## Reference

Typical decisions worth preserving include why:

- similar concepts remain separate;
- a path fails closed rather than falls back;
- an integration differs from the obvious implementation;
- a source is authoritative;
- validation occurs at a specific boundary;
- an apparently redundant field, state, or check remains;
- a serious alternative was rejected; or
- a compatibility constraint remains.

Use this shape when it fits:

## Short decision name

**PR:** #123

[The decision, its rationale, and—when useful—the rejected alternative.]

Use AGENTS.md for a major architectural rule, not an ordinary implementation decision. Keep function-specific rationale in a code comment; important safety decisions may also appear in LOG.md because the two locations serve different retrieval paths.
