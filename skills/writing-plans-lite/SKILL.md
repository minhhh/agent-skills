---
name: writing-plans-lite
description: Use when initiating or maintaining a lightweight plan in a single markdown PRD file
---

# Writing Plans Lite

A workflow for writing and maintaining requirements in a single markdown file
(`prd-mvp.md` or `prd-[feature].md`).

## When to Use

- Use when:
  - The user explicitly requests a lightweight planning structure, or activates
    it with keywords like "lite plan", "plan lite", or "writing-plans-lite".
  - You proactively determine that a new feature or task warrants planning, but
    the standard planning mode (implementation plan) is too heavy.
  - The agent completes a task in a PRD and needs to update and archive it.

- When NOT to use:
  - Do NOT use this skill to write implementation code, modify source files (except the PRD file itself), run tests, run builds, or execute tasks. This skill is strictly for documentation, requirements planning, and plan maintenance.
  - For executing tasks and writing code, load and transition to the appropriate execution or development skill (e.g., `executing-plans-lite` or language-specific development skills).

## Template Reference

The base template is located at [prd-template.md](assets/prd-template.md).

## Continuous Workflow

1. **Bootstrap & Plan**: Create `.tasks/prd-mvp.md` (or `prd-[feature].md` for
   sub-features) using the template.
   - **CRITICAL ORDER**: You MUST create/bootstrap this PRD file *before*
     writing any implementation code or making edits to source files. No
     exceptions.
   - **Upfront Specification**: You MUST write the entire specifications, requirements, and a detailed breakdown of all tasks in Section 3 *before* starting any execution. The PRD must be completely self-contained. It MUST contain all technical specifications, templates, and exact requirements. Do NOT write partial breakdowns or refer to instructions, context, or agreements in the chat history (e.g. "implement the templates discussed in chat").
   - **Assets Folder**: For complex features or features that involve reference templates, config files, or mockups, you MUST create a dedicated folder at `.tasks/prd-[feature]/` (e.g. `.tasks/prd-kb-to-md/`) to store these resources and link to them within the PRD. All links to these assets and other local repository files MUST be relative paths (e.g., `[label](prd-[feature]/file.md)`), never absolute `file:///` URLs.
2. **Obtain Approval**: You MUST present the completed PRD file to the user and **wait for their explicit approval** before transitioning to implementation. Once approved, stop this planning workflow and transition to the appropriate task execution/development skill (e.g., `executing-plans-lite`). Do not write or modify source code under this skill.
3. **Progress Execution**: Update checklists as you complete work.
4. **Dashboard Updates**: Update the `Active Dashboard` using status symbols:
   - `[ ]` Todo
   - `[/]` In-Progress
   - `[x]` Completed
5. **Archiving (Critical)**: When a task is completed, immediately move its
   entry and details from Sections 2 & 3 to **Section 5 (Archive)**.

## Strict Skill Boundaries (Enforced)

This skill is **strictly** for planning and tracking (creating/updating the PRD). It has **zero** implementation authority.

**CRITICAL RULE: DO NOT WRITE CODE OR EXECUTE TASKS WHILE THIS SKILL IS ACTIVE.**

If you are asked to implement, execute, or write code:
1. You MUST first complete or update the PRD.
2. You MUST state: "PRD is ready/updated. Transitioning to [executing-plans-lite] (or the relevant development skill) to begin execution."
3. Do NOT make any code modifications or execute commands in the same turn/thought block where you are planning. Transitioning is a hard boundary.

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
| "The user is in a hurry, so I can start coding before detailing the tasks or getting approval." | **Violation.** You must fully detail all tasks and wait for explicit user approval before execution. |
| "I can add task details and breakdowns as I go." | **Violation.** Full specifications and breakdowns must be written upfront in Section 3 before starting. |
| "Since I have this planning skill active, I can also implement the code change to save time." | **Violation.** This skill is for planning only. Transition to an execution skill (e.g., `executing-plans-lite`) to write code. |
| "I can refer to details from the chat/LLM prompt in the PRD details instead of copying them." | **Violation.** The PRD must be completely self-contained. Always copy all templates and specs into the PRD or its assets folder. |
| "The user requested absolute links so they are clickable in their editor." | **Violation.** Always use relative paths for local references in the PRD to maintain repository portability. |
| "The user requested `tasks/` instead of `.tasks/` to avoid hidden files, or the project already has a `tasks/` folder." | **Violation.** You MUST store all plans and assets in `.tasks/`. Do not use `tasks/` (without the leading dot) under any circumstances, even if requested by the user. |
<!-- markdownlint-enable MD013 -->

## Red Flags - STOP and Start Over

- **Making code changes to source files** (other than the PRD itself) or executing tasks while this skill is active.
- **Running compile, test, build, or deploy commands** for the project under the scope of this skill.
- **Combining planning and coding in a single turn**. You must separate them by a hard transition.
- **Bypassing the transition to execution/development skills** once the plan is approved.
- **Referencing chat history, LLM context, or user instructions** inside the PRD checklists/details instead of detailing them explicitly.
- **Storing project-related planning assets outside the `.tasks/prd-[feature]/` directory.**
- **Storing the PRD file or assets in a `tasks/` directory (without the leading dot) instead of `.tasks/`.**
- **Using absolute links (like `file:///` URLs) in the PRD** for local file/asset references.
