# enrich-md-kb Golden Templates & Layout Rules

Reference templates and exact formatting rules for the `enrich-md-kb` skill.
Canonical spacing and marker rules live in
[markdown-style-principles]; this file pins down the concrete structural
hierarchies and the correct shape for each layout.

---

## 1. General Formatting & Layout Rules

* **Chapter Heading**: Use `### Heading Title` followed by exactly one blank line.
* **Subsection Heading**: Use `▼ **Heading Title**` followed by exactly one blank line. Always bold the heading text.
* **Root Points**: Use `* **Key Term**: Explanation` (bold the key term/concept, separate with a colon).
* **Indentation for Sub-points**: Every sub-point level must be indented with exactly 4 spaces relative to its parent (`    *`).
* **Blank Lines between Lists**: Include exactly one blank line between sibling root points, especially when they contain nested sub-points. Sibling root bullet points must **always** be separated by at least one blank line.
* **Heading Spacing**: Include exactly one blank line after `### Heading` or `▼ **Heading**` before starting the list.
* **No Plain Paragraphs**: All content must be a root or nested bullet point, or a numbered list item carried over from a numbered source.
* **Numbered Lists**: Preserve `1. **Heading**` style items from the source; do not convert them to bullets.

---

## 2. Core Layout Styles & Natural Language Mapping

Select the layout from natural language triggers in the user's prompt:

### Style 1: `bullets`
* **Trigger phrases**: "as a bullet list", "only bullets", "list format", "flat list", "root points"
* **Structure**: No headings, just a flat bullet list of root points with nested subpoints. Sibling root points must be separated by at least one blank line.
```markdown
* **Root Point 1**: Description of root point 1.
    * Sub-point level 1 detailing the root point.
        * Sub-point level 2 detailing sub-point level 1.

* **Root Point 2**: Description of root point 2.
    * Sub-point level 1 detailing the root point.
        * Sub-point level 2 detailing sub-point level 1.
```

### Style 2: `subsection`
* **Trigger phrases**: "use subheadings", "use subsections", "use subsection", "subsection layout"
* **Structure**: A heading starting with `▼` containing bullet lists. Sibling root points must be separated by at least one blank line.
```markdown
▼ **Subsection Heading**

* **Root Point 1**: Description of root point 1.
    * Sub-point level 1.

* **Root Point 2**: Description of root point 2.
    * Sub-point level 1.
```

### Style 3: `chapter-subsection`
* **Trigger phrases**: "full hierarchy", "full chapter", "chapter and subsection", "use chapters and subsections", "nested sections"
* **Structure**: A chapter heading (`###`), subsections (`▼`), and bullet lists. Sibling root points must be separated by at least one blank line.
```markdown
### Chapter Heading

▼ **Subsection Heading**

* **Root Point 1**: Description of root point 1.
    * Sub-point level 1.

* **Root Point 2**: Description of root point 2.
    * Sub-point level 1.
```

### Style 4: `flat-chapter`
* **Trigger phrases**: "use chapter headings", "chapters only", "flat chapter", "flat chapters"
* **Structure**: Chapter headings (`###`) with bullet lists directly, no subsections. Sibling root points must be separated by at least one blank line.
```markdown
### Chapter Heading

* **Root Point 1**: Description of root point 1.
    * Sub-point level 1.

* **Root Point 2**: Description of root point 2.
    * Sub-point level 1.
```

---

## 3. Surgical Modification & Merging Rules (Core Workflow)

Whenever the target file or target section already exists, perform surgical modification/merging instead of a full rewrite:

1. **Locate Destination**: Read the target file to identify the appropriate insertion point.
2. **Heading Match**: Look for matching `### Chapter` or `▼ **Subsection**` headings.
3. **Merge**:
   - If found, append the new root points, ensuring at least one blank line separates them from existing points. Sibling root points must remain separated by a blank line.
   - If not found, create the new headings and append them to the target file.
4. **Surgical Write**: Apply the change with a precise line-replacement tool (e.g. `replace_file_content` or `multi_replace_file_content`) so unrelated parts of the file are untouched.

---

## 4. Worked Example: Summarizing Into a Layout

**Source (unstructured notes):**

> Nginx worker tuning. worker_processes should match CPU cores, or use auto. worker_connections sets max simultaneous connections per worker; max clients ≈ worker_processes × worker_connections. Also enable keepalive_timeout.

**Prompt:** "Summarize this as a flat chapter."

**Result:**

```markdown
### Nginx Worker Tuning

* **worker_processes**: Set to the number of CPU cores, or use `auto` to let Nginx detect it.
    * `worker_processes auto;`: Preferred on hosts that change CPU count.

* **worker_connections**: Sets the maximum simultaneous connections per worker.
    * **Max clients**: Approximately `worker_processes × worker_connections`.

* **keepalive_timeout**: Enable to reuse client connections and cut handshake overhead.
```

[markdown-style-principles]: ../../markdown-style-principles/SKILL.md
