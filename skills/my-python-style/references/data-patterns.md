# Data Pipeline Patterns

## 4. Lossless Result Containers
Prefer returning structured **Result Objects** (using `dataclasses`) over simple filtered lists. Bins every input item into explicit categories (e.g., `valid`, `ignored`, `failed`).

**Benefits:**
- **Auditability**: Prevents silent data loss; the caller knows the fate of every input item.
- **Encapsulated Summary**: UI/CLI logic can trivially generate reports from the binned attributes.

## 5. Priority-Ordered Hybrid Matching
In search or categorization tasks, treat deterministic sources (like SQL FTS5) as the "Source of Truth" and only fall back to stochastic sources (like Vector embeddings) when exact matches are not found.