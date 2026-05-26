---
name: m-skill-lint
description: Use when requested by the user to validate skill directory structure and conventions
---

# M-Skill-Lint: Skill Structure Validator

Validates that skills in `skills/` follow the repository conventions. Run this ONLY when explicitly requested by the user, not when simply modifying a skill or adding a new skill.
 DO NOT run this automatically when modifying or adding a new skill unless explicitly requested

## Validation Checks

### File Existence

| Check | Rule |
| --- | --- |
| `SKILL.md` | Required for every skill |
| `ATTRIBUTION` | Required if skill is third-party; forbidden if first-party |
| `LICENSE` | Required if third-party; forbidden if first-party |
| `scripts/` | Optional — only present if the skill ships helper binaries |
| `references/` | Optional — only present for reference docs linked from SKILL.md |

### SKILL.md Frontmatter

- Must open and close with `---`
- Must contain exactly two fields: `name` and `description` (no extra fields like `version`, `author`, `tags`, etc. allowed)
- Must be <= 1024 characters total length
- The `name` field value must match the parent directory name, using only letters, numbers, and hyphens (`^[a-zA-Z0-9-]+$`)
- The `description` field value must start with "Use when..." (case-insensitive)

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
./.agents/skills/m-skill-lint/scripts/lint
```

## Common Pitfalls — These Thoughts Mean STOP

- "The linter reported a failing skill directory, let me write a script or modify the files to fix it."
- "I should create ATTRIBUTION or LICENSE files for skills that fail validation."
- "Let's automate fixing the directory structure/conventions of other folders."

**All of these mean: STOP. Do NOT try to modify, add, or fix any skill directories or files based on the validation results of the lint script unless the user explicitly requests you to do so. The purpose of the linter is purely to report validation results, not to automatically resolve them.**
