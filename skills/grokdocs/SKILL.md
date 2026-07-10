---
name: grokdocs
description: Use when searching the workspace for documentation, code patterns, or project knowledge — or after modifying files that should be re-indexed for search.
---

# grokdocs

## Overview

**grokdocs** is a local-first search engine indexing your workspace docs and source code for semantic and full-text search.

**Prerequisite:** `grokdocs` must be available in `PATH`.

## Quick Reference

| Action | Command |
|--------|---------|
| Search for knowledge | `grokdocs search "<query>"` |
| Re-index after file changes | `grokdocs sync` |
| Embed vectors (separate step) | `grokdocs embed` |

## Usage

### Searching

```bash
grokdocs search "how does skill linting work"
```

Flags:
- `-m, --mode <hybrid\|fts\|semantic>` — search strategy (default: `hybrid`)
- `--limit <n>` — max results (default: `5`)
- `-c, --collection <name>` — scope to a collection
- `-p, --project <path>` — explicit project root (auto-detected from cwd by default)

### Re-indexing After Changes

When files are added, edited, or deleted:

```bash
grokdocs sync
grokdocs embed   # optional: generate vector embeddings for semantic/hybrid search
```
