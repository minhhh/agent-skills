# Exception Handling & Logging

## Exception Handling

```java
// ✅ GOOD: Catch specific exception, add context
try {
    user = userRepository.findById(id);
} catch (DataAccessException e) {
    throw new ServiceException("Failed to find user: " + id, e);
}

// ❌ BAD: Catch too broad
catch (Exception e) { e.printStackTrace(); }
```

## Null Safety

```java
// ✅ Use Optional for nullable returns
public Optional<User> findById(Long id) {
    return userRepository.findById(id);
}

// ✅ Safe access chain
String name = Optional.ofNullable(user)
    .map(User::getName)
    .orElse("Unknown");

// ✅ Parameter validation
public void updateUser(User user) {
    Objects.requireNonNull(user, "user must not be null");
}
```

## Logging

```java
// ✅ GOOD: Parameterized
log.debug("Finding user by id: {}", userId);
log.error("Failed to process order {}", orderId, exception);

// ❌ BAD: String concatenation
log.debug("Finding user by id: " + userId);
```

- Use parameterized logging (`{}` placeholders) — avoids string construction when level is disabled
- Always include the exception object as the last parameter in `log.error()`
- Use appropriate levels: `ERROR` (failures), `WARN` (unexpected but handled), `INFO` (lifecycle), `DEBUG` (details)
