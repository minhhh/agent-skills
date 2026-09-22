# Memory and Ownership

Every resource in C++ needs exactly one owner and a deterministic release point. RAII is that mechanism: the destructor runs when the owner goes out of scope, on every exit path.

## RAII for every resource

Wrap each resource in a type whose destructor releases it. Then you never write a manual cleanup path.

```cpp
// BAD: leaks the handle on any early return or exception
Handle h = open_handle(path);
if (!validate(h)) return error;
use(h);
close_handle(h);

// GOOD: the destructor runs on every path
struct Handle {
    explicit Handle(const std::string& path) : h_{open_handle(path)} {}
    ~Handle() { close_handle(h_); }
    Handle(const Handle&) = delete;
    Handle& operator=(const Handle&) = delete;
    Handle(Handle&&) noexcept = default;
    Handle& operator=(Handle&&) noexcept = default;
private:
    RawHandle h_;
};
```

This applies to file descriptors, sockets, locks, database connections, GPU buffers, and temporary state changes. The standard library already wraps many of these: `std::fstream`, `std::jthread`, `std::scoped_lock`, `std::unique_ptr`.

## Rule of zero

If your type only holds members that already manage themselves, declare no destructor, no copy, and no move. The compiler generates correct ones.

```cpp
// GOOD: the vector and string manage their own resources
struct User {
    std::string name;
    std::vector<std::string> roles;
};
```

Reaching for the rule of zero keeps the type trivially movable and copyable, which containers and algorithms rely on.

## Rule of five

If you declare any one of destructor, copy constructor, copy assignment, move constructor, or move assignment, declare all five. Declaring only some leaves the others to the compiler, which produces surprising or deleted behavior.

```cpp
class Buffer {
public:
    Buffer() = default;
    ~Buffer() { /* frees raw_ */ }

    Buffer(const Buffer& other);             // copy
    Buffer& operator=(const Buffer& other);  // copy assign
    Buffer(Buffer&& other) noexcept;         // move
    Buffer& operator=(Buffer&& other) noexcept;
private:
    char* raw_ = nullptr;
};
```

Mark moves `noexcept`. `std::vector` moves elements instead of copying during growth only when the move is `noexcept`; otherwise it copies, which can be far slower.

## Choosing a smart pointer

| Pointer | Use for | Not for |
| --- | --- | --- |
| `std::unique_ptr` | Exclusive ownership, the default | Sharing between owners |
| `std::shared_ptr` | Genuine shared ownership with independent lifetimes | Replacing `unique_ptr` out of habit |
| `std::weak_ptr` | Observing a `shared_ptr` without keeping it alive, breaking cycles | Dangling-pointer surveillance |
| raw pointer or reference | Non-owning observation | Owning anything |

```cpp
// GOOD: factory returns exclusive ownership
auto make_widget() { return std::make_unique<Widget>(); }

// GOOD: shared graph node, weak back-edge breaks the cycle
struct Node {
    std::vector<std::shared_ptr<Node>> children;
    std::weak_ptr<Node> parent;
};

// GOOD: observe without owning
void render(const Widget& widget);
```

`shared_ptr` has a control block and an atomic refcount. Use it only when two owners can outlive each other. Passing `shared_ptr` by value into a function that just reads the object is a common overuse; pass `const Widget&` instead.

## Value semantics by default

Return by value. The compiler elides the copy (guaranteed copy elision in C++17), so this is not slower than returning through a pointer.

```cpp
// GOOD: named return value optimization, no copy
std::vector<int> read_values() {
    std::vector<int> out;
    out.reserve(1024);
    // fill
    return out;
}
```

Do not return output parameters unless you have measured a reason:

```cpp
// BAD: two out-params, unclear ownership, callers can pass aliases
void read_values(std::vector<int>& out, bool& ok);

// GOOD: return a struct or use expected
std::expected<std::vector<int>, Error> read_values();
```

## Parameter passing

| Parameter kind | Pass by | Example |
| --- | --- | --- |
| Small scalar | value | `void set(int flags)`, `void f(std::string_view s)` |
| Read-only expensive object | `const&` | `void analyze(const std::string& data)` |
| Sink, will own | value, then move | `void set_name(std::string name)` |
| Optional output | return value | prefer return over out-param |
| Non-owning view | value | `std::span<const T>`, `std::string_view` |

```cpp
// Cheap types by value
void print(int x);
void log(std::string_view msg);

// Expensive read-only input by const&
void parse(const std::string& input);

// Sink: take by value, move into the member
class Sink {
public:
    explicit Sink(std::string name) : name_{std::move(name)} {}
private:
    std::string name_;
};
```

## Move, do not copy

`std::move` is a cast to an rvalue. It does not move anything by itself; it enables the move constructor or move assignment to run.

```cpp
// BAD: copies the vector because the source is an lvalue
v.push_back(other_vector);

// GOOD: moves the elements out of the source
v.push_back(std::move(other_vector));
```

Do not move from an object you still need, and do not `std::move` a return value. Returning a local by value already uses the move or elision, and `return std::move(local)` can defeat elision.

```cpp
// BAD: prevents copy elision and can pessimize
std::vector<int> f() {
    std::vector<int> out;
    return std::move(out);
}

// GOOD: let the compiler elide
std::vector<int> f() {
    std::vector<int> out;
    return out;
}
```
