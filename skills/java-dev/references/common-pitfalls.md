# Common Pitfalls — These Thoughts Mean STOP

If you catch yourself thinking any of these, **STOP** and apply the correct approach:

| Rationalization | Problem | Impact | Fix |
|-----------------|---------|--------|-----|
| "Resource will close automatically" | Missing try-with-resources | FD exhaustion after 20hrs | Wrap in try-with-resources |
| "This is single-threaded, no sync needed" | Undocumented thread model | Future bugs when threading added | Add `// NOT thread-safe` comment |
| "I'll add the test after I finish this" | No test coverage | Gaps never get filled | Add integration test now |
| "This is performance-critical, streams are too slow" | Premature optimization | Bugs from complex code | Measure first with profiler |
| "Just this once I'll catch and ignore the exception" | Swallowed exception | Silent failures, lost data | Log exception or rethrow |
| "ThreadLocal cleanup isn't critical here" | Classloader leak | OOM after 10 deployments | Remove in finally block |
| "The lock order doesn't matter for this simple case" | Undocumented lock order | Deadlock when code grows | Document ordering now |
| "This allocation is trivial" | Boxing in hot loop | GC pressure, latency spikes | Use primitive types |
| "I'll use HashMap, order doesn't matter" | Non-deterministic ordering | Build flakiness | Use LinkedHashMap/TreeMap |
| "Mockito is faster than a real test database" | Mocked database | Mock/prod drift, broken prod | Use real DB (Testcontainers/QuarkusTest) |
| "Let me refactor this code I haven't read yet" | Refactoring unknown code | Breaking working functionality | Read and understand first |
| "I'll just use the class/package name as a String" | String class/package reference | Silently breaks on rename/move | Use `.class` reference |
| "I know this blocks, but it's quick" | Blocking event-loop (Quarkus/Vert.x) | Cascading 503 errors | Use `@Blocking` annotation |
| "I'll use field injection, it's less code" | Field injection prevents testing | Hard to mock, hidden deps | Use constructor injection |
