---
name: markdown-style-principles
description: >
  Use when a skill (like kb-to-md, book-to-md, or enrich-md-kb) references this as a Prerequisites foundation. NOT invoked directly by users.
---

# Markdown Style Principles

Shared formatting and structural rules for summarizing information into high-density Markdown layouts.

## Formatting Hierarchy & Layout Styles

### 1. `bullets`

* No headings; only bullet lists of root points and nested subpoints.

```markdown
* Root Point 1: Description of root point.
    * Sub-point level 1 detailing the root point.

* Root Point 2: Description of root point.
    * Sub-point level 1 detailing the root point.
```

### 2. `subsection`

* Sub-section heading starting with `▼` containing bullet lists.

```markdown
▼ **Subsection Heading**

* **Root Point 1**: Description.
    * Sub-point level 1.
```

### 3. `chapter-subsection`

* Chapter heading (`###`), sub-section headings (`▼`), and bullet lists.

```markdown
### Chapter Heading

▼ **Subsection Heading**

* **Root Point 1**: Description.
    * Sub-point level 1.
```

### 4. `flat-chapter`

* Chapter heading (`###`) with bullet lists directly, no sub-sections.

```markdown
### Chapter Heading

* **Root Point 1**: Description.
    * Sub-point level 1.
```

## Spacing and List Styling Rules

- **Bullet Markers**: Always use asterisks (`*`), never hyphens (`-`).
- **Indentation**: Sub-points must be indented with exactly 4 spaces relative to their parent.
- **Blank Lines**: Sibling root bullet points must be separated by exactly one blank line.
- **Sub-section Heading Format**: Sub-section headings starting with `▼` must always be bolded: `▼ **Heading**`.
- **Heading Spacing**: Exactly one blank line after `### Heading` or `▼ **Heading**`.
- **No Plain Paragraphs**: All text must be structured as root or nested bullet points.
- **Tone**: Direct, technical language. Avoid third-person descriptions (e.g., use "Metric X indicates..." instead of "The author explains...").
