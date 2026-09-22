# Error Handling

C++ has three mechanisms for failure: exceptions, `std::expected` and `std::optional`, and error codes. Pick one style per boundary and stay consistent. Mixing them inconsistently is how errors get dropped.

## Choosing a mechanism

| Situation | Use | Why |
| --- | --- | --- |
| Expected, local, part of normal flow | `std::expected<T, E>` | No stack unwinding, explicit at the call site |
| Absence is normal | `std::optional<T>` | Carries "no value" without an error type |
| Constructor failure | exception | Constructors cannot return a value |
| Operator overloading | exception | Operators cannot carry an error return |
| Non-recoverable invariant | `std::terminate` / `assert` | Continuing is worse than stopping |
| Across a C or ABI boundary | error code | Exceptions must not cross the boundary |

## `expected` for fallible operations

`std::expected` carries either a value or an error. It is C++23.

```cpp
enum class ParseError { Empty, BadDigit };

std::expected<int, ParseError> parse_int(std::string_view in) {
    if (in.empty()) return std::unexpected(ParseError::Empty);
    int value = 0;
    for (char c : in) {
        if (c < '0' || c > '9') return std::unexpected(ParseError::BadDigit);
        value = value * 10 + (c - '0');
    }
    return value;
}

// Caller must handle both
auto result = parse_int(input);
if (!result) {
    return std::unexpected(result.error());
}
use(*result);
```

Use `std::optional` when there is no error to explain, only absence.

```cpp
std::optional<User> find_user(UserId id);
```

Do not return a `bool` plus an out-parameter. The caller loses the reason and the value at the same time.

```cpp
// BAD: why did it fail? what is in cfg on failure?
bool load(const std::string& path, Config& cfg);

// GOOD
std::expected<Config, LoadError> load(const std::string& path);
```

## Exception-safety guarantees

Every function that can throw gives one of three guarantees. Know which one you provide.

| Guarantee | Promise |
| --- | --- |
| No-throw (`noexcept`) | Never throws |
| Strong | On failure, the object is unchanged, the operation failed atomically |
| Basic | On failure, the object is valid but unspecified, no leaks |
| None | Do not offer this |

Prefer the strong guarantee when it is cheap. Compute everything that can fail first, then commit with non-throwing operations.

```cpp
// GOOD: build the new value first, swap in at the end
void update(Config next) {
    validate(next);        // may throw, object untouched
    config_ = std::move(next);  // noexcept commit
}
```

The copy-and-swap idiom gives copy assignment the strong guarantee:

```cpp
class Buffer {
public:
    Buffer& operator=(const Buffer& other) {
        Buffer tmp{other};       // may throw, *this untouched
        swap(tmp);               // noexcept
        return *this;
    }
private:
    void swap(Buffer& other) noexcept { std::swap(raw_, other.raw_); }
    char* raw_ = nullptr;
};
```

## `noexcept`

Mark `noexcept` when a function truly cannot throw. This matters for moves and swaps.

```cpp
// GOOD: vector uses move instead of copy when growing
Buffer(Buffer&& other) noexcept;
Buffer& operator=(Buffer&& other) noexcept;
void swap(Buffer& other) noexcept;
```

Do not mark `noexcept` on a function that can throw. A throw from a `noexcept` function calls `std::terminate`. Destructors are implicitly `noexcept`; do not let one throw.

```cpp
// BAD: a throwing destructor terminates the program
~File() noexcept(false) { flush(); }  // if flush throws, terminate

// GOOD: swallow or log in the destructor, or provide an explicit close()
~File() { try { flush(); } catch (...) { log_error(); } }
```

## Catch by reference, catch narrowly

```cpp
// BAD: catches by value, slices the exception
catch (std::exception e) { }

// BAD: swallows everything, including programming errors
catch (...) { log("failed"); }

// GOOD: catch by const reference, name what you can handle
catch (const std::filesystem::filesystem_error& e) {
    log_error("fs error: {}", e.what());
    throw;
}

// GOOD: catch-all at a boundary, then rethrow or convert deliberately
catch (const std::exception& e) {
    log_error("unhandled: {}", e.what());
    throw;
}
```

Never `throw;` a different exception inside a destructor, and never let an exception escape a destructor.

## Do not use exceptions for control flow

Exceptions are for exceptional, non-local failure. Using them for expected conditions wastes the fast path and hides the control flow.

```cpp
// BAD: "not found" is normal, not exceptional
try {
    return registry.at(key);
} catch (const std::out_of_range&) {
    return default_value;
}

// GOOD: check first
if (auto it = registry.find(key); it != registry.end()) {
    return it->second;
}
return default_value;
```

## Custom error types carry context

A bare enum loses the detail that makes the error actionable.

```cpp
// GOOD: an error type with the context the caller needs
struct LoadError {
    std::string path;
    int line;
    std::string reason;
};
```

Keep error types cheap to copy and free of side effects. Format them for the user at the boundary, log the detail internally.
