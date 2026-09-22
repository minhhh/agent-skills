# Code Quality

Readable C++ is C++ where the intent is visible from the types and names, and where a reviewer does not have to reconstruct ownership or lifetime from memory.

## `const` correctness

Mark everything `const` that does not mutate. Const-correctness documents intent and prevents accidental writes.

```cpp
class Account {
public:
    // GOOD: read-only method, callable on a const Account
    double balance() const { return balance_; }

    // GOOD: mutating method, not callable on a const Account
    void deposit(double amount) { balance_ += amount; }

private:
    double balance_ = 0.0;
};
```

Pass read-only parameters by `const&`. Declare member functions `const` unless they mutate. Use `constexpr` for values known at compile time.

```cpp
constexpr double kPi = 3.141592653589793;
```

Do not use `mutable` to sidestep `const` except for genuine caches guarded by a mutex, and document why.

## Naming

Pick one convention per project and keep it. The point is consistency, not the specific choice.

| Entity | Typical convention |
| --- | --- |
| Type, class, struct, concept | `PascalCase` |
| Function, method, variable | `snake_case` or `camelCase` |
| Constant, compile-time value | `kPascalCase` or `UPPER_SNAKE_CASE` |
| Namespace | `lowercase` |
| Member variable | `snake_case_` or `m_` prefix |
| Template parameter | `T`, `U`, or a descriptive `ValueType` |

Names describe what the thing is, not its type. `user_count` beats `int_var`. `parse_config` beats `do_parse`.

## Headers

Headers are included everywhere, so they must be self-contained and minimal.

```cpp
// BAD: using-directive in a header leaks into every includer
using namespace std;

// GOOD: qualify, or bring in only what you need
#include <string>
void log(const std::string& message);
```

Use include guards or `#pragma once`, and put `#pragma once` at the top.

```cpp
#pragma once

#include <string>
#include <vector>
```

Include what you use. If a header uses `std::string`, include `<string>` even if another include currently provides it. Transitive includes break when the intermediate header changes.

Forward-declare when you only need a reference or pointer in the interface, to cut compile time and coupling.

```cpp
// In the header: forward declaration is enough for a reference parameter
class Widget;
void render(const Widget& widget);
```

Do not rely on implicit conversions in interfaces. Mark single-argument constructors `explicit` unless conversion is genuinely intended.

```cpp
// GOOD: no accidental implicit conversion from int to Widget
explicit Widget(int id);
```

## Definitions and the one-definition rule

Define non-template functions in a `.cpp`, not a header, unless they are `inline`. Header definitions included in several translation units violate the ODR and fail at link time.

```cpp
// BAD in a header: multiply defined across translation units
int counter() { return 42; }

// GOOD: inline or constexpr if it must live in the header
inline int counter() { return 42; }
constexpr int kCounter = 42;
```

## Formatting is not review material

Use `clang-format` and check the config into the repo. Formatting debates waste review time.

```bash
clang-format -i src/*.cpp include/*.hpp
```

## Static analysis

Run `clang-tidy` in CI with a project config. It catches bugs the compiler does not.

```bash
clang-tidy --config-file=.clang-tidy src/*.cpp
```

Suggested checks to start with:

```yaml
Checks: >
  bugprone-*,
  cppcoreguidelines-*,
  modernize-*,
  performance-*,
  readability-*
```

Do not enable every check at once. Turn on a category, fix the findings, then add the next.

## Functions and complexity

Keep functions short and single-purpose. If a function needs a comment to explain what it does, it is doing more than one thing.

- Prefer early return over deep nesting.
- Limit parameters. More than four suggests a config struct.
- Avoid boolean parameters. `render(true)` is unreadable; `render(WithBorder::yes)` is not.

```cpp
// BAD: what does true mean at the call site?
draw(widget, true, false);

// GOOD: name the options
enum class Border { none, shown };
draw(widget, Border::shown, Antialias::off);
```

## Comments explain why

The code says what it does. Comments explain the non-obvious reason, the workaround, or the constraint.

```cpp
// BAD: restates the code
i++;  // increment i

// GOOD: explains a decision the code cannot
// Reserve 64 slots: the vendor API silently drops writes beyond this.
buffer.reserve(64);
```

Keep comments current. A stale comment is worse than no comment.
