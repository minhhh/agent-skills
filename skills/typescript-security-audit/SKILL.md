---
name: typescript-security-audit
description: Use when the user requests a security review or audit in a TypeScript or Node.js project, or when reviewing code that handles authentication, authorization, payment, PII, untrusted input, or user-supplied URLs.
---

# Security Audit for TypeScript and Node.js Applications

You are a security reviewer for TypeScript and Node.js applications. Find vulnerabilities before they reach production, with focus on the OWASP Top 10 adapted for the Node.js ecosystem.

## Prerequisites

This skill builds on [`security-audit-principles`] and [`typescript-dev`].

Apply all rules from:

- **`security-audit-principles`**: OWASP Top 10 vulnerability categories, severity assessment, audit workflow and reporting format
- **`typescript-dev`**: type safety, strict mode, type-driven input validation, and testing practices

Then apply the TypeScript and Node.js security patterns below.

## Workflow

Follow the `security-audit-principles` workflow (Steps 1 to 4). TypeScript-specific Step 3 example:

```text
🔴 CRITICAL — userService.ts:34
SQL Injection: the query is built with string concatenation, so an attacker can
execute arbitrary SQL, dump the database, or modify records.

Suggested fix:
  await db.query('SELECT * FROM users WHERE id = $1', [userId])
```

Step 2 uses the TypeScript and Node.js Security Checklist below in place of the generic categories.

**If no CRITICAL findings:**

> "No critical vulnerabilities found. [N warnings / notes listed above.] Consider addressing warnings for defense in depth."

## Severity Assignment Decision Flow

```mermaid
flowchart TD
    Finding_detected((Finding detected))
    Exploitable_remotely_{Exploitable remotely?}
    Breach_or_RCE_{Leads to data breach or RCE?}
    Needs_auth_{Requires authentication?}
    CRITICAL[CRITICAL]
    WARNING[WARNING]
    NOTE[NOTE]
    Finding_detected --> Exploitable_remotely_
    Exploitable_remotely_ -->|"no (local only)"| NOTE
    Exploitable_remotely_ -->|yes| Breach_or_RCE_
    Breach_or_RCE_ -->|yes| CRITICAL
    Breach_or_RCE_ -->|no| Needs_auth_
    Needs_auth_ -->|"yes (auth required)"| WARNING
    Needs_auth_ -->|"no (unauthenticated)"| CRITICAL
```

## Security Checklist

### 🔴 A01 Injection

**SQL injection.** Never build a query with string concatenation or a template literal.

```typescript
// BAD: SQL injection
const user = await db.query(`SELECT * FROM users WHERE name = '${req.params.name}'`);

// GOOD: parameterized query (node-postgres)
const user = await db.query("SELECT * FROM users WHERE name = $1", [req.params.name]);

// GOOD: Prisma parameterizes automatically
const user = await prisma.user.findFirst({ where: { name: req.params.name } });

// GOOD: TypeORM query builder with parameters
const user = await repo.createQueryBuilder("user")
    .where("user.name = :name", { name: req.params.name })
    .getOne();
```

**Command injection.** Never pass user input to a shell.

```typescript
// BAD: attacker sends "file.jpg; rm -rf /"
exec(`convert ${req.body.filename} output.pdf`);

// GOOD: validate, and use execFile without a shell
import { execFile } from "node:child_process";
if (!/^[\w.-]+$/.test(req.body.filename)) throw new Error("Invalid filename");
execFile("convert", [req.body.filename, "output.pdf"]);
```

**NoSQL injection.** Object queries accept operators, so `{ "$gt": "" }` can bypass a password check. Validate the type and use a schema.

```typescript
// BAD: attacker sends { "$gt": "" } as the password
const user = await User.findOne({ email, password: req.body.password });

// GOOD: parse into a typed shape first
import { z } from "zod";
const LoginSchema = z.object({ email: z.string().email(), password: z.string().min(8) });
const { email, password } = LoginSchema.parse(req.body);
const user = await User.findOne({ email });
```

**Template and XSS injection.** Escape before rendering user content as HTML.

```typescript
// BAD: XSS
const html = `<h1>Hello ${req.body.name}</h1>`;

// GOOD: escape the value
import escapeHtml from "escape-html";
const html = `<h1>Hello ${escapeHtml(req.body.name)}</h1>`;

// Browser: prefer textContent, or sanitize with DOMPurify
element.textContent = userInput;
element.innerHTML = DOMPurify.sanitize(userInput);
```

**`eval` and `new Function`.** Never execute an untrusted string. If a dynamic execution path exists, treat it as CRITICAL.

### 🔴 A02 Broken authentication

**Hardcoded secrets.** No API keys, tokens, or passwords in source. Validate the environment at startup so a missing secret fails fast.

