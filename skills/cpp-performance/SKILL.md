---
name: cpp-performance
description: Use when debugging slow C++ code, profiling performance, reducing latency or memory use, or optimizing CPU-bound loops, allocations, cache behavior, and build-time optimization such as PGO and LTO.
---

# C++ Performance Optimization

Profile first, then optimize the hotspot the profiler found. Intuition about C++ performance is wrong often enough that guessing wastes more time than it saves. Measure, change one thing, measure again.

## When to Use

- A profiler has identified a hot function or a memory bottleneck
- Reducing latency, CPU, or memory in a library or server
- Choosing between data layouts, allocation strategies, or dispatch mechanisms
- Setting up PGO, LTO, or target-specific flags for a release build
- Investigating cache misses or branch mispredictions

Do not use this to optimize a cold path. Startup code, config parsing, and error paths rarely matter.

## Decision Tree: Where Is the Time Going?

| Bottleneck | Signal | Pattern |
| --- | --- | --- |
| Too many allocations | heap profile dominated by `operator new`, `malloc` | Reserve, pool, arena, avoid temporaries |
| Copying large objects | `memcpy` high in the profile | Pass by `const&`, move, return by value |
| Cache misses | high `cache-misses`, low IPC in `perf stat` | SoA layout, smaller working set, iterate linearly |
| Branch mispredicts | high `branch-misses` | Branchless code, sort before branch, table lookup |
| Virtual dispatch | indirect call overhead in hot loop | Devirtualize, template, `final`, variant + visit |
| Lock contention | mutex in the profile under threads | Shrink the critical section, shard, atomic fast path |
| Repeated expensive work | same value computed or fetched again | Cache, precompute, memoize |
| Wrong algorithm | O(n^2) where O(n log n) exists | Fix the algorithm before micro-optimizing |
| Compiler not optimizing | assembly shows redundant work | Inline boundary headers, PGO, LTO, flags |

Rule out the algorithm and the I/O first. If 90 percent of the time is a database query, no amount of cache tuning helps.

## Profiling Tools

### perf

The default Linux profiler. CPU sampling with call graphs, no instrumentation.

```bash
# Record a profile for the process
perf record -g --call-graph dwarf ./app

# Report the hot functions
perf report

# Hardware counters: IPC, cache, and branch behavior
perf stat -e cycles,instructions,cache-misses,branch-misses ./app
```

Read `perf stat` before `perf report`. A low IPC with high `cache-misses` points at the memory hierarchy, not at a bad algorithm.

### Flamegraphs

Turn a `perf` capture into a flamegraph to see where time pools across the call stack.

```bash
perf record -F 99 -g --call-graph dwarf ./app
perf script > out.perf
stackcollapse-perf.pl out.perf > out.folded
flamegraph.pl out.folded > flame.svg
```

### Valgrind and callgrind

Deterministic instruction-level profiling. Slow, so use it on a reduced workload.

```bash
valgrind --tool=callgrind --callgrind-out-file=callgrind.out ./app
callgrind_annotate callgrind.out
```

Use `cachegrind` for cache behavior and `--tool=massif` for heap growth over time.

### heaptrack

Heap profiling with allocation stacks. Finds the allocation that grows without bound.

```bash
heaptrack ./app
heaptrack_gui heaptrack.app.*.gz
```

### Intel VTune and AMD uProf

Hardware-level analysis: microarchitecture, memory access, thread concurrency. Use when `perf` is not granular enough.

### Microbenchmarks

For a single function, use Google Benchmark with a stable machine.

```cpp
#include <benchmark/benchmark.h>

static void BM_Parse(benchmark::State& state) {
    auto input = make_input();
    for (auto _ : state) {
        benchmark::DoNotOptimize(parse(input));
    }
}
BENCHMARK(BM_Parse);
BENCHMARK_MAIN();
```

```bash
./bench --benchmark_format=json > before.json
# change one thing
./bench --benchmark_format=json > after.json
```

