# TypeScript Configuration Reference

tsconfig and build guidance for the `typescript-dev` skill. Load this when setting up a project, fixing type-check failures, or tuning build performance.

## Strict baseline

```json
{
  "compilerOptions": {
    "strict": true,
    "noUncheckedIndexedAccess": true,
    "exactOptionalPropertyTypes": true,
    "noImplicitOverride": true,
    "noImplicitReturns": true,
    "noFallthroughCasesInSwitch": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "forceConsistentCasingInFileNames": true
  }
}
```

`strict` turns on `noImplicitAny`, `strictNullChecks`, `strictFunctionTypes`, `strictBindCallApply`, `strictPropertyInitialization`, `noImplicitThis`, and `alwaysStrict`.

The two flags worth adding beyond `strict`:

- `noUncheckedIndexedAccess` makes `arr[i]` and `obj[key]` include `undefined`. It catches real out-of-bounds bugs, at the cost of more guards.
- `exactOptionalPropertyTypes` separates `prop?: string` from `prop: string | undefined`, so a missing key and an explicit `undefined` are not interchangeable.

## Module resolution

| Target | `module` | `moduleResolution` | Notes |
| --- | --- | --- | --- |
| Node.js | `NodeNext` | `NodeNext` | Honors `exports` and `imports`; requires file extensions or package exports |
| Bundler (Vite, esbuild) | `ESNext` | `bundler` | Allows extensionless imports; pair with `isolatedModules` |
| Legacy | `CommonJS` | `node` | Avoid for new projects |

Set `target` to `ES2022` or later and `lib` to match the runtime (`["ES2022"]` for Node, add `DOM`/`DOM.Iterable` for the browser).

## Project references

Split a monorepo into composite projects so `tsc` only rebuilds what changed.

```json
// Root tsconfig.json
{
  "files": [],
  "references": [
    { "path": "./packages/shared" },
    { "path": "./packages/frontend" },
    { "path": "./packages/backend" }
  ]
}
```

```json
// packages/shared/tsconfig.json
{
  "compilerOptions": {
    "composite": true,
    "outDir": "./dist",
    "rootDir": "./src",
    "declaration": true,
    "declarationMap": true
  },
  "include": ["src/**/*"]
}
```

Downstream packages add `{ "references": [{ "path": "../shared" }] }` and run a solution-wide check with `tsc --build`. Do not run `tsc --build` in CI just to type-check; use a non-emitting solution check so CI does not rewrite artifacts.

## Path mapping

```json
{
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@/*": ["src/*"],
      "@components/*": ["src/components/*"]
    }
  }
}
```

Path aliases must also be configured in the bundler or Node loader (`vite-tsconfig-paths`, `tsconfig-paths`) or the runtime will not resolve them.

## Build performance

Slow type checks come from checking dependency declarations and from rechecking unchanged files. Address both.

| Option | Effect |
| --- | --- |
| `skipLibCheck: true` | Skips `.d.ts` checking in dependencies; the single biggest speedup |
| `incremental: true` | Writes `.tsbuildinfo` and rechecks only changed files |
| `composite: true` | Enables project references and incremental builds |
| `assumeChangesOnlyAffectDirectDependencies: true` | Skips rebuilding transitive dependents; faster but less safe |
| `isolatedModules: true` | Guarantees each file can be transpiled alone; required for most bundlers |
| `noEmit: true` | Type-check only; no output to write |

## Diagnosing slow checks

```bash
tsc --diagnostics                 # counts and total time
tsc --extendedDiagnostics         # per-phase breakdown
tsc --generateTrace trace         # emits an event trace
npx @typescript/analyze-trace trace   # ranks the files that dominate check time
```

Common causes of a slow check: a recursive type several levels deep, a large union that forces repeated work, `@types` packages pulled in by `types` in the config, and a `paths` mapping that makes resolution hit the filesystem on every import.

## Multiple configurations

Use `extends` to share a base and specialize per context.

```json
// tsconfig.build.json
{
  "extends": "./tsconfig.json",
  "compilerOptions": { "sourceMap": false, "declaration": true },
  "exclude": ["**/*.test.ts", "**/*.spec.ts"]
}
```

```json
// tsconfig.test.json
{
  "extends": "./tsconfig.json",
  "compilerOptions": { "types": ["vitest/globals", "node"] },
  "include": ["src/**/*.test.ts", "src/**/*.spec.ts"]
}
```

## Ambient and module declarations

```typescript
// src/types/global.d.ts
declare global {
    interface Window {
        myApp: { version: string };
    }
    namespace NodeJS {
        interface ProcessEnv {
            DATABASE_URL: string;
            NODE_ENV: "development" | "production" | "test";
        }
    }
}
export {};

// src/types/modules.d.ts
declare module "*.css" {
    const classes: Record<string, string>;
    export default classes;
}
```

Prefer `declare module` over `@ts-ignore` for genuinely untyped packages.

## Quick reference

| Option | Purpose |
| --- | --- |
| `strict` | All strict checks |
| `noUncheckedIndexedAccess` | `arr[i]` includes `undefined` |
| `exactOptionalPropertyTypes` | Separates missing from `undefined` |
| `composite` | Project references, incremental |
| `incremental` | Recheck changed files only |
| `skipLibCheck` | Skip dependency `.d.ts` checks |
| `isolatedModules` | Each file transpiles alone |
| `moduleResolution` | How imports are resolved |
| `paths` | Import aliases |
| `noEmit` | Type-check only |
| `declaration` | Emit `.d.ts` |
| `sourceMap` | Emit source maps |