```typescript
// BAD
const jwtSecret = "my-super-secret-key";

// GOOD
import { z } from "zod";
const EnvSchema = z.object({
    JWT_SECRET: z.string().min(32),
    DATABASE_URL: z.string().url(),
    SESSION_SECRET: z.string().min(32),
});
const env = EnvSchema.parse(process.env);
```

**Weak JWT configuration.** Pin the algorithm and set an expiry.

```typescript
// BAD: no expiry, weak secret, no algorithm pinning
const token = jwt.sign({ userId }, "secret");

// GOOD
const token = jwt.sign({ userId }, env.JWT_SECRET, { algorithm: "HS256", expiresIn: "15m" });
const payload = jwt.verify(token, env.JWT_SECRET, { algorithms: ["HS256"] });
```

**Weak password hashing.** Use bcrypt with cost 12 or higher, or argon2.

```typescript
// BAD: MD5 or SHA1
const hash = createHash("md5").update(password).digest("hex");

// GOOD: argon2 (preferred) or bcrypt cost >= 12
const hash = await argon2.hash(password);
const valid = await argon2.verify(hash, password);
```

**Session tokens in `localStorage`.** XSS can read them. Use an httpOnly, Secure, SameSite cookie.

```typescript
res.cookie("session", token, {
    httpOnly: true,
    secure: true,
    sameSite: "strict",
    maxAge: 15 * 60 * 1000,
});
```

**Missing rate limiting.** Auth, registration, and expensive endpoints need a limiter.

```typescript
import rateLimit from "express-rate-limit";
const authLimiter = rateLimit({ windowMs: 15 * 60 * 1000, max: 10 });
app.post("/login", authLimiter, loginHandler);
```

### 🔴 A03 Broken access control

**IDOR.** Authenticated is not authorized. Check ownership before every read, update, or delete.

```typescript
// BAD: any authenticated user can delete any order
app.delete("/orders/:id", authenticate, async (req, res) => {
    await Order.findByIdAndDelete(req.params.id);
    res.sendStatus(204);
});

// GOOD: verify ownership
app.delete("/orders/:id", authenticate, async (req, res) => {
    const order = await Order.findById(req.params.id);
    if (!order || order.userId.toString() !== req.user.id) {
        return res.status(403).json({ error: "Forbidden" });
    }
    await order.deleteOne();
    res.sendStatus(204);
});
```

**Mass assignment.** Never spread `req.body` onto a model. Use an explicit allowlist or a validated DTO that strips unknown fields.

```typescript
// BAD: attacker sends { "isAdmin": true }
await User.findByIdAndUpdate(req.user.id, req.body);

// GOOD: allowlist, or a schema that strips unknown keys
const { name, email } = req.body;
await User.findByIdAndUpdate(req.user.id, { name, email });
```

### 🔴 A04 Cryptographic failures

**Sensitive data in logs.** Log identifiers, never full user objects, tokens, or request bodies.

```typescript
// BAD: leaks password hashes and tokens
console.log("User authenticated:", user);

// GOOD
logger.info("User authenticated", { userId: user.id });
```

**Insecure randomness.** `Math.random()` is predictable. Use `crypto.randomBytes` for tokens and ids.

```typescript
// BAD
const resetToken = Math.random().toString(36).slice(2);

// GOOD
import { randomBytes } from "node:crypto";
const resetToken = randomBytes(32).toString("hex");
```

**Unencrypted PII at rest.** Encrypt sensitive fields before storage and manage the key in the environment or a KMS.

### 🔴 A05 Security misconfiguration

**Stack traces to the client.** Log server-side, return a generic message.

```typescript
// BAD: leaks internals
res.status(500).json({ error: err.message, stack: err.stack });

// GOOD
logger.error({ err, path: req.path }, "Unhandled error");
res.status(500).json({ error: "Internal server error" });
```

**Missing security headers.** Apply `helmet()` to every Express app.

**Permissive CORS.** `origin: "*"` with credentials lets any site make credentialed requests. Use an explicit allowlist.

```typescript
app.use(cors({
    origin: ["https://app.example.com", "https://admin.example.com"],
    methods: ["GET", "POST", "PUT", "DELETE"],
    credentials: true,
}));
```

**Missing input validation.** Validate every external boundary with a schema: HTTP, WebSocket, message queues, file uploads.

### 🟡 A06 Vulnerable and outdated components

```bash
npm audit                 # known CVEs
npm audit fix             # patch and minor, inspect the diff first
npm outdated              # available updates
```

Flag every HIGH or CRITICAL advisory, packages with no publish activity in over two years, and transitive advisories. `npm audit fix --force` can introduce breaking major versions; review the diff before committing. Run `npm audit` in CI and fail on HIGH.

### 🟡 A07 Prototype pollution

