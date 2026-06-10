---
name: kb-to-md
description: Use when asked to summarize a body of knowledge, text, or links into structured markdown layouts, or when inserting/copying markdown sections to a target file. Builds on markdown-style-principles.
---

# Knowledge to Markdown

## Overview

- Summarize information into a minimal, high-density markdown format matching specific structural hierarchies, and copy or merge the content into a target file.
- The default level of details is "full", which means all important details should be captured. User can also specify "lite" level of details, which should be shorter and more concise, however the english should still be standard and not broken

## Prerequisites

**REQUIRED SUB-SKILL:** Load markdown-style-principles via `skill("markdown-style-principles")`
**REQUIRED BACKGROUND:** You MUST understand markdown-style-principles before using this skill.

> **Action required:** This skill does NOT automatically load `markdown-style-principles`. You must explicitly call the `skill` tool to load it before proceeding.

## When to Use

- Translating links, references, or text into bulleted lists with nested details.
- Inserting or merging new technical sections into existing markdown
  documentation.
- When NOT to use: Do not use for simple code generation or non-markdown
  documentation.

## Layout Styles & Selection Rules

Select the formatting layout defined in [markdown-style-principles] based on natural language triggers in the user's prompt:

1. **`bullets`**
   - **Trigger phrases**: "only bullet", "only bullets"
   - **Selection**: Uses the `bullets` style.

2. **`subsection`**
   - **Trigger phrases**: "use subsection"
   - **Selection**: Uses the `subsection` style.

3. **`chapter-subsection`**
   - **Trigger phrases**: "full chapter", "chapter and subsection"
   - **Selection**: Uses the `chapter-subsection` style.

4. **`flat-chapter`**
   - **Trigger phrases**: "flat chapter"
   - **Selection**: Uses the `flat-chapter` style.


## Surgical Modification & Merging Guidelines (Core Workflow)

Whenever a target file or section already exists, do not perform a full
overwrite. Perform a surgical update:

1. Read the target file to locate the exact destination.
2. Locate matching `### Chapter` or `▼ **Subsection**` headings.
3. If found, merge new root points into the list, ensuring at least one blank
   line separates them from existing points. If not found, create the headings
   and append.
4. Use precise line replacement tools (`replace_file_content` or
   `multi_replace_file_content`) to avoid changing unrelated content.
