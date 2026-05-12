# Agent Skills

A collection of specialized skills for AI coding agents (like Antigravity, Claude, or GitHub Copilot) to extend their capabilities with custom tools and workflows.

## Repository Structure

```
agent-skills/
├── skills/
│   └── <skill-name>/
│       ├── SKILL.md        # Documentation and instructions for the agent
│       └── scripts/        # Helper scripts and tools for the skill
├── scripts/
│   ├── link-default.sh    # Link all default skills to agent config paths
│   ├── link-skill.sh      # Link a specific skill
│   └── unlink-skill.sh    # Remove a specific skill link
├── Makefile               # Convenient commands for installation
└── README.md
```

## Available Skills

### [caveman](skills/caveman/SKILL.md)
Terse communication mode for AI responses that strips filler words to reduce token usage. Three intensity levels: `lite`, `full`, and `ultra`. Automatically reverts to normal prose for security warnings and destructive operations.

### [git-conventionalize](skills/git-conventionalize/SKILL.md)
Non-interactively rewrite git commit messages to follow the [Conventional Commits](https://www.conventionalcommits.org/) standard.

### [onboard](skills/onboard/SKILL.md)
Structured 12-step process for rapidly orienting to a new codebase. Parallel discovery of tech stack, structure, entry points, dependencies, and config, culminating in a comprehensive `project_map.md`.

### [summarize-markdown-book](skills/summarize-markdown-book/SKILL.md)
Summarize long-form technical documentation, manuals, or books into a structured hierarchical Markdown format with custom markers (`▼`, `###`, `*`).

## Gemini Extension

### [gemini-superpowers](gemini-superpowers/README.md)
A Gemini CLI extension that ports the [Superpowers](https://github.com/obra/superpowers) project. Bundles skills for planning, TDD, systematic debugging, code review, and subagent-driven development, along with CLI commands and session hooks.

## Installation

### 1. Link Skills to Agents

Use `make` commands to symlink skills to your AI agents' configuration directories. This makes the skills available across different agent interfaces.

```bash
# Link all default skills (listed in scripts/config/default-skills.txt)
make link-default

# Link a specific skill
make link-skill SKILL=git-conventionalize

# Unlink a specific skill
make unlink-skill SKILL=git-conventionalize
```

The scripts currently link to:
- `~/.agents/skills/`
- `~/.claude/skills/`

### 2. Skill-Specific Setup

Some skills require additional binaries or dependencies.

#### Git Conventionalize
The `git-reword` tool is required for this skill. Build and install it:

```bash
go install ./skills/git-conventionalize/scripts/git-reword
```

## Creating New Skills

1. Create a new directory in `skills/`.
2. Add a `SKILL.md` following the [Antigravity skill format](https://github.com/google-deepmind/antigravity).
3. Add any necessary helper scripts in a `scripts/` subdirectory.
4. Add the skill name to `scripts/config/default-skills.txt` if it should be linked by default.


