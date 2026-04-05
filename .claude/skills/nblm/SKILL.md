---
name: nblm
description: >-
  ALWAYS use this skill when the user wants to ask, query, search, or retrieve information from their NotebookLM notebooks — or when they reference "my notebook(s)", "my NotebookLM", "my knowledge base", "my sources", "my research", or a specific notebook by topic (e.g. "check my beads notebook", "what do my notes say about X", "ask my research notebook about Y", "look up X in my stuff", "pull from my notes"). ALWAYS use this skill whenever the user's request could plausibly be answered by their own curated NotebookLM library rather than from general knowledge or the web. This skill runs the `nlm` CLI to return grounded answers with citations `[N]` pointing at sources the user themselves collected. Do NOT use for generic web search, Claude's internal knowledge, local file search, Google Drive queries, or generating new content (podcasts/videos/summaries) — this skill is read-only retrieval from the user's existing NotebookLM notebooks.
---

# NotebookLM RAG via `nlm` CLI

Answer the user's question by querying their own NotebookLM notebooks. Unlike web search or Claude's general knowledge, responses here are **grounded** — they come from sources the user has personally curated, with inline citations pointing at those sources. If the notebook can't answer, NotebookLM says so rather than guessing.

## Prerequisites

The `nlm` binary must be on PATH and authenticated. Verify with `nlm -limit 1 list`. If it fails:

- **Not installed**: tell the user it's missing and stop. Don't try to install it yourself — installation path depends on their Go setup and may need a specific fork/branch.
- **Auth expired** (output mentions auth or 401): tell the user to run `nlm auth` in their terminal. Auth is interactive (opens a browser) so an agent can't complete it.

## Core workflow

### Step 1 — Find the right notebook

List all notebooks to match against the user's request:

```
nlm -limit 0 list
```

The `-limit 0` flag is important: without it, only the first 10 of potentially dozens of notebooks show, and you may miss the relevant one.

Output looks like:
```
Total notebooks: 35

ID                                   TITLE                                       SOURCES LAST UPDATED
a0a01e3b-0eab-4ce6-8b26-f84668a5d4fa 📋 Erggo: Structured Task Backlogs for ...  4       2026-03-26T15:21:14Z
f921dea2-089a-49d4-9cb0-e59f2d9418c2 📿 Beads Best Practices: A Guide to Age... 2       2026-01-01T17:12:02Z
```

**Matching titles**: users name notebooks casually ("my beads notebook"), so use case-insensitive substring matching against the title column. Titles usually start with an emoji plus topic keywords and get truncated at ~40 chars in the table. A typical pipeline:

```bash
nlm -limit 0 list | grep -i "beads"
```

The first column is the notebook ID — that's what you pass to `generate-chat`.

**When multiple notebooks match**: prefer the one with more sources or the most recent `LAST UPDATED` timestamp, but if it's genuinely ambiguous, list the candidates and ask the user which one they meant. Wrong notebook → wrong answer, so it's worth a clarifying question.

**When nothing matches**: tell the user no notebook matches their description, show the closest few titles, and ask. Don't silently pick something unrelated.

### Step 2 — Ask the question

```
nlm generate-chat <notebook-id> "<prompt>"
```

Pass the user's question as-is (or lightly rephrased for clarity). Be generous with context in the prompt — NotebookLM handles full natural-language questions well, including multi-part ones.

Two output streams:
- **stderr**: one line `Generating response for: <prompt>` — progress, ignore
- **stdout**: the grounded answer, with citations like `[1]`, `[1, 3]`, `[1-4]` referring to sources in that notebook

Latency is typically 5–30 seconds. Don't timeout aggressively.

**Quoting**: wrap the prompt in double quotes and escape any embedded doubles, or use a here-string. Prompts with apostrophes work inside double quotes without escaping.

### Step 3 — Relay the answer

Return the answer **verbatim, preserving the `[N]` citations** — they are load-bearing. Users rely on them to audit claims against the source list. Don't summarize away the citations or replace them with your own reasoning.

If the user asks "what are citations [1] and [2]?" or "where does that come from?", resolve them with:

```
nlm sources <notebook-id>
```

Citation `[N]` corresponds to the Nth source in that list (1-indexed).

## Command reference

**Read-only (safe for agents):**
- `nlm -limit 0 list` — all notebooks
- `nlm -limit N list` — first N (default 10)
- `nlm sources <notebook-id>` — sources in a notebook (for citation resolution)
- `nlm notes <notebook-id>` — user-curated notes saved inside the notebook
- `nlm generate-chat <notebook-id> "<prompt>"` — RAG query

**Mutating (confirm with user first, or don't run):**
- `nlm create`, `nlm add`, `nlm rm`, `nlm rm-source`, `nlm rm-note` — modify user data
- `nlm audio-create`, `nlm video-create`, `nlm generate-*` (other than `generate-chat`) — kick off async generation jobs that write back into the notebook

**Skip entirely:**
- `nlm chat` — interactive REPL, will hang an agent session
- `nlm auth` — browser-based, needs the human

## Troubleshooting

- **HTTP 400 from chat**: Google occasionally changes the chat payload format. If it fails consistently, the installed `nlm` may be stale. Tell the user.
- **HTTP 429 / rate limits**: Google throttles heavy use. Back off; don't retry in a tight loop.
- **"Authentication required"**: cookies expired (every 2–4 weeks). User must run `nlm auth`.
- **Empty or truncated answer**: the notebook sources may not contain enough on the topic. Try a broader question or check `nlm sources` to see what's actually in there.

## Examples

**"What do my notes say about SLICOT Riccati solvers?"**
```bash
nlm -limit 0 list | grep -i slicot
# → 22b4774a-...  🧑‍🏫 Single Loop Control Methods: ...  21  ...
nlm generate-chat 22b4774a-0579-48e7-a5b0-a245f37e3ecd \
  "What do my sources say about SLICOT Riccati solvers?"
```

**"Check my beads notebook — how do dependency edges work?"**
```bash
nlm -limit 0 list | grep -i beads
# (if multiple matches, pick most-sources / most-recent or ask)
nlm generate-chat <id> "How do dependency edges work?"
```

**"Ask my research notebook about X and tell me which source says it"**
```bash
nlm generate-chat <id> "X"
# Answer comes back with e.g. [2, 3]
nlm sources <id>
# Read rows 2 and 3 from the output and cite them back to the user
```
