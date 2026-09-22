# Common Pitfalls: These Thoughts Mean STOP

If you catch yourself thinking any of these, stop and apply the correct approach.

| Rationalization | Problem | Impact | Fix |
| --- | --- | --- | --- |
| "A raw pointer is faster here" | Ownership is now unclear, and the fast path is rarely the issue | Leak or double free that appears in production | `unique_ptr`, or a non-owning reference if ownership lives elsewhere |
| "I tested it manually" | Manual testing does not cover lifetime, races, or UB | Silent corruption under load | Write a test, run under AddressSanitizer |
| "The compiler did not warn" | Warnings are off or the bug is in the standard's blind spot | UB that only a sanitizer catches | `-Wall -Wextra -Werror`, run sanitizers in CI |
| "It works on my machine" | Undefined behavior, ABI, or a missing include | Fails on another compiler or optimization level | Build with two compilers and `-O2`, add the missing include |
| "An atomic flag makes it thread-safe" | An atomic protects one value, not compound state | Race that passes every test on x86 | Guard compound state with a mutex |
| "I will initialize it later" | Reading before the assignment is UB | Nondeterministic behavior | Initialize at declaration |
| "The `string_view` is fine, the string is still around" | The view does not extend the owner's lifetime | Use-after-free | Own the string, or prove the owner outlives the view |
| "`shared_ptr` everywhere avoids reasoning about ownership" | Refcounts, cycles, and ownership blur | Leaks from cycles, atomic overhead, unclear design | `unique_ptr` by default, `shared_ptr` only for real sharing |
| "It is only a small copy" | Copies of containers and strings are O(n) | Unexpected slowdown in hot paths | Pass `const&` or a view, or move |
| "`std::move` makes it faster" | `move` on a `const` object copies, `move` on a return value blocks elision | Slower, not faster | Move only non-const locals you are done with |
| "`noexcept` is just a hint" | A throw from `noexcept` calls `std::terminate` | Process death instead of an exception | Mark `noexcept` only when it is true |
| "Catch `...` and continue" | Swallows real bugs and leaves state unknown | Corruption with no signal | Catch specific exceptions, log, rethrow |
| "`volatile` handles threading" | `volatile` is not atomic and does not order memory | Data race | `std::atomic` or a mutex |
| "`detach` the thread, it will finish" | Detached threads outlive the objects they use | Use-after-free at shutdown | `std::jthread`, or join explicitly |
| "I will template it in case it is reused" | One instantiation, now a template error surface | Slower build, harder debug | Concrete type until the second user appears |
| "Constrain later with `enable_if`" | SFINAE errors are unreadable | Hours lost to template diagnostics | A concept from the start |
| "The header is self-contained, it compiled" | Missing includes hidden by transitive ones | Breaks when the intermediate header changes | Include what you use |

## Red flags in code review

These are not opinions. Each one has caused a production incident.

- `new` or `delete` anywhere outside a custom resource wrapper
- A raw owning pointer, or a parameter that is a `T*` where ownership is unclear
- `strcpy`, `strcat`, `sprintf`, `memcpy` without a size
- `reinterpret_cast` or a C-style cast
- A `lock()` without a matching RAII guard
- A `detach()`ed thread
- A `catch (...)` that logs and does not rethrow
- `using namespace` in a header
- A non-`explicit` single-argument constructor
- A destructor that can throw
- A struct passed across a library boundary with public members
- An uninitialized member
