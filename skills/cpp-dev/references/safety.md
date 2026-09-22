# Safety

C++ does not stop you from reading freed memory, indexing out of bounds, or overflowing a signed integer. The compiler is allowed to assume none of those happen, so the bugs are silent and the behavior is undefined. This file covers the classes of UB that reach production most often.

Priority: Safety is first. A safety violation is never traded for speed or brevity.

## Initialize everything

Reading an indeterminate value is undefined behavior.

```cpp
// BAD: value is indeterminate; the read is UB
int value;
if (flag) value = compute();
use(value);

// GOOD: initialize at the point of declaration
int value = flag ? compute() : 0;
use(value);
```

The same rule applies to members. Give every member a default, either inline or in the constructor's initializer list.

```cpp
// BAD: count and name are indeterminate if the constructor body returns early
struct Counter {
    int count;
    std::string name;
};

// GOOD: members have defaults
struct Counter {
    int count = 0;
    std::string name;
};
```

Prefer brace initialization to stop narrowing conversions:

```cpp
double ratio = 0.5;
int truncated{ratio};  // compile error: narrowing, this is what you want
int allowed = static_cast<int>(ratio);  // explicit, deliberate
```

## Lifetime of views and references

`std::string_view`, `std::span`, and references do not own. They are valid only while the owner lives.

```cpp
// BAD: returns a view into a local that dies at the return
std::string_view filename() {
    std::string path = "/tmp/data.txt";
    return path;  // dangling
}

// GOOD: own the result
std::string filename() { return "/tmp/data.txt"; }

// GOOD: accept a view as input when the caller owns the buffer
void log(std::string_view message);
```

A `span` over a temporary container is the same bug:

```cpp
// BAD: the vector dies at the semicolon, the span points at freed memory
std::span<const int> view = get_values();
int total = std::accumulate(view.begin(), view.end(), 0);

// GOOD: keep the owning vector alive alongside the view
auto values = get_values();
std::span<const int> view{values};
```

## Bounds

`operator[]` does not check. Reach for bounds-checked access when the index is not provably valid.

```cpp
// BAD: index comes from input, out-of-range is UB
auto item = items[user_index];

// GOOD: throws std::out_of_range, or use a span with .at semantics
auto item = items.at(user_index);
```

For raw buffers, carry the size and use `std::span` so the bound travels with the pointer.

```cpp
void fill(std::span<int> buffer, int value) {
    std::fill(buffer.begin(), buffer.end(), value);  // size is known
}
```

## Integer overflow and signedness

Signed overflow is undefined behavior. Unsigned overflow wraps, which is often wanted but is a common logic bug.

```cpp
// BAD: signed overflow, UB
int total = a + b;

// GOOD: check before it wraps
if (b > 0 && a > std::numeric_limits<int>::max() - b) {
    return std::unexpected(Overflow{});
}
int total = a + b;
```

Prefer `<cstdint>` fixed-width types for protocols and file formats. Prefer `std::ssize` over `.size()` when you need a signed size, since `.size()` returns unsigned and signed/unsigned comparison is a bug factory.

## Raw allocation

Ownership by raw pointer is the source of leaks and double frees.

```cpp
// BAD: leaks if throw_before_owner runs
Widget* w = new Widget();
throw_before_owner();
std::unique_ptr<Widget> owner{w};

// GOOD: allocation creates the owner
auto owner = std::make_unique<Widget>();
```

Never construct two owners from the same raw pointer. This is the classic double free.

```cpp
// BAD: two unique_ptrs own one object, double free at scope exit
auto* raw = new Widget();
std::unique_ptr<Widget> a{raw};
std::unique_ptr<Widget> b{raw};
```

## C string functions

`strcpy`, `strcat`, `sprintf`, `gets`, and `memcpy` without a size are buffer overflows waiting for input.

```cpp
// BAD: no bound on the copy
char buffer[64];
strcpy(buffer, user_input);

// GOOD: owned, allocated, sized
std::string buffer = user_input;

// GOOD: formatted without a buffer
auto message = std::format("user {} logged in", user_id);
```

For binary copies, prefer `std::copy` with iterators, or `std::span` and `memcpy(dst.data(), src.data(), dst.size_bytes())` only when the sizes are already equal and known.

## Reinterpretation

`reinterpret_cast` and C-style casts silence the type system without a runtime check.

```cpp
// BAD: if the object is not actually a Derived, this is UB
auto* d = (Derived*)base_ptr;
auto* d = reinterpret_cast<Derived*>(base_ptr);

// GOOD: checked downcast, returns nullptr on mismatch
auto* d = dynamic_cast<Derived*>(base_ptr);
```

Reserve `reinterpret_cast` for the narrow cases where it is required, such as talking to a C API or hardware, and isolate it in one place with a comment explaining the guarantee.
