# kb-to-md Skill Creation

---

## 1. Product Specification

*Core vision, requirements, and user needs.*

### Feature Overview

Create a new agent skill named `kb-to-md` that instructs agents how to format synthesized knowledge into a structured markdown layout matching specific use cases, and how to insert/copy those sections into specified target files.

### User Stories / Requirements

- **task-01**: Create `skills/kb-to-md/SKILL.md` containing the markdown formatting rules, examples, and file insertion guidelines.
- **task-02**: Register the new skill in the repository root `README.md` in alphabetical order.
- **task-03**: Refine the `kb-to-md` skill specification.
- **task-04**: Run the skill linter (`./.agents/skills/m-skill-lint/scripts/lint`) to ensure the directory structure and metadata conform to the repository rules.
- **task-05**: Run subagent pressure tests to verify compliance with the new skill under TDD.

---

## 2. Active Dashboard

*Current tracker.*

- [x] **task-03**: Refine the `kb-to-md` skill specification
- [x] **task-04**: Run skill linter validation
- [x] **task-05**: Perform TDD verification with a subagent

---

## 3. Active Task Details

*Only contains details/checklists for active tasks in Section 2.*

---

## 4. Future Roadmap & Backlog

*Placeholders for future tasks.*

- [ ] **task-06**: Extend skill if additional insertion modes are needed.

---

## 5. History / Archive

*When a task is completed, cut-and-paste both its Dashboard status and its Detailed Requirements/Checklist here.*

- [x] **task-01**: Create `skills/kb-to-md/SKILL.md`
  - [x] Create directory `skills/kb-to-md`
  - [x] Write `SKILL.md` with Yaml frontmatter and the 5 formatting templates (Case 1 to 5) as specified in [golden-templates.md](prd-kb-to-md/golden-templates.md)
  - [x] Add spacing and list styling rules (bold concepts, exactly 4-space nested list indentation, and blank line spacing) as defined in [golden-templates.md](prd-kb-to-md/golden-templates.md)
  - [x] Add file insertion and merging guidelines defined in [golden-templates.md](prd-kb-to-md/golden-templates.md)
- [x] **task-03**: Refine the `kb-to-md` skill specification
  - [x] Gather feedback and instructions on the required refinements from the user
  - [x] Wrap all template structures for Cases 1 to 4 in ```markdown code blocks inside `skills/kb-to-md/SKILL.md`
  - [x] Update skill with named layout styles (bullets, subsection, chapter-subsection, chapter) and Option A natural language triggers
  - [x] Move surgical modification case to core workflow guidelines
- [x] **task-04**: Run skill linter validation
  - [x] Execute `./.agents/skills/m-skill-lint/scripts/lint`
  - [x] Verify that all checks pass cleanly
- [x] **task-05**: Perform TDD verification with a subagent
  - [x] Create dummy target markdown file (`scratch/test_readme.md`)
  - [x] Run a baseline test with a subagent (WITHOUT the skill loaded) and observe formatting failures
  - [x] Run a verification test with a subagent (WITH the skill loaded) using natural language triggers (Option A)
  - [x] Clean up temporary dummy files

### task-01: Create `skills/kb-to-md/SKILL.md`

*Objective: Create the main skill file defining the 5 formatting and insertion use cases.*

- **Checklist**:
  - [x] Create directory `skills/kb-to-md`
  - [x] Write `SKILL.md` with Yaml frontmatter and the 5 formatting templates (Case 1 to 5) as specified in [golden-templates.md](prd-kb-to-md/golden-templates.md)
  - [x] Add spacing and list styling rules (bold concepts, exactly 4-space nested list indentation, and blank line spacing) as defined in [golden-templates.md](prd-kb-to-md/golden-templates.md)
  - [x] Add file insertion and merging guidelines defined in [golden-templates.md](prd-kb-to-md/golden-templates.md)

### task-03: Refine the `kb-to-md` skill specification

*Objective: Refine the skill's formatting and layout templates based on user feedback.*

- **Checklist**:
  - [x] Gather feedback and instructions on the required refinements from the user
  - [x] Wrap all template structures for Cases 1 to 4 in ```markdown code blocks inside `skills/kb-to-md/SKILL.md`
  - [x] Update skill with named layout styles (bullets, subsection, chapter-subsection, chapter) and Option A natural language triggers
  - [x] Move surgical modification case to core workflow guidelines

### task-04: Run skill linter validation

*Objective: Ensure the new skill structure passes the repository structure checks.*

- **Checklist**:
  - [x] Execute `./.agents/skills/m-skill-lint/scripts/lint`
  - [x] Verify that all checks pass cleanly

### task-05: Perform TDD verification with a subagent

*Objective: Verify that the skill successfully dictates agent formatting under pressure.*

- **Checklist**:
  - [x] Create dummy target markdown file (`scratch/test_readme.md`)
  - [x] Run a baseline test with a subagent (WITHOUT the skill loaded) and observe formatting failures
  - [x] Run a verification test with a subagent (WITH the skill loaded) using natural language triggers (Option A)
  - [x] Clean up temporary dummy files