Merging unvalidated objects lets an attacker set `__proto__` and pollute `Object.prototype`.

```typescript
// BAD: attacker sends { "__proto__": { "isAdmin": true } }
import merge from "lodash/merge";
const config = merge({}, req.body);

// GOOD: validate with a schema that strips unknown keys
const SafeConfigSchema = z.object({ theme: z.string(), lang: z.string() });
const config = SafeConfigSchema.parse(req.body);

// GOOD: use Object.create(null) for lookup tables
const lookup: Record<string, string> = Object.create(null);
```

Also flag `JSON.parse` of user input followed by direct use without schema validation.

### 🟡 A08 Server-side request forgery

Fetching a user-supplied URL lets an attacker reach internal services, including the cloud metadata endpoint at `169.254.169.254`.

```typescript
// BAD: no allowlist
app.get("/proxy", async (req, res) => {
    const response = await fetch(req.query.url as string);
    res.send(await response.text());
});

// GOOD: explicit allowlist, and reject non-https
const ALLOWED = ["https://api.example.com/", "https://cdn.example.com/"];
app.get("/proxy", async (req, res) => {
    const url = String(req.query.url);
    if (!ALLOWED.some((prefix) => url.startsWith(prefix))) {
        return res.status(400).json({ error: "URL not allowed" });
    }
    const response = await fetch(url);
    res.send(await response.text());
});
```

**Open redirect.** Only allow relative paths for post-login redirects.

```typescript
const redirect = String(req.query.redirect ?? "");
if (redirect && !redirect.startsWith("/")) {
    return res.redirect("/dashboard");
}
res.redirect(redirect || "/dashboard");
```

### 🔵 Defense in depth

- Parse every external boundary with a schema (`zod`, `joi`, `valibot`).
- Rate-limit authentication, registration, and expensive endpoints.
- Use a least-privilege database user per service; never root or superuser at runtime.
- Log authentication attempts (success and failure), authorization failures, and PII or payment access, with userId, timestamp, action, and resource.
- Validate environment variables at startup and fail fast when a required secret is missing.
- Rotate secrets and keep them out of the repository and out of logs.

## Node.js security tooling

| Tool | Purpose | Notes |
| --- | --- | --- |
| `helmet` | Secure HTTP headers in one call | Use defaults, then tune CSP |
| `express-rate-limit` | Rate limiting | Apply to auth and expensive routes |
| `zod` / `joi` / `valibot` | Runtime input validation | `zod` gives type inference |
| `bcrypt` / `argon2` | Password hashing | argon2 preferred; bcrypt cost >= 12 |
| `jsonwebtoken` | JWT sign and verify | Pin `algorithms`, always set `expiresIn` |
| `node:crypto` | Secure random and encryption | `randomBytes` for tokens, never `Math.random()` |
| `escape-html` / `DOMPurify` | Output encoding | Escape before rendering user content |
| `dotenv` plus a schema | Env validation | Parse `process.env` once at startup |

```typescript
// Fail fast rather than run misconfigured
const EnvSchema = z.object({
    NODE_ENV: z.enum(["development", "test", "production"]),
    JWT_SECRET: z.string().min(32),
    DATABASE_URL: z.string().url(),
    SESSION_SECRET: z.string().min(32),
});
const env = EnvSchema.parse(process.env);
```

## Common Pitfalls

| Mistake | Impact | Fix |
| --- | --- | --- |
| SQL built with string concatenation | Full database compromise | Parameterized queries or an ORM query builder |
| Spreading `req.body` onto a model | Mass assignment, privilege escalation | Explicit allowlist or a validated DTO |
| Secrets hardcoded in source | Credential exposure | Environment variables validated at startup |
| `Math.random()` for tokens | Predictable tokens, account takeover | `crypto.randomBytes(32)` |
| Logging full user objects | PII and secrets in log aggregators | Log safe identifiers only |
| Session token in `localStorage` | XSS steals the token | httpOnly, Secure, SameSite cookie |
| `cors({ origin: "*" })` in production | Credentialed cross-origin requests | Explicit origin allowlist |
| No `helmet` | Clickjacking, MIME sniffing | Apply `helmet()` |
| Fetching user URLs without an allowlist | SSRF to internal services | Explicit allowlist, https only |
| `npm audit` only at release | Known CVEs ship | Run in CI, fail on HIGH |

## Skill Chaining

**Invoked by:** `typescript-code-review` when a diff touches authentication, authorization, payment, or PII

**Invokes:** none; this is the terminal security skill in the chain

**Can be invoked independently:** the user says "security review", "audit security", or targets a security-critical TypeScript or Node.js implementation

**Works alongside:** `security-audit-principles` for the shared workflow and severity model, `typescript-code-review` for the non-security findings
