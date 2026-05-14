---
name: extend-markdown-summary
description: Extend a technical summary with deep technical details. Use when asked to "add more depth", "fill gaps", or "extend" an existing summary.
---

# Extend Markdown Summary

## Overview
This skill focuses on the **surgical enrichment** of technical summaries. It prioritizes technical depth, actionable commands, and architectural rationale while strictly adhering to the established structural hierarchy.

## Workflow

### 1. Gap Analysis
- Compare the existing summary against the source material (TOC or full text).
- **Identify Technical Voids**: Look for missing specific commands, configuration flags, performance tuning parameters, or "under the hood" explanations.
- **Pinpoint Omissions**: Find sections where the "why" behind a best practice is missing.

### 2. Targeted Context Retrieval
- **File Access**: Use `read_file` or `grep_search` if the reference is a local file.
- **External Specs**: Use `google_web_search` or `web_fetch` ONLY if the primary source is insufficient for the requested depth.

### 3. Surgical Integration
- **Follow the Hierarchy**:
    - Chapter: `###`
    - Sub-section: `▼ **Bold Title**`
    - Key Point: `* **Bold Term**: Detailed explanation.`
    - Nuance/Command: `    *` (Nested list)
- **Deep Insertion**: Add new bullet points or nested nuances directly into existing sub-sections to maintain logical continuity.
- **Maintain Tone**: Use direct, technical language. **Avoid third-person phrasing** (e.g., use "Metric X indicates..." instead of "The author explains that metric X indicates...").

## Guidelines
- **Actionable Specificity**: Every extension must add value. Prefer `sar -n DEV 1` over "Check network statistics."
- **Visual Consistency**: Always use the `▼` marker for sub-sections.
- **Code Blocks**: Wrap all commands and config snippets in language-specific blocks (e.g., ```bash).
- **No Redundancy**: Do not repeat information already present in the summary.

## Example: Enrichment of an Existing Section

**Original:**
▼ **Disk I/O Latency**
* Use `iostat` to check for disk bottlenecks.

**Extended:**
▼ **Disk I/O Latency**
* **Interpreting `%util`**: High utilization in `iostat -x` suggests the disk is saturated.
    * **The `await` Metric**: If `await` is significantly higher than `svctm`, the application is blocked waiting for I/O.
    * `iostat -xz 1`: Use this to filter out idle disks and focus on active hotspots.
* **Kernel I/O Schedulers**: Check `/sys/block/[dev]/queue/scheduler` to ensure the correct policy (e.g., `none` for NVMe, `mq-deadline` for SSDs) is applied.
