---
name: writing-plans-lite
description: Use when initiating or maintaining a lightweight plan in a single markdown PRD file
---

# Writing Plans Lite

A workflow for writing, editing, and maintaining requirements and task lists in
a single markdown file (`prd-mvp.md` or `prd-[feature].md`). Designed for
autonomous execution with minimal process overhead.

## When to Use

- Use when:
  - The user explicitly requests a lightweight planning structure, or activates
    it with keywords like "lite plan", "plan lite", or "writing-plans-lite".
  - You proactively determine that a new feature or task warrants planning, but
    the standard planning mode (implementation plan) is too heavy.
  - The agent completes a task in a PRD and needs to update and archive it.

## Template Reference

The base template is located at [prd-template.md](assets/prd-template.md).

## Continuous Workflow

1. **Bootstrap**: Create `.tasks/prd-mvp.md` (or `prd-[feature].md` for
   sub-features) using the template.
   - **CRITICAL ORDER**: You MUST create/bootstrap this PRD file *before*
     writing any implementation code or making edits to source files. No
     exceptions.
2. **Progress Execution**: Update checklists as you complete work.
3. **Dashboard Updates**: Update the `Active Dashboard` using status symbols:
   - `[ ]` Todo
   - `[/]` In-Progress
   - `[x]` Completed
4. **Archiving (Critical)**: When a task is completed, immediately cut its entry
   from Section 2 (Active Dashboard) and its details from Section 3 (Active Task
   Details), and append them to **Section 5 (History / Archive)**.

## Usage Tips

- **Safeguard Item**: If you want to keep a task active (In-Progress) while you
  think, gather context, or wait for input, you can optionally add a final
  safeguard item to the task checklist: `- [ ] Wait for further instruction
  from user`.

## Common Mistakes & Loops

<!-- markdownlint-disable MD013 -->
| Excuse / Rationalization | Reality |
| --- | --- |
| "The user requested database/types first; I can write the PRD later." | **Violation.** You must bootstrap the plan file first. Coding without a plan leads to design mismatch and lost context. |
| "This is too small to need a plan file." | If it requires multiple steps, edits to multiple files, or has any ambiguity, it needs a plan file. |
| "I'll archive completed tasks at the end of the session." | **Violation.** Archive completed tasks immediately after they are finished to keep active context clean and token-efficient. |
<!-- markdownlint-enable MD013 -->
