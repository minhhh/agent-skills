---
name: deconstruct
description: Incrementally deconstruct the codebase into a detailed reference document
---

Objective:
* Detailed breakdown of the current codebase across level from architectural to specific code flow
* Useful reference to humans and llms alike


Discovery Phase (run these in parallel):
* Entry Point Mapping: Identify main entry points and trace their initial setup/bootstrap logic.
* Core Symbol Discovery: Locate and summarize primary classes, interfaces, and data structures.
* Control Flow Analysis: Trace the "happy path" for the system's 2-3 most critical operations.
* Feature Grouping: Map directories and files to high-level functional domains.
* Side Effect Audit: Identify where the system interacts with external boundaries (DB, API, Disk).
* Testing Pattern Audit: Identify how core logic is verified (unit vs. integration vs. e2e).

Use the following format to document flow
* leaderboard.py
    -> save_scores(data: LeaderboardData) -> None
        -> populate_database(data=data, top=100, filter=None)
            -> DatabaseManager.update(data=data, top=top) -> ServerRequest[Data]
                -> Call server endpoint `v2/db/update`


Documentation Phase:
* Report: If deconstructed.md does not exist, write summary to it with these sections:
  - Architecture Overview: High-level design and layering.
  - Entry Points: How the system starts and handles initial requests.
  - Important Classes/Symbols: The "nouns" of the system and their responsibilities.
  - Important Flows: The "verbs" - how data moves through the core logic.
  - Specific Features: A mapping of business features to code locations.
  - Configurations: How the system is tuned and where secrets/env vars live.
  - Testing: The verification strategy and key test suites.
* Incremental Report: If deconstructed.md exists, integrate new findings into the existing sections. Use "New Discovery" tags if the information contradicts or significantly expands previous entries.


Operational Guidelines:
* Traceability: Always link code flows to specific file paths and line ranges.
* LLM Context: Keep descriptions concise and structured so they can be easily parsed for future prompts.
* Dependency Context: Note when a flow is heavily dependent on a specific third-party library.
