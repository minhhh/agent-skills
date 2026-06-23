# Testing

## Preferred Stack

- **JUnit 5** — the standard test runner
- **AssertJ** — fluent, readable assertions
- **MockServer / MockWebServer** — HTTP-level mocking of external services
- **Testcontainers** — real database in integration tests

## Principles

- Prefer real wiring over mocking. Reach for Mockito only when a dependency genuinely cannot be substituted with a real or in-memory implementation
- Strive for a fully automated integration test using a real database
- Add unit tests for classes with complex logic or data transformations. Skip unit tests when they only duplicate integration test coverage

```java
class UserServiceTest {
    @Test
    @DisplayName("findById returns user when found")
    void findById_whenUserExists_returnsUser() {
        // given
        when(userRepository.findById(1L)).thenReturn(Optional.of(expected));

        // when
        Optional<User> result = userService.findById(1L);

        // then
        assertThat(result).isPresent();
        assertThat(result.get().getName()).isEqualTo("test");
    }
}
```

## Framework-Specific

### Spring Boot

| Annotation | Use Case |
|------------|----------|
| `@SpringBootTest` | Full application context |
| `@WebMvcTest` | Web layer only (controllers) |
| `@DataJpaTest` | JPA repository slice |
| `@MockBean` / `@MockitoBean` | Replace a bean with mock |

**Test utilities:** `MockMvc` (web layer), `TestRestTemplate` / `RestClient` (REST integration), `@TestConfiguration` (overrides)

### Quarkus

| Annotation | Use Case |
|------------|----------|
| `@QuarkusTest` | Full CDI container |
| `@QuarkusIntegrationTest` | Black-box against built jar/native |
| `@QuarkusComponentTest` | Single bean, lightweight |

**Test utilities:** REST Assured (HTTP tests), `@TestHTTPResource` (test URLs), `@TestMock` / `@MockitoConfig` (mocking)
