# Performance

## Hot Path Optimization

```java
// ❌ BAD: Stream overhead in per-request path
public List<String> getActive() {
    return items.stream()
        .filter(Item::isActive)
        .map(Item::getName)
        .collect(Collectors.toList());
}

// ✅ GOOD: Simple loop for hot path
public List<String> getActive() {
    List<String> result = new ArrayList<>(items.size());
    for (Item item : items) {
        if (item.isActive()) {
            result.add(item.getName());
        }
    }
    return result;
}
```

## Avoid Unnecessary Boxing

```java
// ❌ BAD: Boxing creates GC pressure
List<Integer> counts = getCounts();
int sum = 0;
for (Integer count : counts) { sum += count; }

// ✅ GOOD: Primitives when possible
int[] counts = getCounts();
int sum = 0;
for (int count : counts) { sum += count; }
```

**What counts as performance-critical:** tight loops, per-request processing, high-frequency paths. Config parsing, startup code, and build-time logic are generally not critical.

For hot paths, measure before optimizing — don't pre-optimize cold code.

## N+1 Query Prevention

```java
// ❌ BAD: N+1 — one query per iteration
records.forEach(record -> {
    long count = deviceRepo.countByDeviceId(record.getDeviceId());
    record.setDeviceCount(count);
});

// ✅ GOOD: Batch query + Map lookup
List<String> deviceIds = records.stream()
    .map(Record::getDeviceId).distinct().collect(Collectors.toList());
Map<String, Long> countMap = deviceRepo.countByDeviceIdIn(deviceIds).stream()
    .collect(Collectors.toMap(CountDTO::getDeviceId, CountDTO::getCount));
records.forEach(r -> r.setDeviceCount(countMap.getOrDefault(r.getDeviceId(), 0L)));
```

| Scenario | Loop (❌) | Batch (✅) |
|----------|-----------|------------|
| count | `repo.countByXxx(id)` | `repo.countByXxxIn(ids)` → `Map<id, count>` |
| findById | `repo.findById(id)` | `repo.findByIdIn(ids)` → `Map<id, entity>` |
| exists | `repo.existsByXxx(id)` | `repo.findXxxIn(ids)` → `Set<id>` |

## Batch Query Limits

```java
// ❌ BAD: IN clause with thousands of parameters
List<User> users = userRepository.findByIdIn(allIds);  // SQL too large, unstable plan

// ✅ GOOD: Batch query helper
public static <T, R> List<R> batchQuery(List<T> params, int batchSize,
                                         Function<List<T>, List<R>> queryFn) {
    List<R> result = new ArrayList<>();
    for (int i = 0; i < params.size(); i += batchSize) {
        List<T> batch = params.subList(i, Math.min(i + batchSize, params.size()));
        result.addAll(queryFn.apply(batch));
    }
    return result;
}

List<User> users = batchQuery(allIds, 500, ids -> userRepository.findByIdIn(ids));
```

**Limit:** Maximum 500 parameters per IN clause.

## Collection Initial Capacity

```java
// ❌ BAD: Default capacity, causes repeated resizing
List<String> names = new ArrayList<>();
for (Item item : items) { names.add(item.getName()); }

// ✅ GOOD: Pre-sized
List<String> names = new ArrayList<>(items.size());
```

## String Building in Loops

```java
// ❌ BAD: String concatenation in loop creates many intermediates
String s = "";
for (int i = 0; i < 1000; i++) { s += i; }

// ✅ GOOD: StringBuilder
StringBuilder sb = new StringBuilder(4000);
for (int i = 0; i < 1000; i++) { sb.append(i); }
String s = sb.toString();
```