Compare with a statistical summary, not a single run. Report the delta and the variance.

## Optimization Patterns

### Pattern 1: Reserve container capacity

Growing a vector reallocates and copies. Reserve when the size is known.

```cpp
// BAD: multiple reallocations as it grows
std::vector<Result> results;
for (const auto& item : items) {
    results.push_back(process(item));
}

// GOOD: one allocation
std::vector<Result> results;
results.reserve(items.size());
for (const auto& item : items) {
    results.push_back(process(item));
}
```

### Pattern 2: Pass and return without copies

```cpp
// BAD: copies the string in and out
std::string normalize(std::string input);

// GOOD: view in, value out
std::string normalize(std::string_view input);
```

Return by value and let the compiler elide the copy. Do not use an out-parameter to "avoid a copy"; it usually adds one and hurts readability.

```cpp
// GOOD: named return value optimization, no copy
std::vector<int> load() {
    std::vector<int> out;
    // fill
    return out;
}
```

### Pattern 3: Move instead of copy

```cpp
// BAD: copies the buffer
sink.push_back(buffer);

// GOOD: moves it
sink.push_back(std::move(buffer));
```

Mark move operations `noexcept` so `std::vector` moves rather than copies during growth.

```cpp
Buffer(Buffer&&) noexcept;
Buffer& operator=(Buffer&&) noexcept;
```

### Pattern 4: Avoid temporaries in loops

```cpp
// BAD: constructs a string every iteration
for (int i = 0; i < n; ++i) {
    log("value " + std::to_string(i));
}

// GOOD: reuse a buffer
std::string line;
line.reserve(32);
for (int i = 0; i < n; ++i) {
    line.clear();
    std::format_to(std::back_inserter(line), "value {}", i);
    log(line);
}
```

### Pattern 5: Layout for the cache

Data that is used together should live together. Loop over contiguous memory.

```cpp
// AoS: each iteration pulls a full Record, most of it unused
struct Particle { float x, y, z, mass; };
std::vector<Particle> particles;
for (auto& p : particles) total += p.x;  // touches y, z, mass too

// SoA: each iteration touches only the field in use
struct Particles {
    std::vector<float> x, y, z, mass;
};
for (float v : particles.x) total += v;
```

Use SoA when the hot loop touches one or two fields and the count is large. Use AoS when the loop uses most fields together, since AoS keeps one object's fields on the same cache line.

### Pattern 6: Branchless where it matters

A mispredicted branch costs far more than a few arithmetic operations.

```cpp
// BAD: unpredictable branch in a hot loop
if (value < 0) value = 0;

// GOOD: branchless clamp
value = value < 0 ? 0 : value;  // compilers often make this a cmov

// GOOD: arithmetic form when the compiler will not
value = std::max(value, 0);
```

Sorting data before a branch-heavy pass can also help, by making the branch predictable.

### Pattern 7: Reduce virtual dispatch in hot paths

A virtual call in a tight loop blocks inlining and prediction.

```cpp
// BAD: virtual call per element
for (const auto& shape : shapes) total += shape->area();

// GOOD: one virtual call, then a switch that inlines
for (const auto& shape : shapes) {
    switch (shape->kind()) {
        case Kind::circle:  total += static_cast<const Circle&>(*shape).area(); break;
        case Kind::square:  total += static_cast<const Square&>(*shape).area(); break;
    }
}
```

Mark leaf classes `final` so the compiler can devirtualize. Prefer `std::variant` plus `std::visit` for a closed set of types.

### Pattern 8: Custom allocators and pools

If the profile shows `operator new` at the top, reuse memory.

```cpp
// Arena for short-lived nodes, freed all at once
class Arena {
public:
    void* allocate(std::size_t n);
    void reset();  // release everything
private:
    std::vector<std::unique_ptr<std::byte[]>> blocks_;
};
```

For a fixed-size type allocated and freed often, a free list beats the general allocator.

```cpp
// Thread-local pool avoids cross-thread contention
thread_local std::vector<std::unique_ptr<Node>> node_pool;
```

