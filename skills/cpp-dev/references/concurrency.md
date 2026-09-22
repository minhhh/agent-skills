# Concurrency

A data race is undefined behavior, not a slow path. Two threads touching the same memory with at least one write and no synchronization is a bug the compiler is allowed to exploit. Symptoms show up under load, once in millions of requests.

## Threads manage their own lifetime

A `std::thread` that is destroyed without `join()` or `detach()` calls `std::terminate`. Prefer `std::jthread`, which joins on destruction and carries a stop token.

```cpp
// BAD: if this throws before join(), the program terminates
std::thread worker{do_work};
do_more();
worker.join();

// GOOD: jthread joins automatically on every exit path
std::jthread worker{[](std::stop_token st) {
    while (!st.stop_requested()) {
        step();
    }
}};
```

Do not detach threads. A detached thread can outlive the objects it references.

## Stop tokens

Pass work a `std::stop_token` so it can stop, and request a stop from the owner.

```cpp
// The callable receives the token automatically when its first parameter is stop_token
std::jthread worker{[](std::stop_token st, Queue& q) {
    while (!st.stop_requested()) {
        process(q.pop());
    }
}, std::ref(queue)};

// Ask it to stop
worker.request_stop();
// jthread destructor joins here
```

## Locks own themselves

Never call `lock()` and `unlock()` by hand. An exception or early return between them deadlocks.

```cpp
// BAD: throw in do_work leaves the mutex locked
mu.lock();
do_work();
mu.unlock();

// GOOD: RAII release on every path
{
    std::scoped_lock lock{mu};
    do_work();
}
```

Use `std::scoped_lock` for one or more mutexes. It uses a deadlock-avoiding algorithm when given several.

```cpp
// GOOD: lock both without deadlock, regardless of order at other call sites
std::scoped_lock lock{a, b};
```

For a `shared_mutex`, read with `std::shared_lock` and write with `std::unique_lock`.

```cpp
std::shared_mutex config_mutex;

Config read_config() {
    std::shared_lock lock{config_mutex};  // many readers
    return config_;
}

void write_config(Config next) {
    std::unique_lock lock{config_mutex};  // one writer
    config_ = std::move(next);
}
```

## One lock order

Even with `scoped_lock`, a manual lock path elsewhere can deadlock. Pick one global order for acquiring multiple mutexes and document it. Never acquire B then A if any path acquires A then B.

## Atomics are for flags and counters

`std::atomic` gives you an indivisible operation on one value. It does not protect a compound structure.

```cpp
// BAD: the atomic flag does not make the map operations safe
std::atomic<bool> updating{false};
if (!updating.exchange(true)) {
    cache[key] = value;  // still a race with readers
    updating.store(false);
}

// GOOD: if two threads touch it, guard it with a mutex
std::mutex cache_mutex;
{
    std::scoped_lock lock{cache_mutex};
    cache[key] = value;
}
```

Atomic counters and flags are fine:

```cpp
std::atomic<uint64_t> requests_served{0};
requests_served.fetch_add(1, std::memory_order_relaxed);

std::atomic<bool> shutting_down{false};
if (shutting_down.load(std::memory_order_acquire)) return;
```

Default to `std::memory_order_seq_cst`. Loosen the ordering only with a documented reason, because a wrong ordering is a bug that passes every test on x86 and fails on ARM.

## Condition variables

Pair a condition variable with a mutex and a predicate. Always wait in a loop, because spurious wakeups are allowed.

```cpp
// BAD: one wait, spurious wakeup proceeds with nothing to consume
cv.wait(lock);

// GOOD: wait until the predicate holds
cv.wait(lock, [&] { return !queue.empty() || stopped; });
if (stopped) return;
auto item = queue.front();
queue.pop();
```

Change the predicate under the same mutex, then notify:

```cpp
{
    std::scoped_lock lock{mu};
    queue.push(item);
}
cv.notify_one();
```

## Do not share more than you must

The best fix for a race is to remove the sharing. Prefer thread-local state, message passing, or an immutable snapshot.

```cpp
// GOOD: each thread has its own scratch buffer
thread_local std::vector<char> scratch;

// GOOD: publish an immutable snapshot behind shared_ptr, readers never lock
std::shared_ptr<const Config> current_config;
```

## Checklist

- Is every shared mutable variable guarded by a mutex, an atomic, or made thread-local?
- Does every thread join, either explicitly or through `jthread`?
- Is every lock RAII, never manual?
- Is there a single documented lock order?
- Is shared state read while it could be written?
- Is the code free of `volatile` used as synchronization? `volatile` is not atomic and does not order memory.
