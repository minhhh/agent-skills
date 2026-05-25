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
## Private Skills

To maintain personal or proprietary skills separate from the public repository, this project uses a Git submodule pointing to a private repository (`private-agent-skills`):

```
private-agent-skills/
└── skills/
    └── <private-skill-name>/
        ├── SKILL.md
        └── scripts/
```

### Why Use a Private Skills Submodule?
- **Separation of Concerns**: Keeps personal coding workflows, domain-specific instructions, or proprietary enterprise automation separate from general-purpose, open-source skills.
- **Security & Privacy**: Avoids accidental leakage of internal credentials, private code patterns, or proprietary business logic to public repositories.
- **Unified Interface**: Allows local tooling to manage and symlink both public and private skills using the same Makefile target.

## Available Skills

### First-Party Skills

#### [deconstruct](skills/deconstruct/SKILL.md)
Incrementally decompose a codebase into a detailed reference document — from architecture overview down to specific code flows. Runs parallel discovery phases (entry points, core symbols, control flow, feature grouping, side effects, testing patterns) and writes findings to a structured `deconstructed.md`.

#### [extend-markdown-summary](skills/extend-markdown-summary/SKILL.md)
Surgically enrich existing technical summaries with deeper detail — missing commands, configuration flags, performance tuning parameters, and architectural rationale — while preserving the established structural hierarchy.

#### [general-code-review](skills/general-code-review/SKILL.md)
Conduct a comprehensive automated code review covering logical correctness, security vulnerabilities, performance bottlenecks, and architectural alignment. Skips syntax/formatting nits and provides ranked critical issues with surgical fix snippets.

