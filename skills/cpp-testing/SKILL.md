---
name: cpp-testing
description: Use when writing, reviewing, or fixing C++ tests, setting up GoogleTest, GoogleMock, or CTest with CMake, diagnosing failing or flaky tests, or enabling sanitizers, fuzzing, benchmarks, and coverage.
---

# C++ Testing

GoogleTest and GoogleMock with CMake and CTest. C++ has no runtime to catch memory and lifetime mistakes, so tests plus sanitizers are how you get confidence.

## When to Activate

- Writing new C++ tests or fixing existing ones
- Setting up a test target in CMake with CTest discovery
- Diagnosing a failing or flaky test
- Adding sanitizer, fuzzing, benchmark, or coverage workflows
- Reviewing a change that needs a regression test

## TDD Workflow

Follow RED, GREEN, REFACTOR.

1. RED: write a failing test that describes the new behavior. Run it and watch it fail.
2. GREEN: write the smallest change that makes it pass.
3. REFACTOR: clean up while the test stays green.

```cpp
// tests/add_test.cpp
#include <gtest/gtest.h>

int Add(int a, int b);  // production code

TEST(AddTest, AddsTwoNumbers) {  // RED: fails to link or fails the assert
    EXPECT_EQ(Add(2, 3), 5);
}

// src/add.cpp
int Add(int a, int b) { return a + b; }  // GREEN
```

## Basic tests

`EXPECT_*` records a failure and continues. `ASSERT_*` records a failure and returns from the test. Use `ASSERT` when the rest of the test cannot run after the failure.

```cpp
TEST(CalculatorTest, AddsTwoNumbers) {
    EXPECT_EQ(Add(2, 3), 5);
}

TEST(ParserTest, RejectsEmptyInput) {
    auto result = parse("");
    ASSERT_FALSE(result.has_value());  // stop if it unexpectedly parsed
    EXPECT_EQ(result.error(), ParseError::Empty);
}
```

Prefer `EXPECT_THAT` with matchers for readable assertions.

```cpp
#include <gmock/gmock.h>

using ::testing::ElementsAre;
using ::testing::HasSubstr;

TEST(UserTest, RolesAreOrdered) {
    auto roles = user.roles();
    EXPECT_THAT(roles, ElementsAre("admin", "editor"));
}

TEST(LogTest, ContainsUser) {
    EXPECT_THAT(last_log_message(), HasSubstr("alice"));
}
```

## Fixtures

Use a fixture when several tests share setup. `SetUp` runs before each test, `TearDown` after.

```cpp
class UserStoreTest : public ::testing::Test {
protected:
    void SetUp() override {
        store_ = std::make_unique<UserStore>(":memory:");
        store_->seed({{"alice"}, {"bob"}});
    }

    // If a test can throw, name the cleanup method TearDown, not a destructor.
    std::unique_ptr<UserStore> store_;
};

TEST_F(UserStoreTest, FindsExistingUser) {
    auto user = store_->find("alice");
    ASSERT_TRUE(user.has_value());
    EXPECT_EQ(user->name, "alice");
}

TEST_F(UserStoreTest, MissesUnknownUser) {
    EXPECT_FALSE(store_->find("carol").has_value());
}
```

Do not put logic that can fail in the fixture destructor. A throwing destructor terminates the test binary. Use `TearDown` and `ASSERT` there, or use `TearDown` plus `EXPECT`.

## Parameterized tests

Table-driven tests with a name generator.

```cpp
struct ParseCase {
    std::string input;
    std::optional<int> expected;
};

class ParseTest : public ::testing::TestWithParam<ParseCase> {};

TEST_P(ParseTest, ParsesOrRejects) {
    const auto& param = GetParam();
    auto result = parse_int(param.input);
    if (param.expected) {
        ASSERT_TRUE(result.has_value()) << param.input;
        EXPECT_EQ(*result, *param.expected);
    } else {
        EXPECT_FALSE(result.has_value()) << param.input;
    }
}

INSTANTIATE_TEST_SUITE_P(
    ParseCases,
    ParseTest,
    ::testing::Values(
        ParseCase{"42", 42},
        ParseCase{"0", 0},
        ParseCase{"", std::nullopt},
        ParseCase{"12x", std::nullopt}
    ),
    [](const ::testing::TestParamInfo<ParseCase>& info) {
        return info.param.input.empty() ? "empty" : "in_" + info.param.input;
    });
```

The name generator must return a string without invalid characters. Failure output shows the name, so make it meaningful.

## Mocks and fakes

Prefer a fake with real behavior over a mock that records calls. Mocks drift from production. Use a mock when you need to assert an interaction or simulate an external boundary.

```cpp
class UserRepository {
public:
    virtual ~UserRepository() = default;
    virtual std::optional<User> find(const std::string& name) = 0;
};

class MockUserRepository : public UserRepository {
public:
    MOCK_METHOD(std::optional<User>, find, (const std::string& name), (override));
};

TEST(UserServiceTest, ReturnsNotFound) {
    MockUserRepository repo;
    EXPECT_CALL(repo, find("alice")).WillOnce(::testing::Return(std::nullopt));

    UserService service{repo};
    EXPECT_FALSE(service.profile("alice").has_value());
}
```

For stateful behavior, a hand-written fake is clearer:

```cpp
class InMemoryUserRepository : public UserRepository {
public:
    std::optional<User> find(const std::string& name) override {
        auto it = users_.find(name);
        if (it == users_.end()) return std::nullopt;
        return it->second;
    }
    void save(User user) { users_[user.name] = std::move(user); }
private:
    std::map<std::string, User> users_;
};
```

