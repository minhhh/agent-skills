# Concurrency

## Thread Model Documentation

Most of our state is confined to a single thread. Always document the thread model.

```java
// ❌ BAD: No indication of thread model
public class EventProcessor {
    private List<Event> buffer = new ArrayList<>();  // Is this shared?

    public void add(Event e) {
        buffer.add(e);
    }
}

// ✅ GOOD: Explicit thread model
public class EventProcessor {
    // NOT thread-safe — designed for single-threaded use only
    private List<Event> buffer = new ArrayList<>();

    public void add(Event e) {
        buffer.add(e);
    }
}
```

## Read-Modify-Write

```java
// ❌ BAD: read-modify-write loses updates under concurrency
PointsAccount account = accountRepo.findById(id);
account.setBalance(account.getBalance() + points);
accountRepo.save(account);  // Two concurrent reads see same balance, second write wins

// ✅ GOOD: Atomic update in the DB
@Modifying
@Query("UPDATE PointsAccount SET balance = balance + :points WHERE id = :id")
int addBalance(@Param("id") Long id, @Param("points") int points);

// ✅ GOOD: Optimistic locking with JPA @Version
@Version
private Long version;  // Hibernate retries on conflict
```

## Check-Then-Act

```java
// ❌ BAD: Race between check and insert
if (!rewardRepo.existsByTenantIdAndPeriod(tenantId, period)) {
    rewardRepo.save(new RankingReward(...));  // Two threads both pass check
}

// ✅ GOOD: Unique constraint as last line of defence
// DDL: UNIQUE INDEX uk_tenant_period (tenant_id, period)
try {
    rewardRepo.save(new RankingReward(...));
} catch (DataIntegrityViolationException e) {
    log.warn("Duplicate prevented by unique index: {}", tenantId);
}
```

## Executor Service

```java
// ❌ BAD: Direct thread creation
new Thread(() -> doWork()).start();

// ✅ GOOD: ExecutorService
ExecutorService executor = Executors.newFixedThreadPool(10);
Future<Result> future = executor.submit(() -> doWork());

// ✅ GOOD: CompletableFuture
CompletableFuture<User> future = CompletableFuture
    .supplyAsync(() -> findUser(id))
    .thenApply(user -> enrichUser(user));
```

### Thread Pool Configuration

```java
private static final int CPU_COUNT = Runtime.getRuntime().availableProcessors();

// I/O-heavy tasks
private final ThreadPoolExecutor ioExecutor = new ThreadPoolExecutor(
    CPU_COUNT * 2,      // core threads
    CPU_COUNT * 4,      // max threads (bounded!)
    60L, TimeUnit.SECONDS,
    new LinkedBlockingQueue<>(10000),  // bounded queue
    new ThreadPoolExecutor.AbortPolicy()
);

// CPU-heavy tasks
private final ThreadPoolExecutor cpuExecutor = new ThreadPoolExecutor(
    CPU_COUNT,
    CPU_COUNT + 1,
    60L, TimeUnit.SECONDS,
    new LinkedBlockingQueue<>(1000),
    new ThreadPoolExecutor.CallerRunsPolicy()  // backpressure
);
```

Avoid `Executors.newFixedThreadPool()` (unbounded queue → OOM) and `Executors.newCachedThreadPool()` (unbounded thread count).

## Framework-Specific

- **Quarkus/Vert.x:** Never block the I/O thread. Use `@Blocking` for JDBC calls.
- **Spring Boot:** Use `@Async` + `ThreadPoolTaskExecutor` for background tasks.
