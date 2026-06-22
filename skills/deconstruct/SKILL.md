---
name: deconstruct
description: Use when analyzing a codebase structure, mapping domains, tracing execution paths, or writing/updating a deconstructed codebase reference document (e.g., deconstructed.md).
---

# Codebase Deconstruction

## Overview
Deconstructing a codebase produces a comprehensive, highly traceable reference document (`deconstructed.md`) that maps high-level architecture to low-level implementation details. This serves as an onboarding guide and a high-fidelity context source for humans and agents alike.

---

## When to Use
- **New Codebase Entry**: Onboarding to a new repository to map its boundaries and entry points.
- **Architectural Reference**: Documenting how system components decoupled from each other interact.
- **Traceability Needs**: Providing exact line-level references for critical flows, constants, and data schemas.

### When NOT to Use
- Do not run deconstruction on small, trivial scripts.
- Do not run deconstruction if the goal is only single-line debugging or minor bug fixing.

---

## Discovery Workflow (Using code-review-graph MCP)
Follow this structured approach to gather insights using the `code-review-graph` MCP tools before writing or editing the document:

1. **Verify / Update Graph**: Run `build_or_update_graph_tool` to ensure the Tree-sitter knowledge graph is fresh. Check stats with `list_graph_stats_tool`.
2. **Analyze Communities**:
   - Run `get_architecture_overview_tool` (with `detail_level: "minimal"`) to discover major directory communities, cohesion metrics, and package coupling.
   - Use `list_communities_tool` and `get_community_tool` to identify files/classes in key domains.
3. **Trace Execution Flows**:
   - Run `list_flows_tool` to find major execution entry points, paths, and criticality scores.
   - Run `get_flow_tool` on high-criticality flows (e.g., `SyncCollection`, `searchCmd`) to extract precise call graphs, line ranges, and function details.
4. **Identify Hubs and Boundaries**:
   - Run `get_hub_nodes_tool` to find central classes or structures.
   - Run `get_bridge_nodes_tool` to identify interfaces, API boundary layers, and interaction specifications.

---

## Target Reference Document Structure (deconstructed.md)
The generated or updated `deconstructed.md` file must be saved under the `docs/` folder (or project root if `docs/` is absent) and adhere to the following 8-section layout:

### 1. High-Level Architecture & Domain Map
- **Diagram**: An ASCII/Unicode or Mermaid architecture block diagram showing entry points, parser engines, embedding pipelines, storage, etc.
- **Domain Directory Mapping**: A detailed bulleted list matching codebase directories to high-level functional domains, summarizing their roles.

### 2. System Parameters & Configuration Precedence
- **Core Magic Variables**: Bullet points of hardcoded parameters, constants, and thresholds with their default values and exact `file:line` locations.
- **CLI Parameter Precedence**: Documentation of configuration overrides, defaults, and priority rules (e.g., CLI flags vs. configuration file settings).

### 3. Traceability Mapping Matrix
- A markdown table mapping system commands/features to their package entry points, primary data store operations, and core target structures or files.
- Column structure: `| Feature/Command | Package Entry Point | Primary DB/IO Operations | Core Target Struct / Files |`

### 4. Data Models & Interface Specifications
- **Schema Details**: SQL schema definitions (tables, columns, foreign keys, triggers) or key data structures.
- **Interface Mappings**: Relational mappings, key structures, and FAISS/index formats with exact `file:line` locations.

### 5. Entry Points & Bootstrapping
- Detailed explanation of boot mechanisms (e.g., CLI library setup, router/server init).
- Use list items targeting specific files, functions (e.g., `Execute()`, `init()`), and line ranges.

### 6. Critical Flow Walkthroughs
- Discover as many important execution flows as possible, preferably through the use of `code-review-graph` (utilizing community structure and critical flow lists to guide selection). Include step-by-step logic walks for each identified workflow.
- **Workflow Format**: Use `->` (or `→`) combined with tab/space indentation to denote caller-callee execution workflow steps, instead of drawing complex ASCII/Unicode tree hierarchies.
  Example:
  ```markdown
  * `ingest.SyncCollection`
      -> `ingest.ingestFile`
          -> `parser.GetParser`
              -> `parser.MarkdownParser.Parse`
  ```
- **Code-Review-Graph Compatible Naming**: Use query-friendly symbol names for the `code-review-graph` MCP (e.g. `Package.Function` or `Struct.Method`). Fully qualified structures like `File::Class.Method` are not necessary in most cases and should only be used when ambiguity exists (e.g., when two files contain the same class and method names). Avoid vague descriptions.

### 7. Package Boundaries & Interaction Interfaces
- **Boundary Map**: A Mermaid diagram or ASCII flow depicting how packages (e.g., `parser`, `ingest`, `embed`, `project`) communicate and decouple.
- **Decoupling Details**: Explanations of boundary patterns (e.g., parser registries, relational data mapping, status update channels).
- **Clickable Symbol Links**: Every key boundary symbol or interface should be linked using relative links (e.g., `[Parser](../internal/parser/parser.go#L25)`).

### 8. Development, Testing, and Performance Profiling
- **Build Instructions**: Commands to compile the project (including build tags like `-tags fts5`).
- **Running Tests**: Command instructions to run specific test suites and tag combinations.
- **Benchmarking & Profiling**: Specific commands to execute benchmark functions and use Go's `pprof` system to analyze CPU/Memory profiles (`go tool pprof`).

---

## Operational Guidelines
- **Prefer Relative Links**: All references to local files, interfaces, classes, structs, and methods in the generated markdown reference file should use relative paths (e.g., `[fts.go](../internal/project/fts.go#L84)`) to ensure portability across different host environments.
- **External Resources**: Use standard HTTPS URLs for documentation or assets hosted on the internet.
- **Be Concise and Structured**: Avoid hand-wavy explanations or block text. Prefer tables, lists, code fences, and clear block diagrams.
- **Incremental Updates**: If `deconstructed.md` already exists, merge new discoveries cleanly into the existing sections without breaking the 8-section layout.
