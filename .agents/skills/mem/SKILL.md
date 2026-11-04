---
name: mem
description: Rules and steps for managing project knowledge. All persistent project memory lives in a structured, AI-owned markdown repository located at `.mem/`. The mem is the AI's perfect memory and the only way to stay aligned over weeks/months.
disable-model-invocation: false
---

### Core Principles You Never Break

- The human owns the code and makes final decisions. You are the memory and high-speed executor.
- Anything worth implementing is worth permanently recording in the `.mem` directory.
- The mem directory is for YOU (the AI). Summarize mem contents rather than dumping them verbatim, unless the user requests a specific file by path.

### Authority inside .mem/

- You may freely create, update, rename, move, or delete files.
- You may create new top-level directories when the project evolves.
- You may delete a file only if it exists in the repo and has no uncommitted changes.
- All diagrams must be Mermaid only.
- If mem content contradicts actual code, summarize the disparity, prioritize the code as the source of truth, and ask the user to confirm your suggested mem fix.

### Mandatory structure (create missing parts as needed)

.mem/
    summary.md               # one-paragraph living snapshot
    terminology.md           # a repository of short (term - meaning) lines describing the domain language
    practices.md             # patterns and practices relevant to this project
    mem-map.md               # hierarchical index of all mem files
    plans/                   # roadmaps & TODOs
    tmp/                     # git-ignored session scraps
    [any-domain]/            # e.g. parser/, auth/, ui/, billing/
        summary.md + *.md    # one focused topic per file (kebab-case)

### Every mem file must

- cover exactly one topic
- contain concrete code examples + Mermaid diagrams
- link to related mems with relative paths
- document invariants, contracts, rationale, and lessons learned
- stay under 250 lines; if larger, decompose into focused sub-files

### Mandatory workflow (gently enforce)

1. Seed sessions with the most relevant mem files.
2. Use chat mode for exploration and design; never jump straight to code.
3. Implement only after a clear decision.
4. The instant the user says "looks good / ship it / this is final", immediately update or create the corresponding mem entries so the mem reflects reality.
5. After big changes, check if mem structure still mirrors the codebase and refactor if needed.

### Important Behaviours

- Session scraps go in .mem/tmp/ (git-ignored)
- Only permanent learnings go in main mem files
- If you're documenting something you'll need in future sessions, it goes in the mem
- If it's just 'how I solved today's problem,' it stays in chat
- Information in the mem directory is a description of the current state of the system.
- There is no need to leave behind summaries of completed work. Instead, update the mem appropriately.
- Your performance over time is determined by the quality of your code and the mem.
- After completing any user request that modifies code behavior or structure, immediately update the corresponding mem file before moving to the next task.
- Your success is measured by mem accuracy after each session: the mem must reflect current system state, not a history of changes.

### Example - mem entry after adding retry logic to API client:

BAD (changelog style):
  "Added retry logic to api-client.ts on 2024-01-15. Previously requests
   would fail immediately. Now they retry 3 times with exponential backoff."
GOOD (current state):
  "The API client retries failed requests up to 3 times with exponential
   backoff (100ms, 200ms, 400ms). Retries apply only to 5xx and network
   errors; 4xx responses fail immediately."
If you need to capture changelog-style information, save it in mem/tmp/.

### Session handoffs

When the user requests a handoff, create a handoff document in mem/tmp/. This document should provide all relevant knowledge from the current session that will be useful when resuming in a fresh session, including current task state, decisions made, approaches tried, blockers encountered, and next steps. The goal is to let a fresh session continue seamlessly without losing momentum.

### New Sessions

At session start, read mem/mem-map.md, mem/terminology.md, and mem/summary.md.
IMPORTANT: Before exploring the codebase or searching for files, ALWAYS check mem/mem-map.md first. It's your index to all project documentation. Use it to find relevant mem files before diving into code.
When the session starts, briefly show that you have domain knowledge before attending to the first request.
if the mem/ does not exist, ask the user if you should create one.
