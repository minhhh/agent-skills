# kb-to-markdown Golden Templates & Layout Rules

This document specifies the exact formatting rules and structural hierarchies required for the `kb-to-markdown` skill.

---

## 1. General Formatting & Layout Rules

* **Chapter Heading**: Use `### Heading Title` followed by exactly one blank line.
* **Subsection Heading**: Use `▼ **Heading Title**` (or `▼ Heading Title`) followed by exactly one blank line.
* **Root Points**: Use `* **Key Term**: Explanation` (bold the key term/concept, separate with colon).
* **Indentation for Sub-points**: Every sub-point level must be indented with exactly 4 spaces relative to its parent (`    *`).
* **Blank Lines between Lists**: Include exactly one blank line between sibling root points to improve readability, especially when they contain nested sub-points.
* **Heading Spacing**: Include exactly one blank line after `### Heading` or `▼ **Heading**` before starting the list.

---

## 2. Core Layout Styles & Natural Language Mapping

The skill dynamically selects the appropriate layout style based on natural language triggers in the user's prompt:

### Style 1: `bullets`
* **Trigger phrases**: "as a bullet list", "only bullets", "list format", "flat list", "root points"
* **Structure**: No headings, just a flat bullet list of root points with nested subpoints.
```markdown
* **Root Point**: Description of root point.
    * Sub-point level 1 detailing the root point.
        * Sub-point level 2 detailing sub-point level 1.
```

### Style 2: `subsection`
* **Trigger phrases**: "use subheadings", "use subsections", "subsection layout"
* **Structure**: A heading starting with `▼` containing bullet lists.
```markdown
▼ **Subsection Heading**

* **Root Point**: Description of root point.
    * Sub-point level 1.
```

### Style 3: `chapter-subsection`
* **Trigger phrases**: "full hierarchy", "use chapters and subsections", "nested sections"
* **Structure**: A chapter heading (`###`), subsections (`▼`), and bullet lists.
```markdown
### Chapter Heading

▼ **Subsection Heading**

* **Root Point**: Description of root point.
    * Sub-point level 1.
```

### Style 4: `chapter`
* **Trigger phrases**: "use chapter headings", "chapters only", "flat chapters"
* **Structure**: Chapter headings (`###`) with bullet lists directly, no subsections.
```markdown
### Chapter Heading

* **Root Point**: Description of root point.
    * Sub-point level 1.
```

---

## 3. Surgical Modification & Merging Rules (Core Workflow)

Whenever the target file or target section already exists, the skill must perform surgical modification/merging instead of a full rewrite:

1. **Locate Destination**: Read the target file to identify the appropriate insertion point.
2. **Heading Match**: Look for matching `### Chapter` or `▼ **Subsection**` headings.
3. **Merge**:
   - If found, append the new root points, ensuring a blank line separates them from existing points.
   - If not found, create the new headings and append them to the target file.
4. **Surgical Write**: Use `replace_file_content` or `multi_replace_file_content` to apply changes without affecting unrelated parts of the file.
