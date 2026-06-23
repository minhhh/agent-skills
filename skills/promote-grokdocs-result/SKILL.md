---
name: promote-grokdocs-result
description: Use when a grokdocs search ranks a relevant section too low and you want to boost it by adding incremental search tags as a bullet point under the target heading.
---

# Promote Grokdocs Result

## Overview
Boost a specific section's ranking in grokdocs by adding a `* **Tags:**` bullet point under its heading. Each search term you want to promote gets appended to this line. Tags stay in the section's chunk (grokdocs chunks by headings), improving both FTS5 BM25 and semantic embedding.

The tags line is natural markdown, invisible to rendered readers, and incrementally modifiable — keep adding terms over time.

## How grokdocs Chunking Works

grokdocs splits markdown files into **chunks at heading boundaries**. Each `#`/`##`/`###` heading starts a new chunk. A `### Go Slices` section is its own indexed chunk with its own BM25 score and embedding vector.

Adding a tags bullet under the heading seeds both:
- **FTS5 BM25** — higher term frequency in that chunk
- **Semantic embedding (FAISS)** — vector shifts toward the query

## Tag Format

```markdown
### Slices

* **Tags:** Go slices, dynamic array, slice header, backing array, append, capacity
```

This is:
- `*` — a bullet list item (natural markdown)
- `**Tags:**` — bold label clearly marking it as metadata
- Comma-separated keywords — easy to read, easy to append

## Workflow

### Step 1: Search
```bash
grokdocs search "<query>" --mode hybrid --limit 10
```

### Step 2: Identify the section to promote
```
[3] doc/go-slices.md > ### Go Slices [L45-L60] — score: 0.016
```

### Step 3: Add or append tags
Read the file, find the heading, and add a tags bullet if missing, or append new terms to the existing one:

```markdown
### Slices

* **Tags:** Go slices, dynamic array, slice header, backing array, append, capacity
```

Next time you search for a different concept (e.g., "nil slice"), just append:

```markdown
* **Tags:** Go slices, dynamic array, slice header, backing array, append, capacity, nil slice, reslice, make
```

### Step 4: Re-index
```bash
grokdocs sync --embed
```

### Step 5: Verify
```bash
grokdocs search "<query>" --limit 5
```

## Disambiguation

The tags line naturally disambiguates sections with the same heading across languages:

| Doc | Tags |
|-----|------|
| `go/slices.md` → `### Slices` | `* **Tags:** Go slices, slice header, backing array` |
| `python/sequences.md` → `### Slices` | `* **Tags:** Python slices, sequence, start:stop:step` |
| `js/array.md` → `### Slices` | `* **Tags:** JavaScript slice, shallow copy, array method` |

Each tags line stays in its own section's chunk and biases the embedding toward that language.

## Keyword Strategy

- **Lead with the disambiguator**: If searching for "Go slices", start tags with "Go slices"
- **Cover synonyms**: "dynamic array", "variable-length", "flexible array"
- **Include action/API terms**: "append", "make", "copy", "reslice" — these match how people search
- **One line, many terms**: Keep it on a single bullet. Append new terms with a comma.
- **5-15 terms** per line is fine. Don't pad with unrelated terms.

## Common Mistakes

- **Tags under the wrong heading**: Verify the heading text. Tags under `## Arrays` won't boost the `## Slices` chunk.
- **Multiple tags lines**: Keep one single `* **Tags:**` bullet. Multiple lines dilute per-line density.
- **Forgetting `--embed`**: `grokdocs sync` without `--embed` skips vector embedding.
- **Non-editable files**: Skip generated docs, vendored code, third-party content.
