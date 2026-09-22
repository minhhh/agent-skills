# Build and Tooling

Half of C++ failures are not logic bugs. They are link errors, missing warnings, or a build that works on one machine and not another. These rules make the build deterministic and noisy in the right way.

## CMake, target-based

Modern CMake attaches properties to targets, not to global directories. Global settings leak across the whole project and break when dependencies change.

```cmake
# BAD: global settings affect every target, including dependencies
include_directories(${PROJECT_SOURCE_DIR}/include)
add_compile_options(-Wall -Wextra)

# GOOD: properties on the target that needs them
add_library(core src/core.cpp)
target_include_directories(core PUBLIC include)
target_compile_features(core PUBLIC cxx_std_20)
target_compile_options(core PRIVATE -Wall -Wextra -Wpedantic -Werror)
target_link_libraries(core PRIVATE fmt::fmt)
```

Public vs private matters. `PUBLIC` propagates to targets that link you. `PRIVATE` does not. Use `PRIVATE` for warnings so you do not force them on consumers.

```cmake
# Header is part of the interface: PUBLIC
target_include_directories(mylib PUBLIC include)

# Implementation-only dependency: PRIVATE
target_link_libraries(mylib PRIVATE spdlog::spdlog)

# The dependency's headers appear in mylib's public headers: PUBLIC
target_link_libraries(mylib PUBLIC fmt::fmt)
```

Use a `CMakePresets.json` so every machine uses the same generator, flags, and build directory.

```json
{
  "version": 6,
  "configurePresets": [
    {
      "name": "debug",
      "generator": "Ninja",
      "binaryDir": "${sourceDir}/build/debug",
      "cacheVariables": {
        "CMAKE_BUILD_TYPE": "Debug",
        "CMAKE_CXX_FLAGS": "-Wall -Wextra -Wpedantic -Werror"
      }
    }
  ]
}
```

Link with usage requirements, not by hand. `target_link_libraries(mylib PUBLIC fmt::fmt)` carries include directories, definitions, and link flags.

## Compiler warnings

Turn warnings on and treat them as errors. A warning that is ignored today is the bug that ships tomorrow.

```cmake
# GCC and Clang
-Wall -Wextra -Wpedantic -Werror
```

```cmake
# MSVC
/W4 /WX /permissive-
```

Add more warnings as the codebase allows: `-Wshadow`, `-Wconversion`, `-Wold-style-cast`, `-Wnon-virtual-dtor`, `-Woverloaded-virtual`. Do not enable `-Wconversion` on a large legacy build in one step; fix it in a new target first.

## Sanitizers

Sanitizers find at runtime what the type system cannot. Run them on every test in CI.

```cmake
# Address + undefined behavior, the default pair for tests
-fsanitize=address,undefined -fno-omit-frame-pointer -g
```

| Sanitizer | Finds | Flag |
| --- | --- | --- |
| AddressSanitizer | use-after-free, buffer overflow, leaks | `-fsanitize=address` |
| UndefinedBehaviorSanitizer | signed overflow, misaligned access, null deref | `-fsanitize=undefined` |
| ThreadSanitizer | data races | `-fsanitize=thread` |
| MemorySanitizer | uninitialized reads (Clang only) | `-fsanitize=memory` |

AddressSanitizer and ThreadSanitizer cannot run together. Run separate CI jobs.

Be aware of cost. AddressSanitizer slows execution roughly 2x and memory use more. That is why it runs in CI and tests, not in the release build.

## The one-definition rule

Every function, type, and variable has exactly one definition across the whole program, with exceptions for `inline` functions and templates.

```cpp
// BAD in a header: defined once per translation unit that includes it
int registry_size() { return 42; }

// GOOD: inline lets multiple definitions coexist, the linker picks one
inline int registry_size() { return 42; }
```

Symptoms of an ODR violation: duplicate symbol link errors, or worse, different translation units behaving differently with no error at all. The latter happens when the same type or inline function has different definitions in different files. Fix by making the definition identical and visible in one header.

## ABI stability

Do not change the layout of a type that crosses a binary boundary. Adding a member, reordering members, changing a virtual function, or changing an inline function's body can break callers compiled against the old version.

Rules for a shared library with a stable ABI:

- Do not add or reorder data members.
- Do not add, remove, or reorder virtual functions.
- Do not change the signature of an exported function.
- Keep inline functions and templates out of the exported interface, or accept that all users must recompile.

Use the pimpl idiom to hide implementation changes:

```cpp
// header stable across implementation changes
class Engine {
public:
    Engine();
    ~Engine();
    Engine(const Engine&) = delete;
    Engine& operator=(const Engine&) = delete;
    void run();
private:
    struct Impl;
    std::unique_ptr<Impl> impl_;
};
```

## Link errors and the include order

If the linker cannot find a symbol, check in order:

1. Is the definition compiled into some object file or library?
2. Is that library in `target_link_libraries` for this target?
3. For a static library, does the link order matter? Put dependents before dependencies.
4. Is the symbol wrapped in an anonymous namespace or marked `static`, hiding it from other translation units?

If the include is not found, check `target_include_directories` for the target and whether the include path was exported as `PUBLIC` or `INTERFACE`.

## Precompiled headers and unity builds

These speed up compilation but hide missing includes, because one translation unit's includes leak into another. Do not use a unity build to paper over missing includes. Fix the includes first.
