# Advanced TypeScript Testing Patterns

Detailed patterns for the `typescript-testing` skill: API integration, database repositories, React components and hooks, snapshot testing, coverage, and timer utilities.

## Contents

- [API integration tests](#api-integration-tests)
- [Database repository tests](#database-repository-tests)
- [React component tests](#react-component-tests)
- [Hook tests](#hook-tests)
- [Snapshot testing](#snapshot-testing)
- [Coverage configuration](#coverage-configuration)
- [Timer and promise utilities](#timer-and-promise-utilities)

## API integration tests

Use `supertest` against the app instance, a real test database, and a truncate in `beforeEach`. Assert the status, the shape of the body, and that secrets never appear in a response.

```typescript
import request from "supertest";
import { app } from "../../src/app";
import { pool } from "../../src/config/database";

describe("User API", () => {
    beforeAll(async () => {
        await pool.query(`CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            email VARCHAR(255) UNIQUE NOT NULL,
            password VARCHAR(255) NOT NULL
        )`);
    });

    afterAll(async () => {
        await pool.query("DROP TABLE IF EXISTS users");
        await pool.end();
    });

    beforeEach(async () => {
        await pool.query("TRUNCATE TABLE users CASCADE");
    });

    it("creates a user and omits the password", async () => {
        const response = await request(app)
            .post("/api/users")
            .send({ name: "John Doe", email: "john@example.com", password: "password123" })
            .expect(201);

        expect(response.body).toMatchObject({ name: "John Doe", email: "john@example.com" });
        expect(response.body).toHaveProperty("id");
        expect(response.body).not.toHaveProperty("password");
    });

    it("rejects an invalid email", async () => {
        await request(app)
            .post("/api/users")
            .send({ name: "John Doe", email: "invalid-email", password: "password123" })
            .expect(400);
    });

    it("rejects a duplicate email", async () => {
        const body = { name: "John Doe", email: "john@example.com", password: "password123" };
        await request(app).post("/api/users").send(body);
        await request(app).post("/api/users").send(body).expect(409);
    });

    it("requires authentication for protected routes", async () => {
        await request(app).get("/api/users/me").expect(401);
    });

    it("allows access with a valid token", async () => {
        await request(app)
            .post("/api/users")
            .send({ name: "John Doe", email: "john@example.com", password: "password123" });

        const login = await request(app)
            .post("/api/auth/login")
            .send({ email: "john@example.com", password: "password123" })
            .expect(200);

        await request(app)
            .get("/api/users/me")
            .set("Authorization", `Bearer ${login.body.token}`)
            .expect(200);
    });
});
```

## Database repository tests

Test the repository against the real database engine. Assert the not-found path returns `null`, not an exception.

```typescript
import { describe, it, expect, beforeAll, afterAll, beforeEach } from "vitest";
import { Pool } from "pg";
import { UserRepository } from "../../src/repositories/user.repository";

describe("UserRepository", () => {
    let pool: Pool;
    let repository: UserRepository;

    beforeAll(async () => {
        pool = new Pool({ connectionString: process.env.TEST_DATABASE_URL });
        repository = new UserRepository(pool);
        await pool.query(`CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            email VARCHAR(255) UNIQUE NOT NULL
        )`);
    });

    afterAll(async () => {
        await pool.query("DROP TABLE IF EXISTS users");
        await pool.end();
    });

    beforeEach(async () => {
        await pool.query("TRUNCATE TABLE users CASCADE");
    });

    it("creates and finds a user by email", async () => {
        await repository.create({ name: "John Doe", email: "john@example.com" });
        const user = await repository.findByEmail("john@example.com");
        expect(user?.name).toBe("John Doe");
    });

    it("returns null when the user is missing", async () => {
        expect(await repository.findByEmail("none@example.com")).toBeNull();
    });
});
```

Use `testcontainers` to start a disposable database in CI instead of pointing at a shared instance.

## React component tests

Query by role or label. Avoid `data-testid` unless no semantic query exists.

```tsx
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { UserForm } from "./UserForm";

it("passes the entered values to onSubmit", async () => {
    const onSubmit = vi.fn();
    render(<UserForm onSubmit={onSubmit} />);

    await userEvent.type(screen.getByLabelText(/name/i), "John");
    await userEvent.type(screen.getByLabelText(/email/i), "john@example.com");
    await userEvent.click(screen.getByRole("button", { name: /submit/i }));

    expect(onSubmit).toHaveBeenCalledWith({ name: "John", email: "john@example.com" });
});
```

Cover the async states: loading, empty, and error. A component that renders data is not tested until those three branches are exercised.

## Hook tests

Use `renderHook` and wrap state updates in `act`.

```tsx
import { renderHook, act, waitFor } from "@testing-library/react";
import { useUsers } from "./useUsers";

it("loads users", async () => {
    const { result } = renderHook(() => useUsers());

    await waitFor(() => expect(result.current.loading).toBe(false));
    expect(result.current.users).toHaveLength(2);
});

it("updates the query", () => {
    const { result } = renderHook(() => useSearch());
    act(() => result.current.setQuery("abc"));
    expect(result.current.query).toBe("abc");
});
```

## Snapshot testing

Snapshots are useful for small, stable output. They are a liability for large objects: reviewers rubber-stamp updates and the assertion stops meaning anything.

```typescript
// GOOD: small, focused snapshot
it("renders the empty state", () => {
    const { container } = render(<EmptyState />);
    expect(container.firstChild).toMatchSnapshot();
});

// BAD: a 400-line object snapshot nobody reads
expect(entireApiResponse).toMatchSnapshot();
```

Prefer explicit field assertions for API payloads. When using snapshots, keep them small and run in CI with `--ci` so new snapshots fail instead of being written silently.

## Coverage configuration

Set thresholds where CI enforces them, and exclude generated and config files.

```typescript
// vitest.config.ts
export default defineConfig({
    test: {
        coverage: {
            provider: "v8",
            reporter: ["text", "json-summary", "html"],
            include: ["src/**/*.ts", "src/**/*.tsx"],
            exclude: ["src/**/*.d.ts", "src/**/*.config.ts", "src/main.tsx"],
            thresholds: { branches: 80, functions: 80, lines: 80, statements: 80 },
        },
    },
});
```

```bash
vitest run --coverage
jest --ci --coverage
```

## Timer and promise utilities

Fake timers make time-dependent code deterministic. Use the async advance so pending promises settle between ticks.

```typescript
import { vi } from "vitest";

it("retries after a delay", async () => {
    vi.useFakeTimers();
    const fn = vi.fn().mockRejectedValueOnce(new Error("fail")).mockResolvedValue("ok");

    const promise = retry(fn, { attempts: 2, delayMs: 100 });
    await vi.advanceTimersByTimeAsync(100);
    await expect(promise).resolves.toBe("ok");

    vi.useRealTimers();
});
```

Assert rejection shapes rather than just "it threw":

```typescript
await expect(loadUser("999")).rejects.toThrowError(
    expect.objectContaining({ name: "NotFoundError" }),
);
```

Flush pending microtasks when a test schedules work without a returned promise:

```typescript
await vi.waitFor(() => expect(onDone).toHaveBeenCalled());
```
