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

#### Codebase Orientation & Onboarding

##### [deconstruct](skills/deconstruct/SKILL.md)

Incrementally decompose a codebase into a detailed reference document — from architecture overview down to specific code flows. Runs parallel discovery phases (entry points, core symbols, control flow, feature grouping, side effects, testing patterns) and writes findings to a structured `deconstructed.md`.

##### [onboard](skills/onboard/SKILL.md)

Structured 12-step process for rapidly orienting to a new codebase. Parallel discovery of tech stack, structure, entry points, dependencies, and config, culminating in a comprehensive `project_map.md`.

#### Planning & Task Execution (Lite)

##### [executing-plans-lite](skills/executing-plans-lite/SKILL.md)

Lightweight workflow for executing tasks defined in a single markdown PRD/plan, keeping progress updated, and archiving completed tasks immediately.

##### [writing-plans-lite](skills/writing-plans-lite/SKILL.md)

Flexible, lightweight task and specification tracker in a unified markdown file for starting new projects (greenfield) or adding features (brownfield).

#### Development & Workflow Automation

##### [git-conventionalize](skills/git-conventionalize/SKILL.md)

Non-interactively rewrite git commit messages to follow the [Conventional Commits](https://www.conventionalcommits.org/) standard.

##### [golang-dev](skills/golang-dev/SKILL.md)

Comprehensive Go development skill for writing new code, fixing bugs, refactoring, or adding tests. Covers safety (nil maps, error checking, context propagation), error handling (wrapping, Is/As, single handling rule), concurrency (goroutine lifecycle, errgroup, channel discipline), testing (table-driven, golden files, race detector), code quality (formatting, interfaces, packaging), and performance (preallocation, profiling, strings.Builder). Includes a priority decision flow: Safety > Error Handling > Concurrency > Code Quality.

##### [java-dev](skills/java-dev/SKILL.md)

Comprehensive Java development skill for writing new code, fixing bugs, refactoring, or adding tests. Applies to Spring Boot, Quarkus, and plain Java projects. Covers safety (resource leaks, classloader leaks, silent corruption), concurrency (thread models, read-modify-write, optimistic locking, executors), performance (hot paths, N+1 prevention, batch queries), testing (JUnit 5, AssertJ, framework-specific options), code quality (naming, final, imports, text blocks), exception handling, logging, input validation, DTO conventions, and common pitfalls. Includes a priority decision flow: Safety > Concurrency > Performance > Code Quality.

##### [golang-performance](skills/golang-performance/SKILL.md)

Profile and optimize Go code using pprof, benchstat, and fgprof. Covers CPU profiling, memory profiling, GC tuning, and optimization patterns (slice preallocation, strings.Builder, struct alignment, sync.Pool, GOMEMLIMIT, HTTP transport tuning).

#### Code Quality & Review

##### [general-code-review](skills/general-code-review/SKILL.md)

Conduct a comprehensive automated code review covering logical correctness, security vulnerabilities, performance bottlenecks, and architectural alignment. Skips syntax/formatting nits and provides ranked critical issues with surgical fix snippets.

##### [golang-code-review](skills/golang-code-review/SKILL.md)

Go-specific code review skill that builds on `code-review-principles`. Focuses on safety violations (nil maps, error discarding, missing defer), error handling (unwrapped errors, log-and-return, panic misuse), concurrency (goroutine leaks, loop variable capture, context in struct), and code quality. Provides a severity decision flow and checklist for reviewing Go code.

#### Markdown & Documentation Utilities

##### [book-to-md](skills/book-to-md/SKILL.md)

Summarize long-form technical documentation, manuals, or books into a structured hierarchical Markdown format with custom markers (`▼`, `###`, `*`).

##### [enrich-md-kb](skills/enrich-md-kb/SKILL.md)

Surgically enrich existing technical summaries with deeper detail — missing commands, configuration flags, performance tuning parameters, and architectural rationale — while preserving the established structural hierarchy.

##### [fix-markdown](skills/fix-markdown/SKILL.md)

Use when the user explicitly requests to fix markdown style or formatting errors using markdownlint.

##### [kb-to-md](skills/kb-to-md/SKILL.md)

Summarize technical information or text into specific structured markdown layouts (bullets, subsection, chapter-subsection, chapter) and surgically copy/merge it into a target file.

##### [markdown-style-principles](skills/markdown-style-principles/SKILL.md)

Foundation skill providing shared formatting and layout rules for high-density Markdown documents. Not invoked directly; loaded as a prerequisite by layout-specific skills.

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

#### [writing-skills](skills/writing-skills/SKILL.md) — from [barretstorck/gemini-superpowers](https://github.com/barretstorck/gemini-superpowers)

Framework and guidelines for writing, testing, and verifying new agent skills. Follows a test-driven approach for agent documentation.

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
- `~/.gemini/skills/`

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

## Tooling & Validation

### 1. Skill Structure Linter

To validate that all skills follow the repository's structure and configuration standards, run the built-in skill linter:

```bash
python3 .agents/skills/m-skill-lint/scripts/lint
```

The linter validates:

- The presence of `SKILL.md` and its frontmatter structure.
- Correct license and attribution files for third-party skills.
- The alignment of the README.md index with actual skill folders.
- Integrity of `scripts/config/default-skills.txt`.

### 2. Markdown Style Linter

Developers can optionally validate the style and formatting of the repository's markdown files manually by running `markdownlint-cli2`:

```bash
npx markdownlint-cli2 "**/*.md" "#node_modules"
```
