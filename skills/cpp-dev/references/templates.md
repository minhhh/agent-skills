# Templates and Generic Programming

Templates are compile-time contracts. The goal is a constraint that fails at the call site with a readable reason, not an error 40 frames deep inside a library body.

## Constrain with concepts

`requires` replaced almost all SFINAE in C++20. A concept names the requirement, so the error message names it too.

```cpp
// BAD: any T compiles until it fails deep in the body
template <typename T>
T max_of(T a, T b) { return a < b ? b : a; }

// GOOD: the constraint fails at the call site
template <typename T>
    requires std::totally_ordered<T>
T max_of(T a, T b) { return a < b ? b : a; }
```

Define a concept when the requirement is a reusable idea:

```cpp
template <typename T>
concept Numeric = std::integral<T> || std::floating_point<T>;

template <Numeric T>
T clamp(T value, T lo, T hi) {
    if (value < lo) return lo;
    if (value < hi) return value;
    return hi;
}
```

Prefer standard concepts from `<concepts>`: `std::integral`, `std::floating_point`, `std::totally_ordered`, `std::invocable`, `std::ranges::range`, `std::copyable`, `std::movable`.

## Shorthand and requires-expressions

```cpp
// Shorthand: a single concept before the parameter list
void advance(std::integral auto& counter);

// requires-expression: assert a set of expressions is valid
template <typename T>
concept StringLike = requires(T t) {
    { t.size() } -> std::convertible_to<std::size_t>;
    { t.data() } -> std::convertible_to<const char*>;
};

// requires-clause with a compound condition
template <typename T>
    requires std::ranges::range<T> && std::copyable<std::ranges::range_value_t<T>>
void store(const T& values);
```

## Avoid SFINAE in new code

`std::enable_if` and `std::void_t` still appear in older code and in template scaffolding, but new interfaces should use concepts.

```cpp
// OLD: SFINAE, unreadable errors
template <typename T,
          typename = std::enable_if_t<std::is_integral_v<T>>>
T add(T a, T b);

// NEW: concept
template <std::integral T>
T add(T a, T b);
```

If you must read SFINAE, the rule is: enable the overload only when the substitution is valid, so a bad match is removed from the candidate set rather than producing a hard error.

## `constexpr` and `consteval`

Move work to compile time when the inputs are known then. Compile-time checks are strictly better than runtime checks.

```cpp
// GOOD: evaluated at compile time when possible
constexpr int factorial(int n) {
    return n <= 1 ? 1 : n * factorial(n - 1);
}

constexpr int answer = factorial(5);  // computed by the compiler
```

Use `consteval` when a function must run only at compile time:

```cpp
consteval std::size_t length(const char* s) {
    return *s ? 1 + length(s + 1) : 0;
}
```

`if constexpr` selects a branch at compile time and discards the other, which is how you specialize without tag dispatch:

```cpp
template <typename T>
std::string describe(T value) {
    if constexpr (std::integral<T>) {
        return std::format("int {}", value);
    } else if constexpr (std::floating_point<T>) {
        return std::format("float {}", value);
    } else {
        return "other";
    }
}
```

Remember the discarded branch must still parse. Errors in the not-taken branch can occur if they depend on template parameters.

## Non-type template parameters

Use them for sizes and compile-time configuration so the value is part of the type.

```cpp
template <std::size_t N>
struct FixedBuffer {
    std::array<std::byte, N> data{};
};
```

## Keep template error messages short

The failure mode of templates is a wall of instantiation backtrace. The fixes:

- Constrain with concepts so the error names the requirement.
- Add `static_assert` with a message for invariants the concept cannot express.
- Prefer concrete types at the interface and templates in the implementation.
- Give the type meaningful names.

```cpp
template <typename T>
void process(T value) {
    static_assert(Numeric<T>, "process requires a numeric type");
    // ...
}
```

## Do not over-generalize

Do not template code that only ever has one instantiation. A concrete type is easier to read, faster to compile, and easier to debug.

```cpp
// BAD: a template with exactly one user
template <typename T>
class UserRepository { /* ... */ };

// GOOD: a concrete class
class UserRepository { /* ... */ };
```

Generalize when the second use appears, not before.
