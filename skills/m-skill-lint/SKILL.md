---
name: m-skill-lint
description: >
  Validate skill directory structure and conventions. Checks for missing files,
  correct frontmatter, proper attribution for third-party skills. Also guides
  the workflow for creating new skills.
---

# M-Skill-Lint: Skill Structure Validator

Validates that skills in `skills/` follow the repository conventions. Run this
when adding a new skill, after bulk changes, or during PR review.

## Workflow: Creating a New Skill

When creating a new skill, follow these steps in order. The validation checks
below confirm each step is complete.

### Step 1 — Create directory

```text
skills/<skill-name>/
```

### Step 2 — Write SKILL.md

- YAML frontmatter with exactly `name` (kebab-case, matches directory) and `description`.
- Body follows the [Antigravity skill format](https://github.com/google-deepmind/antigravity).
- If the skill builds on another skill (e.g. `python-code-review` builds on `code-review-principles`), declare a `## Prerequisites` section.

### Step 3 — Add ATTRIBUTION + LICENSE (third-party only)

- Read the original `SKILL.md` from the source repository.
- Create `ATTRIBUTION` with the standard template (see below).
- Copy the original `LICENSE` file verbatim.

**ATTRIBUTION format:**
```markdown
# Attribution

This skill is derived from the <source> collection of coding conventions and best practices.

## Original Source

- **Repository:** https://github.com/<org>/<repo>
- **Author:** <name>
- **License:** <license-type>
```

### Step 4 — Add README.md entry

- If first-party: add under `### First-Party Skills`, alphabetically sorted.
- If third-party: add under `### Third-Party Skills`, alphabetically sorted, with ` — from [source](url)` suffix.

### Step 5 — Register (optional)

- Add to `scripts/config/default-skills.txt` if it should be linked by default for all agents.

### Step 6 — Verify

Run the checks below to confirm everything is correct.

## Validation Checks

### File Existence

| Check | Rule |
|---|---|
| `SKILL.md` | Required for every skill |
| `ATTRIBUTION` | Required if skill is third-party; forbidden if first-party |
| `LICENSE` | Required if third-party; forbidden if first-party |
| `scripts/` | Optional — only present if the skill ships helper binaries |
| `references/` | Optional — only present for reference docs linked from SKILL.md |

### SKILL.md Frontmatter

- Must open and close with `---`
- Must contain `name:` — value must match the parent directory name
- Must contain `description:` — single-line or YAML folded block scalar (`>`)
- Must not contain extra fields (`version`, `author`, `tags`, etc.)
- The rest of the file is free-form Markdown

### ATTRIBUTION Format

- Must follow the standard template (see Step 3 above)
- Repository URL must be valid (reachable)
- License must match the LICENSE file in the same directory

### README.md Entry

- Must be in the correct section (first-party or third-party)
- Must match the skill's actual classification
- Third-party entries must include the `— from [name](url)` suffix
- Must be alphabetically sorted within its section

### Default Skills

- Skills in `scripts/config/default-skills.txt` must actually exist as directories in `skills/`
- No duplicate entries
- No blank lines or inline comments (lines starting with `#` are allowed)

## Lint Command

```bash
# Run all checks from repo root
./skills/m-skill-lint/scripts/lint
```
