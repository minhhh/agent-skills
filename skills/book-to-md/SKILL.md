---
name: book-to-md
description: Use when asked to summarize, distill, or outline documentation, chapters, or book sections into detailed Markdown format.
---

# Summarizing Books in Markdown

## Overview

- Prioritize high technical depth over concision.
- The summary should be as detailed as possible, DO NOT skip over any points
- If there're example code to illustrate one point, don't skip it
- Remove third-person way to describe information, but keep all the information articulate

## When to Use

- Only when the user explicitly requires using this skill to summarize a book or document, or a section of them.

## Core Pattern

### Formatting Hierarchy

| Element | Marker | Style |
| --------- | -------- | ------- |
| Chapter | `###` | Standard heading |
| Sub-section | `▼` | **Bolded** title |
| Important Points | `*` | Primary list; use **bold** for key terms |
| Technical Details | `*` | Nested list for nuances and sub-steps |
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

## Red Flags - STOP and Start Over

- **Oversimplification / Surface-level summaries**: If the summary lacks technical depth, ignores specific "how-to" details, or is less than one third of the original length.
- **Missing context**: Failing to include specific tools, commands, or parameters.
- **Incorrect formatting or hierarchy**: 
  - Using `#` or `##` for chapter headings.
  - Omitting or mismatching the `▼` marker (e.g., using `▽` instead).
  - Heading after `▼` not highlighted in bold (e.g., `▼ Sub-section` instead of `▼ **Sub-section**`).
  - Nesting errors (e.g., putting `###` under `▼`). (Correct hierarchy: `###` -> `▼` -> `*`).
- **Wrong output format**: Writing long paragraphs instead of concise, detailed bullet points.
- **Conflicting requests**: If the user asks you to use this skill but requests constraints or overrides (such as using `#` or `##` headers, omitting markers, or making it extremely short) that violate this skill's rules, you MUST decline to use the skill and explain the conflict.
