---
name: golang-performance
description: Use when debugging slow Go code, profiling performance, optimizing bottlenecks, or improving application memory and CPU usage.
---

# Go Performance Optimization

Comprehensive guide to profiling, analyzing, and optimizing Go code for better performance, including CPU profiling, memory optimization, and implementation best practices.

## When to Use This Skill

- Identifying performance bottlenecks in Go applications
- Reducing application latency and response times
- Optimizing CPU-intensive operations
- Reducing memory consumption and GC pressure
- Tuning HTTP clients and connection pools
- Improving data processing pipeline throughput
- Container memory tuning (GOMEMLIMIT, GOGC)

## Decision Tree: Where Is Time Spent?

Before optimizing, identify the bottleneck type to pick the right pattern:

| Bottleneck | Signal (from pprof) | Action |
| --- | --- | --- |
| Too many allocations | `alloc_objects` high in heap profile | [Pattern 1: Slice Preallocation](#pattern-1-slice-preallocation), [Pattern 4: sync.Pool](#pattern-4-syncpool-for-hot-allocations), [Pattern 10: Map Size Hints, Slice Reuse, and Direct Indexing](#pattern-10-map-size-hints-slice-reuse-and-direct-indexing) |
| CPU-bound hot loop | function dominates CPU profile | Check algorithmic complexity, [Pattern 11: CPU Cache Locality & Layout](#pattern-11-cpu-cache-locality--layout-row-major-vs-column-major-soa-vs-aos), [Pattern 12: Function Inlining](#pattern-12-function-inlining--scheduler-preemption) |
| GC pauses / OOM | high GC%, container limits | [Pattern 5: GOMEMLIMIT](#pattern-5-gomemlimit-for-containerized-environments), [Pattern 9: Backing Array Memory Leaks](#pattern-9-backing-array-memory-leaks-slices-substrings-and-maps) |
| Network / I/O latency | goroutines blocked on I/O | [Pattern 6: HTTP Transport Tuning](#pattern-6-http-transport-tuning), [Pattern 13: High-Concurrency HTTP & I/O Tuning](#pattern-13-high-concurrency-http--io-tuning-draining-bodies-streaming-vs-buffering) |
| Repeated expensive work | same computation/fetch multiple times | [Pattern 14: Compiled Patterns, Precomputed Tables, & Singleflight](#pattern-14-compiled-patterns-precomputed-tables--singleflight) |
| Wrong algorithm | O(n²) where O(n) exists | Use slices/maps instead of linear search |
| Lock contention | mutex/block profile hot | Reduce lock scope, use atomic or channel |
| Slow external calls / scheduler starvation | DB/network time dominates or runnable goroutines | Tune query, add connection pool, [Pattern 12: Scheduler Preemption](#pattern-12-function-inlining--scheduler-preemption), [Pattern 15: PGO & GOMAXPROCS](#pattern-15-profile-guided-optimization-pgo--gomaxprocs-in-containers) |

**Rule out external bottlenecks first**: If 90% of latency is a slow DB query, reducing allocations won't help. Use fgprof to check off-CPU time.

## Core Concepts

### 1. Profiling Types

- **CPU Profiling**: Identify time-consuming functions with pprof
- **Heap/Memory Profiling**: Track allocation hot spots and GC pressure
- **Goroutine Profiling**: Detect leaked or excessive goroutines
- **Block/Mutex Profiling**: Find lock contention and synchronization bottlenecks
- **Tracing**: Capture the full timeline of events with `go tool trace`

### 2. Performance Metrics

- **Execution Time**: How long operations take (ns/op)
- **Allocations**: Number and size of allocations per operation (B/op, allocs/op)
- **GC Pause**: Time spent in garbage collection
- **Goroutine Count**: Indicator of leaks or excessive concurrency
- **Lock Contention**: Time spent waiting on mutexes

### 3. Optimization Strategies

- **Allocation reduction**: Avoid heap allocations in hot paths
- **Algorithmic**: Better data structures and algorithms
- **Concurrency**: Fan-out, worker pools, connection pooling
- **Caching**: Avoid redundant computation or I/O
- **Runtime tuning**: GOMEMLIMIT, GOGC, GOMAXPROCS

## Quick Start

### Basic Timing

```go
func measureTime() {
    start := time.Now()

    result := 0
    for i := 0; i < 1_000_000; i++ {
        result += i
    }

    elapsed := time.Since(start)
    fmt.Printf("Execution time: %v\n", elapsed)
    _ = result
}
```

### Benchmark (preferred — statistically rigorous)

```go
func BenchmarkSum(b *testing.B) {
    for i := 0; i < b.N; i++ {
        result := 0
        for j := 0; j < 1_000_000; j++ {
            result += j
        }
        _ = result
    }
}
```

```bash
go test -bench=BenchmarkSum -benchmem -count=6 . | tee report-before.txt
```

## Profiling Tools

### Tool 1: pprof — CPU and Memory Profiling

**CPU profile — find hot functions:**

```go
import "runtime/pprof"

func main() {
    f, _ := os.Create("cpu.pprof")
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()

    // Your code here
    slowOperation()
}
```

```bash
# Or use the testing framework
go test -cpuprofile=cpu.pprof -memprofile=mem.pprof -bench=. .

# Analyse
go tool pprof -http=:8080 cpu.pprof
go tool pprof -http=:8080 mem.pprof
```

**Heap profile — find allocation hot spots:**

```bash
# Live heap profile from a running server
curl http://localhost:8080/debug/pprof/heap -o heap.pprof
go tool pprof -alloc_objects heap.pprof
go tool pprof -alloc_space heap.pprof
```

Every HTTP server should register `net/http/pprof`:

```go
import _ "net/http/pprof"
```

### Tool 2: benchstat — Compare Benchmarks Statistically

```bash
# Run baseline and candidate
go test -bench=BenchmarkMyFunc -benchmem -count=6 ./pkg > before.txt
# Make changes...
go test -bench=BenchmarkMyFunc -benchmem -count=6 ./pkg > after.txt

# Compare with statistical significance
benchstat before.txt after.txt
```

`benchstat` accounts for noise and tells you if the change is statistically significant. Never compare single benchmark runs.

### Tool 3: fgprof — On-CPU + Off-CPU Profiling

`fgprof` captures time spent both on CPU and waiting (I/O, channels, locks). Use when standard pprof CPU profile suggests the CPU is idle but latency is high.

```go
import "github.com/felixge/fgprof"

http.Handle("/debug/fgprof", fgprof.Handler())
```

```bash
curl http://localhost:8080/debug/fgprof -o fgprof.pprof
go tool pprof -http=:8080 fgprof.pprof
```

### Tool 4: go tool trace — Full Execution Timeline

```go
import "runtime/trace"

func main() {
    f, _ := os.Create("trace.out")
    trace.Start(f)
    defer trace.Stop()

    yourCode()
}
```

```bash
go tool trace trace.out
```

The trace viewer shows goroutine lifecycles, network blocking, GC events, and syscall activity in a timeline. Use when pprof shows a function is slow but not *why* (blocked on I/O? GC? channel?).

## Optimization Patterns

### Pattern 1: Slice Preallocation

```go
// Bad: Grows backing array multiple times
func collect(users []User) []string {
    var ids []string
    for _, u := range users {
        ids = append(ids, u.ID)
    }
    return ids
}

// Good: Single allocation
func collect(users []User) []string {
    ids := make([]string, 0, len(users))
    for _, u := range users {
        ids = append(ids, u.ID)
    }
    return ids
}
```

**Impact**: Avoids repeated slice growth (allocate → copy → GC old backing array).

### Pattern 2: strings.Builder Instead of `+` in Loops

```go
// Bad: O(n²) — each += allocates a new string
func join(parts []string) string {
    var result string
    for _, p := range parts {
        result += p
    }
    return result
}

// Good: Single growing buffer
func join(parts []string) string {
    var sb strings.Builder
    for _, p := range parts {
        sb.WriteString(p)
    }
    return sb.String()
}

// Best: Use standard library
func join(parts []string) string {
    return strings.Join(parts, "")
}
```

**Impact**: `strings.Builder` avoids O(n²) string copying. Always prefer `strings.Join` when a delimiter is used.

### Pattern 3: Struct Field Alignment

Field order affects struct size due to padding. Arrange fields from largest to smallest alignment requirement.

```go
// Bad: 40 bytes — padding between bool (1) and string (16)
type Bad struct {
    active bool    // 1 byte + 7 padding
    name    string // 16 bytes
    count   int64  // 8 bytes
    debug   bool   // 1 byte + 7 padding
}

// Good: 32 bytes — no wasted padding
type Good struct {
    name   string // 16 bytes
    count  int64  // 8 bytes
    active bool   // 1 byte
    debug  bool   // 1 byte
    // 6 bytes padding (end of struct)
}
```

**Diagnose**: `fieldalignment` (from `golang.org/x/tools`):

```bash
go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
fieldalignment -fix ./pkg/...
```

### Pattern 4: sync.Pool for Hot Allocations

```go
var bufPool = sync.Pool{
    New: func() any { return new(bytes.Buffer) },
}

func handle(data []byte) []byte {
    buf := bufPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufPool.Put(buf)
    }()
    buf.Write(data)
    return buf.Bytes()
}
```

**Impact**: Reduces GC pressure on frequently allocated objects. **Only use when profiling shows allocation is a bottleneck** — adds complexity.

### Pattern 5: GOMEMLIMIT for Containerized Environments

```bash
# Set to 80-90% of container memory limit
export GOMEMLIMIT=$((80 * $(cat /sys/fs/cgroup/memory.max) / 100))
```

**Without GOMEMLIMIT**: Go's GC triggers at 2× live heap (GOGC=100). In memory-constrained containers, the heap can grow past the container limit before GC catches up → OOM kill.

**Impact**: Prevents OOM kills in containerized Go services. The soft memory limit gives the GC a target to stay below, with a 50% margin before hard limit.

### Pattern 6: HTTP Transport Tuning

The default `http.Client.Transport` caps `MaxIdleConnsPerHost` at 2 — a common production bottleneck.

```go
// Bad: Default transport — 2 connections per host
client := http.DefaultClient

// Good: Tuned for concurrency
transport := &http.Transport{
    MaxIdleConns:        100,
    MaxIdleConnsPerHost: 10,
    IdleConnTimeout:     90 * time.Second,
    DialContext: (&net.Dialer{
        Timeout:   5 * time.Second,
        KeepAlive: 30 * time.Second,
    }).DialContext,
}
client := &http.Client{Transport: transport}
```

**Impact**: Eliminates connection-establishment latency under concurrency.

### Pattern 7: Avoid `reflect.DeepEqual` in Hot Paths

```go
// Bad: 50-200x slower than typed comparison
if reflect.DeepEqual(a, b) {
    return true
}

// Good: Standard library typed helpers (Go 1.21+)
if slices.Equal(a, b) { ... }
if maps.Equal(a, b) { ... }
if bytes.Equal(a, b) { ... }
```

**Impact**: `reflect.DeepEqual` uses reflection, allocates, and cannot be inlined. The standard library helpers are generics-based, type-checked at compile time, and allocation-free.

### Pattern 8: Logging in Hot Loops

```go
// Bad: fmt.Print* allocates; log.Print* allocates and serializes
for _, item := range items {
    log.Printf("processing item %s: %d", item.ID, item.Value)
}

// Good: Sum outside loop, log once
count := 0
for _, item := range items {
    count++
}
log.Printf("processed %d items", count)

// For debug logging that must stay: use slog.LogAttrs (no allocation when disabled)
if logger.Enabled(ctx, slog.LevelDebug) {
    logger.LogAttrs(ctx, slog.LevelDebug, "processing item",
        slog.String("id", item.ID),
        slog.Int("value", item.Value),
    )
}
```

**Impact**: Logging inside loops prevents inlining, allocates even at disabled log levels (with older logging libraries), and distorts benchmark results.

### Pattern 9: Backing Array Memory Leaks (Slices, Substrings, and Maps)

#### Slice Reslicing

A small reslice of a large slice keeps the entire original backing array in memory, preventing garbage collection.

```go
// Bad: Retains entire megabyte-sized backing array
func getHeader(data []byte) []byte {
    return data[:16]
}

// Good: Independent copy, original backing array can be GC'd
func getHeader(data []byte) []byte {
    header := make([]byte, 16)
    copy(header, data[:16])
    return header
}
```

#### Substring Leaks

Substrings share the backing array of the original string.

```go
// Bad: Keeps the entire parent string in memory
func extractID(msg string) string {
    return msg[:8]
}

// Good: Independent copy (Go 1.20+)
func extractID(msg string) string {
    return strings.Clone(msg[:8])
}
```

#### Map Compaction

Go maps grow but never release bucket memory when entries are deleted. A map that once held millions of entries retains its memory footprint forever.

```go
// Good: Recreate map periodically to reclaim bucket memory
func compact(old map[string]Data) map[string]Data {
    m := make(map[string]Data, len(old))
    for k, v := range old {
        m[k] = v
    }
    return m // old map becomes eligible for GC
}
```

**Impact**: Prevents holding onto large, unused backing allocations, reducing heap residency and memory footprint.

### Pattern 10: Map Size Hints, Slice Reuse, and Direct Indexing

#### Map Size Hints

Initializing a map with `make(map[K]V)` starts with a default small bucket count. As the map grows, the runtime allocates new bucket spaces and rehashes all keys, causing CPU and memory overhead.

```go
// Bad: Repeated rehashing and bucket allocations
m := make(map[string]int)
for _, item := range items {
    m[item.Key] = item.Val
}

// Good: Single allocation, no rehashing
m := make(map[string]int, len(items))
for _, item := range items {
    m[item.Key] = item.Val
}
```

#### Reuse Slices via `append(s[:0], ...)`

Reslicing to zero length (`s[:0]`) keeps the backing array but resets the length to 0. This allows reusing the same allocation across operations without garbage collection overhead.

```go
// Bad: Allocates a new slice on every iteration
for _, item := range batch {
    results := []T{item} // allocates every time
    process(results)
}

// Good: Reuses existing backing array capacity (0 allocations after the first)
var buf []T
for _, item := range batch {
    buf = append(buf[:0], item)
    process(buf)
}
```

#### Direct Indexing vs Append

When the output slice size is exactly known and matches the input size, allocate the slice with the final size and assign directly via index. This avoids the length increment overhead and per-element bounds-checking optimization blocks of `append`.

```go
// Slower: append overhead per element
result := make([]T, 0, len(input))
for i := range input {
    result = append(result, transform(input[i]))
}

// Faster: direct assignment
result := make([]T, len(input))
for i := range input {
    result[i] = transform(input[i])
}
```

**Impact**: Eliminates map rehashing overhead and micro-allocations in slice generation loops.

### Pattern 11: CPU Cache Locality & Layout (Row-Major vs Column-Major, SoA vs AoS)

#### Row-Major Traversal

Go stores 2D arrays/slices in row-major order. Iterating column-first jumps across memory boundaries, causing severe CPU cache misses.

```go
// Bad: Column-first traversal (frequent cache misses)
for col := 0; col < 1024; col++ {
    for row := 0; row < 1024; row++ {
        sum += matrix[row][col]
    }
}

// Good: Row-first traversal (sequential memory access)
for row := 0; row < 1024; row++ {
    for col := 0; col < 1024; col++ {
        sum += matrix[row][col]
    }
}
```

#### Contiguous 2D Slice Allocation

Allocating each row of a 2D matrix individually scatters them across the heap, degrading cache performance.

```go
// Bad: N separate allocations, poor cache locality
matrix := make([][]float64, rows)
for i := range matrix {
    matrix[i] = make([]float64, cols)
}

// Good: Single contiguous allocation, cache-friendly
data := make([]float64, rows*cols)
matrix := make([][]float64, rows)
for i := range matrix {
    matrix[i] = data[i*cols : (i+1)*cols]
}
```

#### Struct of Arrays (SoA) vs Array of Structs (AoS)

When iterating over a single field of a large struct inside a slice, loading the entire Array of Structs (AoS) into cache lines wastes cache capacity. Struct of Arrays (SoA) keeps the specific field values contiguous in memory.

```go
// AoS: Loading Point (24 bytes) to read only X (8 bytes) = 66% cache line waste
type Point struct { x, y, z float64 }
points := make([]Point, n)
for i := range points {
    sum += points[i].x
}

// SoA: Contiguous X values, 100% cache line utilization
type Points struct {
    xs []float64
    ys []float64
    zs []float64
}
for i := range ps.xs {
    sum += ps.xs[i]
}
```

**Impact**: Maximizes CPU L1/L2 cache line utilization, yielding up to 10-50x speedups for memory-bound iterations.

### Pattern 12: Function Inlining & Scheduler Preemption

#### Inline-Friendly Coding

The Go compiler automatically inlines small, simple functions to eliminate function call overhead. Complexity like loops, long switch statements, or logging/side-effects blocks inlining.

```go
// Bad: Logging call blocks compiler inlining
func abs(x int) int {
    if x < 0 {
        log.Printf("negative value: %d", x)
        return -x
    }
    return x
}

// Good: Tiny, inlineable helper. Keep logging/metrics out of hot path helpers
func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}
```

#### Value Receivers for Method Chaining

Fluent builder patterns or method chains using pointer receivers block the compiler from fully inlining the calls due to pointer indirection. Using value receivers allows complete inlining.

```go
// Bad: Pointer receiver blocks inlining of fluent chain calls
func (c *config) WithTimeout(d time.Duration) *config {
    c.timeout = d
    return c
}

// Good: Value receiver allows full inlining and up to 80% speedup in chains
func (c config) WithTimeout(d time.Duration) config {
    c.timeout = d
    return c
}
```

#### Scheduler Preemption in Tight Loops

Goroutines running long, CPU-intensive tight loops with fully inlined operations and no function calls can monopolize the thread, starving other goroutines (starving scheduler preemption). You can force a preemption yield point using `//go:noinline` on a helper function called inside the loop.

```go
// Bad: Can starve scheduler if execution runs for hundreds of milliseconds
for {
    x = x*a + b // fully inlined, pure arithmetic, no preemption check
}

// Good: Forces function call boundary, which acts as a scheduler preemption point
//go:noinline
func processBatch(item WorkItem) {
    // CPU-intensive work here
}

for item := range work {
    processBatch(item) // guaranteed preemption point
}
```

**Impact**: Eliminates call overhead for tiny helpers via inlining, while preventing tight, fully inlined loops from causing goroutine starvation.

### Pattern 13: High-Concurrency HTTP & I/O Tuning (Draining Bodies, Streaming vs Buffering)

#### Draining Response Bodies for Connection Reuse

The HTTP connection pool only reuses connections that have been fully read to EOF and closed. If you discard the response body without draining it, the connection is closed rather than returned to the pool.

```go
// Bad: Connection is closed and cannot be reused
resp, err := client.Get(url)
if err != nil {
    return err
}
defer resp.Body.Close()

// Good: Fully read and discard to reuse the connection
resp, err := client.Get(url)
if err != nil {
    return err
}
defer resp.Body.Close()
_, _ = io.Copy(io.Discard, resp.Body) // drain body
```

#### Stream Processing vs Buffering

Avoid buffering entire files or response streams in memory with `io.ReadAll`. Instead, use `bufio.Scanner` or stream contents directly to process data with O(1) memory footprint.

```go
// Bad: Buffers entire file in memory
data, _ := io.ReadAll(file)
for _, line := range parseLines(data) {
    process(line)
}

// Good: Process line-by-line using streaming buffers
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    process(scanner.Bytes())
}
```

#### Streaming JSON Decoders

For large JSON streams or array payloads, use `json.NewDecoder` to parse objects one-by-one. This avoids allocating a huge contiguous block of memory to fit the entire JSON payload in `json.Unmarshal`.

```go
// Bad: Buffers full JSON body in memory
var items []Item
if err := json.Unmarshal(bodyBytes, &items); err != nil {
    return err
}

// Good: Stream decoder parses objects sequentially
dec := json.NewDecoder(resp.Body)
for dec.More() {
    var item Item
    if err := dec.Decode(&item); err != nil {
        return err
    }
    process(item)
}
```

**Impact**: Eliminates out-of-memory spikes from massive I/O buffers, and prevents connection pool depletion in concurrent network microservices.

### Pattern 14: Compiled Patterns, Precomputed Tables, & Singleflight

#### Package-Level Compiled Regular Expressions & Templates

Parsing regex patterns or template files on every function invocation is extremely expensive. Compile them once at package initialization.

```go
// Bad: Compiles the regex state machine on every function call (~5700ns)
func isValid(email string) bool {
    re := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
    return re.MatchString(email)
}

// Good: Compiles once at package startup; MatchString is safe for concurrent use (~450ns)
var emailRegex = regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)

func isValid(email string) bool {
    return emailRegex.MatchString(email)
}
```

#### Precomputed Lookup Tables

For pure computations with a small input space, replace the logic/branches with an array lookup. If it fits in L1/L2 cache, lookup is faster than computing.

```go
// Good: Lookup table avoids branching and conversions
var hexDigit = [16]byte{'0','1','2','3','4','5','6','7','8','9','a','b','c','d','e','f'}

func byteToHex(b byte) (byte, byte) {
    return hexDigit[b>>4], hexDigit[b&0x0f]
}
```

#### Deduplicating Work with `singleflight`

When a cache entry expires under high concurrent load, a cache stampede occurs where many goroutines fetch the same key simultaneously. Use `singleflight` to collapse duplicate concurrent fetches into a single execution.

```go
import "golang.org/x/sync/singleflight"

var (
    cache sync.Map
    sf    singleflight.Group
)

func GetWeather(city string) (string, error) {
    if val, ok := cache.Load(city); ok {
        return val.(string), nil
    }

    // Only one goroutine fetches; others block on the same key and receive the same result
    result, err, _ := sf.Do(city, func() (any, error) {
        data, err := fetchFromAPI(city)
        if err == nil {
            cache.Store(city, data)
        }
        return data, err
    })
    return result.(string), err
}
```

**Impact**: Prevents redundant CPU compilation work and database/network stampedes during cache misses.

### Pattern 15: Profile-Guided Optimization (PGO) & GOMAXPROCS in Containers

#### Profile-Guided Optimization (PGO)

Go 1.21+ supports PGO, letting the compiler optimize hot execution paths (devirtualizing interfaces and inlining hot functions) based on actual production CPU profiles. PGO typically yields a 2-7% CPU performance improvement.

1. Collect a CPU profile from your production service under representative load:

   ```bash
   curl http://localhost:6060/debug/pprof/profile?seconds=60 > cpu.pprof
   ```

2. Place the profile file as `default.pgo` in your main package directory.
3. Build the binary. The `go build` tool automatically detects `default.pgo` and applies PGO optimizations:

   ```bash
   go build ./cmd/myapp
   ```

#### GOMAXPROCS in Containers

By default, the Go runtime sets `GOMAXPROCS` to the host CPU core count. In a containerized environment (e.g., Kubernetes, Docker) with CPU limits, this causes the Go scheduler to spawn too many OS threads, leading to high context-switching overhead and CPU throttling.

- **Go 1.25+**: Automatically detects container CPU limits under cgroup v2, correctly setting `GOMAXPROCS`.
- **Go 1.24 and earlier**: Import the `automaxprocs` library to dynamically adjust `GOMAXPROCS` to the container CPU limits.

```go
import _ "go.uber.org/automaxprocs" // automatically configures GOMAXPROCS to match container quota

func main() {
    // runtime.GOMAXPROCS is now correctly set to the container's CPU quota
    startServer()
}
```

**Impact**: Automatically optimizes production binaries for real CPU workloads, and prevents scheduler thread bloat/throttling in containerized cloud services.

## Best Practices

1. **Profile before optimizing** — intuition about bottlenecks is wrong ~80% of the time
2. **Allocation reduction has the biggest ROI** — Go's GC is fast but not free
3. **One change at a time** — change, benchmark, compare, commit
4. **Use `slices` and `maps` standard packages** — typed, allocation-free, inlinable
5. **Preallocate when capacity is known** — avoid repeated slice growth
6. **Tune `http.Transport` for production** — default MaxIdleConnsPerHost=2 is too low
7. **Set `GOMEMLIMIT` in containers** — prevents OOM kills
8. **Document optimizations** — add code comments with benchmark numbers so readers know why a non-obvious pattern is used

## Common Pitfalls

- Optimizing without profiling first
- Using default `http.Client` in production without tuning the Transport
- Logging in hot loops (prevents inlining, allocates even when level is disabled)
- `panic`/`recover` as control flow (panic allocates a stack trace)
- Premature optimization — clarity first, then profile, then optimize
- `reflect.DeepEqual` in hot paths instead of `slices.Equal` / `maps.Equal`
- Ignoring field alignment — bloated structs waste memory and CPU cache
