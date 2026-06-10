---
name: book-to-md
description: Use when asked to summarize, distill, or outline documentation, chapters, or book sections into detailed Markdown format. Builds on markdown-style-principles.
---

# Summarizing Books in Markdown

## Overview

- Prioritize high technical depth over concision.
- The summary should be as detailed as possible, DO NOT skip over any points
- If there're example code to illustrate one point, don't skip it
- Remove third-person way to describe information, but keep all the information articulate

## Prerequisites

**REQUIRED SUB-SKILL:** Load markdown-style-principles via `skill("markdown-style-principles")`
**REQUIRED BACKGROUND:** You MUST understand markdown-style-principles before using this skill.

> **Action required:** This skill does NOT automatically load `markdown-style-principles`. You must explicitly call the `skill` tool to load it before proceeding.

## When to Use

- Only when the user explicitly requires using this skill to summarize a book or document, or a section of them.

## Formatting Style Selection

- **Required Style**: This skill exclusively uses the **`chapter-subsection`** layout style defined in [markdown-style-principles].


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