### Pattern 9: Contiguity and small types

Store hot, small values by value in a contiguous container. An index into a vector is cheaper than a pointer chase and keeps the data close.

```cpp
// BAD: pointer per element, one cache miss each
std::vector<std::unique_ptr<Entity>> entities;

// GOOD: values in one block, iterate linearly
std::vector<Entity> entities;
```

When you must reference across a vector, store indices, not pointers, and accept the reallocation risk deliberately.

### Pattern 10: Compile-time work

Move computation the compiler can do to compile time.

```cpp
// GOOD: computed at compile time
constexpr auto table = build_table();
```

`if constexpr` removes dead branches from the instantiation, and `consteval` guarantees compile-time evaluation.

### Pattern 11: Inline small hot functions

The compiler inlines within a translation unit. Put small, hot functions in the header, or rely on LTO to inline across units.

```cpp
// Header, small, hot
inline int clamp(int value, int lo, int hi) {
    return value < lo ? lo : (value > hi ? hi : value);
}
```

Do not inline large functions. It grows the binary and can hurt the instruction cache.

### Pattern 12: SIMD

Vectorize independent arithmetic. Sometimes the compiler does it with `-O3 -march=native`; sometimes it needs a hint or intrinsics.

```cpp
// Good codegen candidate: independent, contiguous, no aliasing
void scale(float* data, std::size_t n, float factor) {
    for (std::size_t i = 0; i < n; ++i) data[i] *= factor;
}
```

`std::execution::par_unseq` policies can vectorize algorithms. For explicit control, use intrinsics or a portable library and measure.

### Pattern 13: Link-time optimization and PGO

LTO inlines across translation units. PGO tells the compiler which branches are hot.

```bash
# LTO
cmake -DCMAKE_INTERPROCEDURAL_OPTIMIZATION=ON ..

# PGO: build instrumented, run the workload, rebuild with the profile
clang++ -fprofile-generate -o app_pgo app.cpp
./app_pgo < representative-workload
clang++ -fprofile-use -o app_opt app.cpp
```

PGO typically gains 5 to 15 percent on branchy code. Use a representative workload, not a trivial one.

### Pattern 14: Tune compiler flags deliberately

```cmake
# Release defaults
-O2 -DNDEBUG
```

`-O3` can increase code size and hurt the instruction cache. `-march=native` improves vectorization but ties the binary to a CPU; for distribution use a baseline target like `-march=x86-64-v3`.

## Best Practices

- Profile before and after every change. Keep the benchmark harness.
- Change one thing at a time and re-measure.
- Optimize the algorithm before the constant factor.
- Prefer clearer code when the measured difference is noise. Maintainability is also a cost.
- Benchmark on a machine that represents production. A laptop result can mislead.
- Watch code size. Faster per call but larger can still be slower overall through the instruction cache.

## Common Pitfalls

| Mistake | Why it is wrong | Fix |
| --- | --- | --- |
| Optimizing without a profile | Guessing picks the wrong hotspot | `perf record`, then optimize what it shows |
| Benchmarking with optimization off | Measures debug code, not release | Build with `-O2` and `-DNDEBUG` |
| One benchmark run | Noise dominates the result | Repeat and report the distribution |
| Micro-optimizing a cold path | Complexity with no measurable gain | Only optimize what the profile shows |
| Preferring `shared_ptr` for speed of design | Atomic refcounts and cache misses | `unique_ptr` and references |
| `-O3` everywhere by reflex | Code bloat, instruction-cache misses | Bench `-O2` against `-O3` |
| `-march=native` for a distributed binary | Illegal instruction on older CPUs | Use a baseline target |
| Ignoring `noexcept` on moves | Vector growth copies instead of moves | Mark moves and swaps `noexcept` |
| Using `std::endl` in a loop | Flushes on every call | Use `'\n'`, flush deliberately |
| Chasing micro-wins before fixing layout | Cache misses dwarf instruction counts | Fix the data layout first |
