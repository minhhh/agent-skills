---
name: kb-to-md
description: Use when asked to summarize a body of knowledge, text, or links into structured markdown layouts (bullets, subsection, chapter-subsection, chapter), or when inserting/copying markdown sections to a target file.
---

# Knowledge to Markdown

## Overview

- Summarize information into a minimal, high-density markdown format matching specific structural hierarchies, and copy or merge the content into a target file.
- The default level of details is "full", which means all important details should be captured. User can also specify "lite" level of details, which should be shorter and more concise, however the english should still be standard and not broken

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

- **Trigger phrases**: "only bullet", "only bullets"
- **Structure**: No headings, just a bullet list of root points with nested
  subpoints. Sibling root bullet points must be separated from each other by
  at least one blank line.
- **Multiple Root Points rule**: Unless the user explicitly requests exactly
  one or a single root bullet point, the agent should feel free to use multiple
  root bullet points at the same level if it is logical to represent different
  topics, categories, or items. Do not force distinct topics under a single
  artificial root bullet point.
- **Bold highlight**: You should also choose to highlight important points
  first in the root point like in Format 2. Or you could choose to highlight
  the whole root point sentence, or you could choose to not highlight anything
  at all if the root point is trivial and does not require attention, like in format 1.

- Format 1:

  ```markdown
  * Description of root point 1
      * Sub-point level 1 detailing the root point.
          * Sub-point level 2 detailing sub-point level 1.

  * Description of root point 2
      * Sub-point level 1 detailing the root point.
          * Sub-point level 2 detailing sub-point level 1.
  ```

- Format 2:

  ```markdown
  * Root Point 1: Description of root point.
      * Sub-point level 1 detailing the root point.
          * Sub-point level 2 detailing sub-point level 1.

  * Root Point 2: Description of root point.
      * Sub-point level 1 detailing the root point.
          * Sub-point level 2 detailing sub-point level 1.
  ```

### 2. `subsection`

- **Trigger phrases**: "use subsection"
- **Structure**: A heading starting with `▼` containing bullet lists. Sibling
  root bullet points must be separated from each other by at least one blank
  line.

```markdown
▼ **Subsection Heading**

* **Root Point 1**: Description of root point 1.
    * Sub-point level 1.

* **Root Point 2**: Description of root point 2.
    * Sub-point level 1.
```

### 3. `chapter-subsection`

- **Trigger phrases**: "full chapter", "chapter and subsection"
- **Structure**: A chapter heading (`###`), subsections (`▼`), and bullet
  lists. Sibling root bullet points must be separated from each other by at
  least one blank line.

```markdown
### Chapter Heading

▼ **Subsection Heading**

* **Root Point 1**: Description of root point 1.
    * Sub-point level 1.

* **Root Point 2**: Description of root point 2.
    * Sub-point level 1.
```

### 4. `chapter`

- **Trigger phrases**: "flat chapter"
- **Structure**: Chapter headings (`###`) with bullet lists directly, no
  subsections. Sibling root bullet points must be separated from each other
  by at least one blank line.

```markdown
### Chapter Heading

* **Root Point 1**: Description of root point 1.
    * Sub-point level 1.

* **Root Point 2**: Description of root point 2.
    * Sub-point level 1.
```

## Spacing and List Styling Rules

- **Indentation for Sub-points**: Every sub-point level must be indented with
  exactly 4 spaces relative to its parent (`*`).
- **Blank Lines**: Sibling root bullet points must **always** be separated
  from each other by at least one blank line (exactly one blank line) to
  improve readability, even when they are part of the same list.
- **Heading Spacing**: Include exactly one blank line after `### Heading` or
  `▼ **Heading**` before starting the list.

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