#### [git-conventionalize](skills/git-conventionalize/SKILL.md)
Non-interactively rewrite git commit messages to follow the [Conventional Commits](https://www.conventionalcommits.org/) standard.

#### [golang-code-review](skills/golang-code-review/SKILL.md)
Go-specific code review skill that builds on `code-review-principles`. Focuses on safety violations (nil maps, error discarding, missing defer), error handling (unwrapped errors, log-and-return, panic misuse), concurrency (goroutine leaks, loop variable capture, context in struct), and code quality. Provides a severity decision flow and checklist for reviewing Go code.

#### [golang-dev](skills/golang-dev/SKILL.md)
Comprehensive Go development skill for writing new code, fixing bugs, refactoring, or adding tests. Covers safety (nil maps, error checking, context propagation), error handling (wrapping, Is/As, single handling rule), concurrency (goroutine lifecycle, errgroup, channel discipline), testing (table-driven, golden files, race detector), code quality (formatting, interfaces, packaging), and performance (preallocation, profiling, strings.Builder). Includes a priority decision flow: Safety > Error Handling > Concurrency > Code Quality.

#### [golang-performance](skills/golang-performance/SKILL.md)
Profile and optimize Go code using pprof, benchstat, and fgprof. Covers CPU profiling, memory profiling, GC tuning, and optimization patterns (slice preallocation, strings.Builder, struct alignment, sync.Pool, GOMEMLIMIT, HTTP transport tuning).

#### [m-skill-lint](skills/m-skill-lint/SKILL.md)
Validate skill directory structure and conventions. Checks for missing files, correct frontmatter, proper attribution for third-party skills, and README consistency. Includes the full workflow for creating new skills.

#### [onboard](skills/onboard/SKILL.md)
Structured 12-step process for rapidly orienting to a new codebase. Parallel discovery of tech stack, structure, entry points, dependencies, and config, culminating in a comprehensive `project_map.md`.

#### [summarize-to-markdown](skills/summarize-to-markdown/SKILL.md)
Summarize long-form technical documentation, manuals, or books into a structured hierarchical Markdown format with custom markers (`▼`, `###`, `*`).

### Third-Party Skills

#### [caveman](skills/caveman/SKILL.md) — from [juliusbrussee/caveman](https://github.com/juliusbrussee/caveman)
Terse communication mode for AI responses that strips filler words to reduce token usage. Three intensity levels: `lite`, `full`, and `ultra`. Automatically reverts to normal prose for security warnings and destructive operations.

#### [code-review-principles](skills/code-review-principles/SKILL.md) — from [mdproctor/cc-praxis](https://github.com/mdproctor/cc-praxis)
Foundation skill providing universal code review principles for catching critical issues — safety violations, concurrency bugs, and silent data corruption. Not invoked directly; loaded as a prerequisite by language-specific code review skills (e.g., `python-code-review`, `java-code-review`).

#### [golang-testing](skills/golang-testing/SKILL.md) — from [affaan-m/ECC](https://github.com/affaan-m/ECC/)
Go testing patterns: table-driven tests, subtests, benchmarks, fuzzing, and test coverage. Follows idiomatic Go testing practices.

#### [project-health](skills/project-health/SKILL.md) — from [mdproctor/cc-praxis](https://github.com/mdproctor/cc-praxis)
Answer whether a project is correct, complete, and consistent. Auto-detects project type (Python, Java, TypeScript, skills, or generic), then runs universal quality checks across 9 categories — docs sync, logic, config, security, release, user journey, git, artifacts, and framework. Supports `--commit` mode for quick pre-commit validators and chains automatically to type-specific health skills.

#### [python-code-review](skills/python-code-review/SKILL.md) — from [mdproctor/cc-praxis](https://github.com/mdproctor/cc-praxis)
Python-specific code review skill that builds on `code-review-principles`. Focuses on safety violations (mutable defaults, bare except), type safety, async correctness, and testing patterns. Provides a severity decision flow and checklist for reviewing Python code.

#### [python-dev](skills/python-dev/SKILL.md) — from [mdproctor/cc-praxis](https://github.com/mdproctor/cc-praxis)
Python development skill for writing new code, fixing bugs, refactoring, or adding tests. Covers safety rules (context managers, no eval), type safety (type hints, mypy --strict), async patterns, testing best practices (pytest fixtures, parametrized tests), and code quality. Includes a priority decision flow: Safety > Type Safety > Async Correctness > Code Quality.

#### [python-optimization](skills/python-optimization/SKILL.md) — from [mdproctor/cc-praxis](https://github.com/mdproctor/cc-praxis)
Profile and optimize Python code using cProfile, memory profilers, and performance best practices. Covers CPU profiling, memory optimization, line profiling, production profiling with py-spy, and optimization patterns including data structures, comprehensions, generators, and local variable access.

#### [python-project-health](skills/python-project-health/SKILL.md) — from [mdproctor/cc-praxis](https://github.com/mdproctor/cc-praxis)
Python-specific health extension that builds on `project-health`. Adds checks for type safety (mypy), dependency health (lockfile integrity, vulnerability scanning), code quality (ruff/flake8, anti-patterns), test health (coverage, test hygiene), and build integrity (import safety, packaging). Findings are prefixed with `[python]` in the combined report.

#### [python-security-audit](skills/python-security-audit/SKILL.md) — from [mdproctor/cc-praxis](https://github.com/mdproctor/cc-praxis)
Security audit skill for Python server applications, following OWASP Top 10 (2021). Covers injection, broken access control, cryptographic failures, insecure design, security misconfiguration, vulnerable components, identification/auth failures, software/data integrity failures, security logging, and SSRF. Builds on `security-audit-principles` and provides Python-specific examples and remediation patterns.

#### [security-audit-principles](skills/security-audit-principles/SKILL.md) — from [mdproctor/cc-praxis](https://github.com/mdproctor/cc-praxis)
Foundation skill providing universal security audit principles. Covers the core security mindset — authentication, authorization, input validation, cryptography, session management, and secure configuration. Not invoked directly; loaded as a prerequisite by language-specific security audit skills (e.g., `python-security-audit`).

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


