---
name: git-conventionalize
description: Use when you need to non-interactively rewrite existing commit messages already in the git history to follow the Conventional Commit standard.
---

# Git Conventionalize Skill

Non-interactively rewrite git commit messages using the Conventional Commit standard.

## Tools

```bash
git-reword status               # Check if repo is clean for rebase
git-reword analyze <hash>       # Get commit diff and message for analysis
git-reword lint "<message>"     # Validate message format
git-reword apply <hash> "<msg>"          # Execute the rewrite
git-reword apply <hash> "<msg>" --author "Name <email>" # Rewrite msg and author
```

## When to Use

- User asks to "conventionalize" or reword an existing commit message already in the git history

## When NOT to Use

- Do NOT use when generating, writing, or drafting a new commit message for staged or unstaged changes
- Do NOT use when creating a commit message for the first time

## Conventional Commit Format

```
<type>(<scope>): <description>

Types: feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert
Scope: optional, typically the affected module or file
```

## Workflow

1. **Check Readiness**: Run `git-reword status` to ensure no uncommitted changes.
2. **Gather Context**: Run `git-reword analyze <hash>` to see the current message and the actual code changes.
3. **Generate Message**: Draft a conventional message based on the analysis.
4. **Validate**: (Optional) Run `git-reword lint "<new-message>"` to ensure your draft is correct.
5. **Execute**: Run `git-reword apply <hash> "<new-message>"` to perform the rebase.

## Multiple Commits

For multiple commits, process them one by one. Note that hashes will change for all commits *after* the one you just rewrote if they are descendants. It is often better to start from the oldest commit and work forward.

## Notes

- The tool uses `git rebase -i` under the hood.
- Rewriting history changes commit hashes; warn the user if they have already pushed.
- Original commit hash changes after rebase.
- If conflicts occur during rebase, the tool will stop and require manual intervention.
