---
name: summarize-markdown
description: user ask "summarize", "distill" doc, chapter, book section → Markdown. Output deep technical depth + standard hierarchy.
---

# Summarizing Books in Markdown

## Overview
- Prioritize high technical depth over concision.
- The summary should be as detailed as possible, DO NOT skip over any points
- If there're example code to illustrate one point, don't skip it
- Remove third-person way to describe information, but keep all the information articulate

## When to Use
- Distilling technical chapters into comprehensive summaries where **nuance and specific details matter**.
- Creating standardized documentation for server troubleshooting or system architecture.
- When the target file requires a specific `###` and `▼` hierarchy.

## Core Pattern

### Formatting Hierarchy
| Element | Marker | Style |
|---------|--------|-------|
| Chapter | `###` | Standard heading |
| Sub-section | `▼` | **Bolded** title |
| Important Points| `*` | Primary list; use **bold** for key terms |
| Technical Details| `    *` | Nested list for nuances and sub-steps |
| Code/Commands | \`\`\` | Language-specific block |

### Example
```markdown
### Chapter Title

▼ **Sub-section Title**

* **Primary Key Point**: Capturing a high-level technical concept.
    * **Specific Nuance**: Explaining the "why" and "how" behind the point.
    * **Critical Exception**: Noting any important corner cases.
* `code_snippet` for practical application.
```

## Common Mistakes
- **Oversimplification**: Ignoring the technical "why" or "how" mentioned in the source text.
- **Marker mismatch**: Using `▽` instead of `▼`.
- **Nesting errors**: Putting `###` under `▼`. (Hierarchy: `###` -> `▼` -> `*`).
- **Missing context**: Failing to include the specific tools, commands, or parameters mentioned in the text.

## Red Flags - STOP and Start Over
- **Surface-level or vague summaries**: If the summary lacks technical depth or ignores specific "how-to" details.
- The final summary is less than one third of the original
- Using `#` or `##` for chapter headings.
- Omitting the `▼` marker for sub-sections.
- Writing long paragraphs instead of concise, detailed bullet points.

