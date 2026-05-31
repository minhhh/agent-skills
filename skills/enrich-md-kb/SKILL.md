---
name: enrich-md-kb
description: Use only when explicitly requested by the user to enrich, extend, or deepen specific sections of a technical summary.
---

# Enrich MD KB

## Overview

This skill focuses on the **surgical enrichment** of technical summaries. It prioritizes technical depth, actionable commands, and architectural rationale while strictly adhering to the established structural hierarchy.

## When to Use

Use this skill **only** when the user explicitly requests to enrich, extend, deepen, or add technical details to a specific section of a technical summary.

* **Triggering Conditions**:
  - The user explicitly instructs you to "enrich", "extend", "deepen", "add depth", or "fill gaps" in a technical summary.
  - The user specifies a particular target they want expanded with technical commands or architectural details, which could be a chapter (`###` heading), a sub-section (`▼` heading), or a specific root bullet point.
* **When NOT to Use**:
  - **Do NOT** proactively enrich summaries, sections, or commands unless explicitly asked, even if you want to show engineering value or believe a senior user persona expects deep details.
  - **Do NOT** trigger this skill automatically for general summarization requests (e.g., "Summarize this log" or "Write a summary of DB performance").
  - **Do NOT** assume that just because the skill is loaded, you should use it; skills are tools, and only use them when their specific triggering conditions are met.

## The Two-Phase Workflow

The enrichment process must **always** be executed in two distinct phases. Do NOT apply changes to any file until Phase 1 is complete and approved by the user.

### Phase 1: Suggesting the Changes
1. **Retrieve and Analyze**: Identify gaps or fetch URL content according to the requested [Enrichment Styles](#enrichment-styles).
2. **Draft the Proposal**: Draft the proposed additions following the [Formatting and Layout Rules](#formatting-and-layout-rules).
3. **Present Draft**: Present the proposed changes to the user as a draft (e.g., in a markdown diff block) and wait for explicit approval.

### Phase 2: Applying the Changes
1. **Surgical Integration**: Once the user approves the draft, integrate the changes into the target file.
2. **Verify Layout**: Verify the final document structure complies with all layout rules.

## Enrichment Styles

Enrichment follows one of two styles based on the request and existing content:

### Style 1: Standard Gap Filling
Use this when asked to add technical depth, missing commands, configuration flags, or architectural explanation to an existing section.
* Compare the summary against reference documentation or logs to identify technical voids.
* Propose deep, specific details and actionable commands to fill those specific gaps.

### Style 2: URL Reference Merging
Use this when a root bullet point or section in the summary contains a URL reference.
* **Read the Source**: Fetch and read the content of the referenced URL (using `read_url_content` or `read_browser_page`).
* **Summarize and Merge**: Document the summary of the URL. Copy and merge the relevant structured points from that URL directly into the existing bullet points under that section.

## Formatting and Layout Rules

All drafts (Phase 1) and final integrations (Phase 2) must adhere to these rules:

### 1. Hierarchy Rules
* **Chapter**: `### Chapter Name`
* **Sub-section**: `▼ **Bold Title**`
* **Key Point**: `* **Bold Term**: Detailed explanation.`
* **Nuance/Command**: `*` (Nested list)

### 2. Constraints
* Bullet points must **always** use asterisks (`*`), never hyphens (`-`).
* Sibling root bullet points must **always** be separated by exactly one blank line.
* Sub-points must **always** be indented with exactly 4 spaces relative to their parent (four spaces followed by an asterisk).
* **Deep Insertion**: Insert new bullets directly into existing sub-sections to maintain logical continuity.
* **Tone**: Use direct, technical language. **Avoid third-person phrasing** (e.g., use "Metric X indicates..." instead of "The author explains that metric X indicates...").
* **Code Blocks**: Wrap all commands and config snippets in language-specific blocks.

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
