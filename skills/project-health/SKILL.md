---
name: project-health
description: Use when validating project correctness, completeness, and consistency, or when /project-health is explicitly invoked.
---

# Project Health

Answer the question: **is this project correct, complete, and consistent?**

Runs universal quality checks that apply to every project type, then automatically
chains to the type-specific health skill based on the detected project type.

---

## Step 0 — Detect Project Type

Auto-detect the project type by checking for characteristic files:

| Type | Detection |
| --- | --- |
| `java` | `pom.xml` or `build.gradle` at root |
| `python` | `pyproject.toml`, `setup.py`, or `requirements.txt` at root |
| `ts` | `package.json` at root (with `tsconfig.json` or `.ts` in `src/`) |
| `generic` | None of the above |

If detection is ambiguous (e.g. both `package.json` and `pyproject.toml` found),
inform the user and ask which one to use. If nothing is detected, ask the user
to declare it.

Store the type — type-aware checks use it throughout this skill.

---

## Step 1 — Determine Mode

Two modes:

| Flag | What runs |
| --- | --- |
| `--commit` | Quick validators only, then stop |
| *(no flag)* | All universal checks + type-specific checks (tier 3+) |

If `--commit` was passed, run automated validators and exit. Otherwise run
everything.

Also parse:

- `--save` → write report to `YYYY-MM-DD-health-report.md` after output
- Category names (e.g. `docs-sync config`) → run only those categories, at full depth

---

## Step 2 — Commit Mode

In `--commit` mode, run quick automated validators only (linters, schema checks, etc.),
report findings, then **stop here**. Present findings and exit.

---

## Step 3 — Build Document Scan List

In full mode, build the scan scope before running checks.

**Always included:**

- All `.md` files (recursive) under `doc/`, `docs/`, `documentation/` (case-insensitive)
- Root-level `.md` files matching: `readme`, `overview`, `summary`, `index`, `contributing`,
  `governance`, `code_of_conduct`, `changelog`, `history`, `release`, `architecture`,
  `design`, `decisions`, `vision`, `philosophy`, `principles`, `api`, `schema`, `glossary`,
  `security`, `deployment`, `install`, `usage`, `troubleshooting`, `roadmap`, `spec`,
  `requirements`, `quality` (case-insensitive match on filename stem)
- Any root `.md` not on the list is still scanned
- Any `README.md`, `CHANGELOG.md`, `CONTRIBUTING.md`, `ARCHITECTURE.md` anywhere in the tree

**Type-specific additions (use detected type from Step 0):**

- `skills` → all `SKILL.md` files in direct subdirectories
- `java` → `pom.xml`, `build.gradle`, Javadoc comments in `src/`

---

## Step 4 — Run Universal Checks

Run all applicable check categories. If specific categories were requested, run only those.

**Read [check-categories.md](check-categories.md)** for the full quality
checklists for all 9 universal check categories before running checks:
`docs-sync`, `logic`, `config`, `security`, `release`,
`user-journey`, `git`, `artifacts`, `framework`.

---

## Step 5 — Chain to Type-Specific Skill

After universal checks complete, automatically invoke the type-specific health
skill in the same session:

| Project type | Invoke |
| --- | --- |
| `skills` | `skills-project-health` |
| `java` | `java-project-health` |
| `generic` | Skip — universal checks only |

The type-specific skill's output is appended to the report. Do NOT redirect the
user to run a separate command — chain automatically.

If the type-specific skill does not exist yet, note it as a LOW finding:
> [config] Type-specific health skill `{type}-project-health` not yet available

---

## Step 6 — Present Report

```
## project-health report — [categories run] — tier [N]

### CRITICAL (must fix)
- [category] finding description

### HIGH (should fix)
- [category] finding description

### MEDIUM (worth fixing)
- [category] finding description

### LOW (nice to fix)
- [category] finding description

### PASS
✅ category1, category2, ...
```

Universal findings have no extra prefix. Type-specific findings are prefixed
with `[type]` (e.g. `[java]`). If no findings in a severity level, omit that section.

---

## Step 7 — Offer Auto-Fix (Mechanical Issues Only)

For mechanical findings (wrong count in README, stale version number, missing
`commands/<name>.md`), offer:

> **Auto-fixable findings detected.**
>
> Would you like me to apply mechanical fixes now?
>
> - [list of specific fixes]
>
> **(YES / NO — judgment calls are never auto-applied)**

Wait for response. Apply only on YES. Never auto-apply.

---

## Step 8 — Save Report (if --save)

If `--save` was passed, write findings to a date-prefixed file:

```bash
# Format: YYYY-MM-DD-health-report.md
```

Tell user:
> Report saved to `YYYY-MM-DD-health-report.md`. This file is gitignored by default.

Verify `.gitignore` includes `*-health-report.md` or similar. If not, suggest adding it.

---

## Common Pitfalls

| Mistake | Why It's Wrong | Fix |
| --------- | ---------------- | ----- |
| Running type-specific checks before reading project type | Checks have no context | Always read project type in Step 0 first |
| Reporting "plans to implement" as bugs | Intentional design language | Distinguish docs describing intent vs. describing current state |
| Auto-fixing judgment findings | Judgment calls require human decision | Only auto-fix mechanical findings, always with YES confirmation |
| Skipping chain to type-specific skill | Incomplete health picture | Chain is mandatory unless type is generic |
| Treating all findings as equal | CRITICAL blocks release, LOW does not | Use severity consistently |
| Running `docs-sync` without reading the actual source files | Can't verify claims without reading | Read code and docs together |

---

## Success Criteria

Health check is complete when:

- ✅ Project type detected before any checks ran
- ✅ Mode confirmed (`--commit` or full)
- ✅ All applicable universal categories checked
- ✅ Type-specific skill chained (or skipped for generic)
- ✅ Report presented with findings grouped by severity
- ✅ Mechanical auto-fix offered (not applied without YES)
- ✅ Report saved if `--save` was passed

**Not complete until** all applicable categories checked and report presented.

---

## Skill Chaining

**Invoked by:**

- User invokes `/project-health`
- Type-specific health skills invoke this as their prerequisite foundation

**Suggests (not auto-chained):**

- `python-project-health` — suggested when a Python project is detected
- `ts-project-health` — suggested when a TypeScript/Node.js project is detected
- `java-project-health` — suggested when a TypeScript/Node.js project is detected
