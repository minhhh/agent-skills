# APIs, DTOs & Input Validation

## DTO/VO Conventions

| Rule | When |
|------|------|
| Use Lombok `@Data` | Plain DTOs, request/response objects |
| Use Lombok `@Value` | Immutable DTOs |
| Use Lombok `@Builder` | DTOs with many fields |
| Avoid `@Data` on JPA `@Entity` | `equals`/`hashCode` breaks with Hibernate proxies |

```java
// ❌ BAD: Hand-written getter/setter boilerplate
public class UserDTO {
    private Long id;
    private String name;
    public Long getId() { return id; }
    public void setId(Long id) { this.id = id; }
    // ...
}

// ✅ GOOD: Lombok
@Data
public class UserDTO {
    private Long id;
    private String name;
}
```

## Input Validation

```java
// ❌ BAD: No validation
@PostMapping("/ship")
public Result ship(@RequestBody ShippingRequest request) { ... }

// ✅ GOOD: Validated
@PostMapping("/ship")
public Result ship(@RequestBody @Valid ShippingRequest request) { ... }

public record ShippingRequest(
    @NotNull Long orderId,
    @NotBlank @Size(max = 500) String shippingInfo
) {}
```

| Field Type | Required Annotation | Reason |
|-----------|-------------------|--------|
| quantity | `@NotNull @Min(1)` | Prevent 0/negative |
| amount/price | `@NotNull @Positive` | `@DecimalMin("0.01")` |
| pagination size | `@Min(1) @Max(100)` | Prevent DB overload |
| pagination page | `@Min(0)` | 0-based indexing |
| percentage | `@Min(0) @Max(100)` | Bounded range |

## Clean APIs

Prioritise well-designed interfaces over expedient ones. If an existing abstraction is the wrong shape, improve it rather than working around it — workarounds accumulate.

When designing an API or SPI, ask: can someone understand what this does without reading its implementation? Can the implementation change without breaking consumers? If not, the boundary needs work.

## Code Consolidation

Check before writing new helpers or utilities — prefer extension or composition over duplication.

- **Within a class or package** — extract shared logic into a method or utility
- **Across modules in a multi-module project** — find the right owner module
- **Across repos in a multi-repo platform** — the most expensive kind; it diverges silently
