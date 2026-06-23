# Safety

Never compromise on: resource leaks, deadlocks, classloader leaks, silent data corruption.

## Resource Leaks

```java
// ❌ BAD: Stream not closed if exception thrown
FileInputStream fis = new FileInputStream(path);
byte[] data = fis.readAllBytes();
fis.close();  // Never reached if readAllBytes() throws

// ✅ GOOD: Guaranteed cleanup
try (FileInputStream fis = new FileInputStream(path)) {
    byte[] data = fis.readAllBytes();
}
```

Always use try-with-resources for any `Closeable` (file streams, HTTP connections, database resources).

## Classloader Leaks

```java
// ❌ BAD: ThreadLocal never removed, holds classloader reference
ThreadLocal<RequestContext> context = new ThreadLocal<>();
context.set(new RequestContext());
// ... use it ...
// Classloader can't be GC'd after hot reload

// ✅ GOOD: Explicit cleanup
ThreadLocal<RequestContext> context = new ThreadLocal<>();
try {
    context.set(new RequestContext());
    // ... use it ...
} finally {
    context.remove();  // Releases classloader reference
}
```

Critical in application server / hot-reload environments (Quarkus dev mode, Spring Boot DevTools, app servers). Unremoved `ThreadLocal` values pin the classloader → OOM after repeated redeployments.

## Silent Data Corruption

```java
// ❌ BAD: Exception swallowed, order marked complete incorrectly
try {
    processPayment(order);
    order.setStatus(COMPLETE);
} catch (Exception e) { }  // Payment failed but order shows complete

// ✅ GOOD: Log and propagate
try {
    processPayment(order);
    order.setStatus(COMPLETE);
} catch (Exception e) {
    LOG.error("Payment failed for order {}", order.getId(), e);
    order.setStatus(FAILED);
    throw e;
}
```

Never swallow exceptions. At minimum log and rethrow as a runtime exception.

## Deadlocks

- Document lock acquisition order in comments
- Minimize critical section size
- Use `tryLock()` with timeout instead of unlimited `lock()`
- Prefer `ConcurrentHashMap`, atomic classes, and immutable data over explicit locks

## Violation Response

When a violation of these rules is detected in existing code, output a **CRITICAL SAFETY WARNING** block with:
- The specific risk (e.g. "potential deadlock between locks A and B")
- The technical context (code path, thread model)
- Actionable fix suggestions
