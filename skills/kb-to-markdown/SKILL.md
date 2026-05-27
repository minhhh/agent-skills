---
name: kb-to-markdown
description: Use when asked to summarize a body of knowledge, text, or links into structured markdown layouts (bullets, subsection, chapter-subsection, chapter), or when inserting/copying markdown sections to a target file.
---

# KB to Markdown

## Overview

Summarize information into a minimal, high-density markdown format matching
specific structural hierarchies, and copy or merge the content into a target
file.

## When to Use

- Translating links, references, or text into bulleted lists with nested details.
- Inserting or merging new technical sections into existing markdown
  documentation.
- When NOT to use: Do not use for simple code generation or non-markdown
  documentation.

## Layout Styles & Selection Rules

Select the formatting layout based on natural language triggers in the user's
prompt:

### 1. `bullets`

- **Trigger phrases**: "only bullet", "flat list"
- **Structure**: No headings, just a flat bullet list of root points with
  nested subpoints.

```markdown
* **Root Point**: Description of root point.
    * Sub-point level 1 detailing the root point.
        * Sub-point level 2 detailing sub-point level 1.
```

### 2. `subsection`

- **Trigger phrases**: "use subsection"
- **Structure**: A heading starting with `▼` containing bullet lists.

```markdown
▼ **Subsection Heading**

* **Root Point**: Description of root point.
    * Sub-point level 1.
```

### 3. `chapter-subsection`

- **Trigger phrases**: "use chapters and subsections", "nested sections"
- **Structure**: A chapter heading (`###`), subsections (`▼`), and bullet lists.

```markdown
### Chapter Heading

▼ **Subsection Heading**

* **Root Point**: Description of root point.
    * Sub-point level 1.
```

### 4. `chapter`

- **Trigger phrases**: "flat chapter"
- **Structure**: Chapter headings (`###`) with bullet lists directly, no
  subsections.

```markdown
### Chapter Heading

* **Root Point**: Description of root point.
    * Sub-point level 1.
```

## Spacing and List Styling Rules

- **Indentation for Sub-points**: Every sub-point level must be indented with
  exactly 4 spaces relative to its parent (`*`).
- **Blank Lines**: Always include exactly one blank line between sibling root
  points to improve readability, even when they are part of the same list.
- **Heading Spacing**: Include exactly one blank line after `### Heading` or
  `▼ **Heading**` before starting the list.

## Surgical Modification & Merging Guidelines (Core Workflow)

Whenever a target file or section already exists, do not perform a full
overwrite. Perform a surgical update:

1. Read the target file to locate the exact destination.
2. Locate matching `### Chapter` or `▼ **Subsection**` headings.
3. If found, merge new root points into the list, ensuring exactly one blank
   line separates them from existing points. If not found, create the headings
   and append.
4. Use precise line replacement tools (`replace_file_content` or
   `multi_replace_file_content`) to avoid changing unrelated content.
