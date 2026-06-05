---
name: enrich-md-kb
description: Use when explicitly requested by the user to apply the skill to extend (add more relevant points), deepen (add more technical details to existing points), reformat (do not change the content, just fix alignment and formatting issues), or a combination of those 3. Builds on markdown-style-principles.
---

# Enrich MD KB

## Overview

This skill focuses on the **surgical extension, deepening, or reformatting** of technical summaries. It prioritizes technical depth, actionable commands, and architectural rationale while strictly adhering to the established structural hierarchy.

## Prerequisites

**This skill builds on [markdown-style-principles](file:///Users/minh/gitrepos/local_tools/agent-skills/skills/markdown-style-principles/SKILL.md)**.

## When to Use

Use this skill **only** when the user explicitly requests to apply the skill to extend (add more relevant points), deepen (add more technical details to existing points), reformat (do not change the content, just fix alignment and formatting issues), or a combination of those 3.

* **Triggering Conditions**:
  - **Extend**: The user explicitly requests to add more relevant points to the technical summary or specific sections.
  - **Deepen**: The user explicitly requests to add more technical details (such as commands, configurations, or architecture explanations) to existing points.
  - **Reformat**: The user explicitly requests to fix alignment, structure, and formatting issues (without modifying or changing the actual content).
  - Any combination of the above three actions.
* **When NOT to Use**:
  - **Do NOT** proactively extend, deepen, or reformat summaries, sections, or commands unless explicitly asked, even if you want to show engineering value or believe a senior user persona expects deep details.
  - **Do NOT** trigger this skill automatically for general summarization requests (e.g., "Summarize this log" or "Write a summary of DB performance").
  - **Do NOT** assume that just because the skill is loaded, you should use it; skills are tools, and only use them when their specific triggering conditions are met.

## The Two-Phase Workflow

The process must **always** be executed in two distinct phases. Do NOT apply changes to any file until Phase 1 is complete and approved by the user.

### Phase 1: Suggesting the Changes
1. **Retrieve and Analyze**: Identify gaps, fetch URL content, or identify formatting/alignment violations according to the requested [Application Styles](#application-styles).
2. **Draft the Proposal**: Draft the proposed additions or corrections following the [Formatting and Layout Rules](#formatting-and-layout-rules).
3. **Present Draft**: Present the proposed changes to the user as a draft (e.g., in a markdown diff block) and wait for explicit approval.

### Phase 2: Applying the Changes
1. **Surgical Integration**: Once the user approves the draft, integrate the changes into the target file.
2. **Verify Layout**: Verify the final document structure complies with all layout rules.

## Application Styles

Applying the skill follows one of three styles based on the request and existing content:

### Style 1: Standard Gap Filling (Extend/Deepen)
Use this when asked to add technical depth, missing commands, configuration flags, or architectural explanation to an existing section.
* Compare the summary against reference documentation or logs to identify technical voids.
* Propose deep, specific details and actionable commands to fill those specific gaps.

### Style 2: URL Reference Merging (Extend/Deepen)
Use this when a root bullet point or section in the summary contains a URL reference.
* **Read the Source**: Fetch and read the content of the referenced URL (using `read_url_content` or `read_browser_page`).
* **Summarize and Merge**: Document the summary of the URL. Copy and merge the relevant structured points from that URL directly into the existing bullet points under that section.

### Style 3: Reformatting (Reformat)
Use this when asked to fix alignment and formatting issues without changing the content.
* **Preserve Content**: Do NOT add, remove, or modify any existing text content or semantic meaning.
* **Fix Formatting**: Correct bullet markers, nested indentation (ensure exactly 4 spaces relative to the parent bullet), and empty line spacing between root bullet points starting with `*` to match the structural rules. Convert any plain paragraphs into root or nested bullet points starting with `*`. Existing numbered bullet points (e.g., `1. **Heading**`) and their spacing should be preserved and not converted or modified.

## Formatting and Layout Rules

All drafts (Phase 1) and final integrations (Phase 2) must adhere to the style rules defined in [markdown-style-principles](file:///Users/minh/gitrepos/local_tools/agent-skills/skills/markdown-style-principles/SKILL.md).

- **Supported Layouts**: This skill preserves and extends whichever layout style (`bullets`, `subsection`, `chapter-subsection`, or `flat-chapter`) is already established in the target document.
- **Reformatting Rule (Style 3)**: When fixing alignment or formatting issues, convert existing malformed structures to strictly comply with the indentation, spacing, and hierarchy rules defined in [markdown-style-principles](file:///Users/minh/gitrepos/local_tools/agent-skills/skills/markdown-style-principles/SKILL.md).

## Guidelines

- **Actionable Specificity**: Every extension must add value. Prefer `sar -n DEV 1` over "Check network statistics."
- **Visual Consistency**: Always use the `▼` marker for sub-sections.
- **No Redundancy**: Do not repeat information already present in the summary.


## Examples

### Example 1: Standard Gap Filling (Style 1)

**Original:**
```markdown
▼ **Disk I/O Latency**

* **iostat**: Use this command to check for disk bottlenecks.
```

**Proposed Draft (Approved & Applied):**
```markdown
▼ **Disk I/O Latency**

* **iostat**: Use this command to check for disk bottlenecks.
    * **Interpreting `%util`**: High utilization in `iostat -x` suggests the
      disk is saturated.
    * **The `await` Metric**: If `await` is significantly higher than `svctm`,
      the application is blocked waiting for I/O.
    * `iostat -xz 1`: Use this to filter out idle disks and focus on active
      hotspots.

* **Kernel I/O Schedulers**: Check `/sys/block/[dev]/queue/scheduler` to ensure
  the correct policy (e.g., `none` for NVMe, `mq-deadline` for SSDs) is
  applied.
```

### Example 2: URL Reference Merging (Style 2)

**Original:**
```markdown
▼ **Garbage Collection Optimization**

* **GC Tuning**: [Go GC Guide](https://go.dev/doc/gc-guide)
* Go GC Guide 2: https://go.dev/doc/gc-guide-2
```

**Proposed Draft (Approved & Applied):**
```markdown
▼ **Garbage Collection Optimization**

* **GC Tuning**: [Go GC Guide](https://go.dev/doc/gc-guide)
    * **Memory Limit**: Set `GOMEMLIMIT` to configure a soft memory limit, reducing GC frequency near limits without risking Out-Of-Memory panic.
    * **GC Percent**: Adjust `GOGC` to tune the trade-off between CPU overhead and memory footprint.
    * **Latency Goals**: Setting `GOMEMLIMIT` helps the runtime meet latency goals under high memory utilization.

* Go GC Guide 2: https://go.dev/doc/gc-guide-2
    * **Main point**: The content
```
