# project-health — Universal Check Categories

Checklists for all 9 universal check categories.
Referenced by `project-health/SKILL.md` Step 5 and all type-specific sub-skills.

---

## Modes

| Mode | What runs |
|------|-----------|
| `--commit` | Quick validators only, then stop |
| *(no flag)* | All universal checks + type-specific checks |

---

## Severity Scale

| Severity | Meaning |
|----------|---------|
| **CRITICAL** | Correctness failure — should block release |
| **HIGH** | Should fix before shipping |
| **MEDIUM** | Worth fixing in next session |
| **LOW** | Nice to fix, low urgency |

---

### `docs-sync` — Documentation Accuracy

- [ ] No "planned" / "not yet implemented" language for things that already exist
- [ ] Version numbers consistent across all places they appear
- [ ] URLs and external references are correct and reachable
- [ ] No stale "TODO" or "coming soon" references for completed work
- [ ] Release status language matches actual state
- [ ] Temporal claims still accurate ("as of Q1 2026", "in the last 6 months")
- [ ] Environment variable names, config keys, and file paths match actual source

---

### `logic` — Workflow Logic & UX

- [ ] No workflow step references a script or file that doesn't exist
- [ ] No workflow step requires an external tool without verifying it's installed
- [ ] Hook outputs are directive (ACTION REQUIRED) not just informational
- [ ] Hook doesn't fire on non-git directories
- [ ] No workflow blocks progress without giving the user a way forward
- [ ] Error messages include recovery steps
- [ ] No redundant checks (same thing checked twice in same flow)
- [ ] Chained workflows don't create infinite loops
- [ ] Exit codes consistent and documented
- [ ] Workflow steps requiring user judgment specify clear decision criteria
- [ ] Ordered workflows document what happens if steps are skipped or reordered

---

### `config` — Project Configuration Health

- [ ] Every config file parses without syntax errors (invalid TOML, broken YAML, unclosed JSON)

---

### `security` — Security & Safety

- [ ] No hardcoded tokens, passwords, or API keys in any file
- [ ] Shell scripts quote variables to prevent word splitting
- [ ] No `rm -rf` without explicit path validation
- [ ] No `eval` of untrusted input
- [ ] All executable scripts have correct permissions
- [ ] No secrets in git history (check recent commits)
- [ ] Scripts validate inputs before acting on them
- [ ] External tool dependencies documented or checked before use
- [ ] Scripts that write files check the target directory exists first
- [ ] Relative paths in scripts work regardless of calling directory

---

### `release` — Release Readiness

- [ ] No SNAPSHOT, dev, or alpha markers in release artifacts
- [ ] All component versions consistent and intentional
- [ ] GitHub labels set up for release note generation
- [ ] `gh release create --generate-notes` would produce meaningful output
- [ ] No obviously incomplete components (stubs, placeholders, empty sections)
- [ ] All tests passing
- [ ] Release notes reference the issues or PRs being released
- [ ] Release notes don't claim fixes/features not in this release

---

### `user-journey` — End User Experience

- [ ] Getting started path documented and works end-to-end
- [ ] First meaningful action is guided
- [ ] Setup prompts are clear and skippable where optional
- [ ] Error messages explain what went wrong and how to recover
- [ ] No dead ends — every failure state has a next step
- [ ] Entry points (commands, slash commands, scripts) work as documented
- [ ] Getting started validated from a fresh environment
- [ ] Documented error recovery steps actually resolve the stated error

---

### `git` — Repository State

- [ ] No uncommitted changes that should be committed
- [ ] No stale git worktrees
- [ ] Tags consistent with marketplace/package versions (for release)
- [ ] No merge conflict markers in tracked files
- [ ] Branch is up to date with remote

---

### `artifacts` — Required Artifacts Exist

(apply to the correct artifact list for the detected project type)
- [ ] All required files and directories exist
- [ ] Required configuration files are present and parseable
- [ ] Any referenced artifact actually exists at its declared path
- [ ] No required artifact is empty, stubbed, or placeholder only
- [ ] No required artifact appears abandoned

---

### `framework` — Framework Pattern Correctness

- [ ] Code examples use patterns correct for the declared framework
- [ ] Documented workflows account for framework-specific constraints
- [ ] No guidance recommends an approach the framework explicitly discourages
- [ ] Framework-specific annotations, decorators, or conventions used correctly in examples
- [ ] Patterns are current for the version of the framework in use
