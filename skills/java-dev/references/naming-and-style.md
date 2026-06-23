# Naming Conventions & Code Style

## Naming Conventions

| Type | Convention | Examples |
|------|-----------|----------|
| Package | All lowercase, reverse domain | `com.example.project` |
| Class/Interface | PascalCase, noun/noun phrase | `UserService`, `HttpClient`, `Pageable` |
| Method | camelCase, verb prefix | `findById`, `isValid`, `hasPermission` |
| Constant | UPPER_SNAKE_CASE | `MAX_RETRY_COUNT` |
| Boolean methods | is/has/can prefix | `isActive()`, `hasPermission()` |
| DTO/VO | PascalCase with DTO/VO suffix | `UserDTO`, `CreateOrderRequest` |

## `final`

Mark parameters and variables `final` in new code unless mutability is required. Omit `this.` prefix unless required for disambiguation (e.g. constructor field assignments).

## Imports

Use simple class names with imports rather than fully qualified names, unless two classes share the same simple name in the same file.

## Never String-Literal Class/Package References

```java
// ❌ BAD: string literals — silently break on rename/move/repackage
Logger log = Logger.getLogger("com.example.OrderService");
Class<?> clazz = Class.forName("com.example.OrderService");

// ✅ GOOD: derived from .class — rename-safe, compile-time verified
Logger log = Logger.getLogger(OrderService.class.getName());
Class<?> clazz = OrderService.class;

// ✅ ACCEPTABLE: class does not exist on the classpath (optional plugin)
Class<?> clazz = Class.forName("com.thirdparty.OptionalExtension");
```

## Text Blocks for Multi-Line Strings

```java
// ❌ BAD: concatenation and escape sequences
String query = "SELECT id, name\n" +
               "FROM users\n" +
               "WHERE active = true";

// ✅ GOOD: text block
String query = """
        SELECT id, name
        FROM users
        WHERE active = true
        """;
```

## Reproducibility

Prefer deterministic behaviour. In non-performance-critical code, prefer sorted structures over hash-based ones to avoid ordering non-determinism.

In performance-critical runtime paths, efficiency takes precedence — but document the tradeoff. Security requirements (e.g. salted data structures) always take precedence.