## Death tests

Death tests assert that code terminates in a controlled way. Use them sparingly, for `assert` and invariant checks.

```cpp
TEST(ConfigTest, DiesOnMissingRequired) {
    EXPECT_DEATH(Config::load("missing.toml"), "required key");
}
```

Death tests fork the process and are slow and platform-sensitive. Keep them few, and do not use them to test ordinary error handling.

## CMake and CTest

Register tests so `ctest` discovers them. `gtest_discover_tests` runs each test case as its own CTest test and names it after the case.

```cmake
include(GoogleTest)

add_executable(unit_tests
    tests/add_test.cpp
    tests/parser_test.cpp
)
target_link_libraries(unit_tests
    PRIVATE
        mylib
        GTest::gtest_main
        GTest::gmock
)

gtest_discover_tests(unit_tests)
```

Run them:

```bash
cmake --preset debug
cmake --build --preset debug
ctest --preset debug --output-on-failure
```

`--output-on-failure` shows the failing test's output. Without it, CTest reports only pass or fail. Use `-R` to filter and `--repeat until-fail:10` to shake out flaky tests.

```bash
ctest --test-dir build/debug -R ParseTest --output-on-failure
ctest --test-dir build/debug --repeat until-fail:10
```

## Sanitizers

Run the test suite under sanitizers in CI. This is the highest-value C++ testing practice after the tests themselves.

```cmake
add_executable(unit_tests_asan tests/...)
target_compile_options(unit_tests_asan PRIVATE
    -fsanitize=address,undefined -fno-omit-frame-pointer -g)
target_link_options(unit_tests_asan PRIVATE
    -fsanitize=address,undefined)
```

| Sanitizer | Finds | Cannot run with |
| --- | --- | --- |
| `address` | use-after-free, overflow, leaks | `thread` |
| `undefined` | signed overflow, misaligned access, null deref | nothing in particular |
| `thread` | data races | `address` |
| `memory` | uninitialized reads, Clang only | `address` |

Set `ASAN_OPTIONS=detect_leaks=1` and `UBSAN_OPTIONS=print_stacktrace=1:halt_on_error=1` in CI so a finding fails the build.

## Fuzzing

libFuzzer finds input-dependent crashes. Add a fuzz target for any parser or decoder.

```cpp
#include <cstddef>
#include <cstdint>
#include <string_view>

extern "C" int LLVMFuzzerTestOneInput(const std::uint8_t* data, std::size_t size) {
    parse_config(std::string_view{
        reinterpret_cast<const char*>(data), size});
    return 0;
}
```

```cmake
target_compile_options(fuzz_config PRIVATE
    -fsanitize=fuzzer,address,undefined -g)
target_link_options(fuzz_config PRIVATE
    -fsanitize=fuzzer,address,undefined)
```

```bash
./build/debug/fuzz_config -max_total_time=60 -artifact_prefix=./crashes/
```

Fuzz with AddressSanitizer and UndefinedBehaviorSanitizer enabled so crashes point at the exact fault.

## Benchmarks

Use Google Benchmark for microbenchmarks. Measure, do not guess.

```cpp
#include <benchmark/benchmark.h>

static void BM_StringConcat(benchmark::State& state) {
    for (auto _ : state) {
        std::string result = join(parts, "");
        benchmark::DoNotOptimize(result);
    }
}
BENCHMARK(BM_StringConcat);

BENCHMARK_MAIN();
```

`benchmark::DoNotOptimize` stops the compiler from deleting the work. Compare runs with a stable machine and report the delta, not a single number.

## Coverage

```bash
cmake --build build/coverage
gcovr --root . --print-summary build/coverage
```

| Code type | Target |
| --- | --- |
| Critical logic | 100% |
| Public API | 90%+ |
| General code | 80%+ |
| Generated code | exclude |

Coverage shows what is not tested. It does not show whether the tests assert the right thing.

## Best practices

Do:

- Test behavior through the public interface.
- Write the failing test first when fixing a bug.
- Use `ASSERT` before a dereference, `EXPECT` for independent checks.
- Name tests so the failure output is self-explanatory.
- Run sanitizers on every test run in CI.
- Keep tests independent; no shared mutable global state.

Do not:

- Sleep in tests. Use a condition variable or a fake clock.
- Mock everything. Prefer fakes and real in-memory implementations.
- Test private members directly.
- Ignore a flaky test. Fix it or delete it with a note.
- Exclude tests from sanitizer builds.

## CI

```yaml
test:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - name: Configure
      run: cmake --preset debug
    - name: Build
      run: cmake --build --preset debug
    - name: Test
      run: ctest --preset debug --output-on-failure
    - name: Sanitizer build
      run: |
        cmake --preset asan
        cmake --build --preset asan
        ctest --preset asan --output-on-failure
```

Run the plain build first for fast feedback, then the sanitizer build.

## Common pitfalls

| Mistake | Why it hurts | Fix |
| --- | --- | --- |
| `EXPECT` before a dereference | The test crashes instead of reporting | `ASSERT` on the optional or pointer first |
| Logic in a fixture destructor | A throw there terminates the binary | Move it to `TearDown` |
| Tests share a global | Order-dependent failures | Fixture per test, no shared mutable state |
| Testing private methods | Couples tests to implementation | Test through the public API |
| No sanitizer run | Lifetime bugs pass the tests | Add an ASan or UBSan CI job |
| Flaky test left in place | Erodes trust in the suite | Fix the race or the timing assumption |
| Fuzzing without sanitizers | Crashes with no diagnosis | `-fsanitize=fuzzer,address,undefined` |
