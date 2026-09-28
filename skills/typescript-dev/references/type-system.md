# TypeScript Type System Reference

Deep guidance for the type-system rules in `SKILL.md`. Load this when working with generics, conditional and mapped types, type guards, or branded and utility types.

## Contents

- [Generics](#generics)
- [Conditional types](#conditional-types)
- [Mapped types](#mapped-types)
- [Template literal types](#template-literal-types)
- [Utility types](#utility-types)
- [Type guards and narrowing](#type-guards-and-narrowing)
- [Branded types](#branded-types)
- [Constructive modeling](#constructive-modeling)
- [Exhaustiveness](#exhaustiveness)
- [Common mistakes](#common-mistakes)

## Generics

Generics tie an output type to an input type. Use one type parameter when the relationship is one-to-one, and constrain it when the body relies on members.

```typescript
// Constrain so T always has an id
function findById<T extends { id: string }>(items: T[], id: string): T | undefined {
    return items.find((item) => item.id === id);
}

// Preserve the input type rather than returning unknown
function first<T>(items: readonly T[]): T | undefined {
    return items[0];
}
```

Rules:

- Add a generic only when the caller's type must flow through. A generic used once and never related to another parameter is usually an `unknown` or a concrete type.
- Constrain with `extends` instead of casting inside the body.
- Default a type parameter when a common case exists: `<T = string>`.
- Never return `any` from a generic. Return `T` or a type derived from `T`.

## Conditional types

A conditional type picks a type based on an assignability test.

```typescript
type IsArray<T> = T extends readonly unknown[] ? true : false;
type ElementType<T> = T extends readonly (infer U)[] ? U : never;

type A = ElementType<string[]>; // string
type B = ElementType<number>;   // never
```

`infer` introduces a type variable inside the `extends` clause. Distribute over unions with a naked type parameter:

```typescript
type ToArray<T> = T extends unknown ? T[] : never;
type Result = ToArray<string | number>; // string[] | number[]
```

Wrap both sides in a tuple to prevent distribution when you want the whole union treated as one type:

```typescript
type ToArrayNonDist<T> = [T] extends [unknown] ? T[] : never;
```

## Mapped types

Mapped types build one type from another by iterating its keys.

```typescript
type PartialDeep<T> = {
    [K in keyof T]?: T[K] extends object ? PartialDeep<T[K]> : T[K];
};

type ReadonlyDeep<T> = {
    readonly [K in keyof T]: T[K];
};
```

Modifiers: `readonly` adds it, `-readonly` removes it, `?` makes properties optional, `-?` makes them required. Filter keys by remapping with `as`:

```typescript
type Getters<T> = {
    [K in keyof T as `get${Capitalize<string & K>}`]: () => T[K];
};
```

## Template literal types

Template literal types compose string types. They pair with mapped types for typed event systems and route tables.

```typescript
type EventName = "click" | "focus";
type Handler = `on${Capitalize<EventName>}`; // "onClick" | "onFocus"

type Route = `/api/${string}`;
const r: Route = "/api/users"; // ok
```

Watch out: a too-broad template like `` `/${string}` `` accepts almost anything. Constrain the interpolated part to a union of known literals when you need real checking.

## Utility types

Reach for a built-in before declaring a new interface. Deriving from an existing type keeps the two in sync.

| Utility | Use for |
| --- | --- |
| `Partial<T>` | All properties optional (patch DTOs) |
| `Required<T>` | All properties required |
| `Readonly<T>` | Immutable view of a type |
| `Pick<T, K>` | Subset of properties |
| `Omit<T, K>` | Everything except listed properties |
| `Record<K, V>` | Object with a known key set |
| `Exclude<T, U>` | Remove members from a union |
| `Extract<T, U>` | Keep members assignable to `U` |
| `NonNullable<T>` | Strip `null` and `undefined` |
| `ReturnType<F>` | Return type of a function type |
| `Parameters<F>` | Tuple of a function's parameters |
| `Awaited<T>` | Unwrap a promise type |
| `InstanceType<C>` | Instance type of a constructor |

```typescript
// Derive instead of duplicating
type UserSummary = Pick<User, "id" | "email">;
type UpdateUser = Partial<Omit<User, "id" | "createdAt">>;

async function load(): Promise<User> { /* ... */ }
type LoadedUser = Awaited<ReturnType<typeof load>>;
```

Custom utilities are worth writing when the same derivation appears three times:

```typescript
type ValueOf<T> = T[keyof T];
type DeepPartial<T> = { [K in keyof T]?: DeepPartial<T[K]> };
```

## Type guards and narrowing

A type guard is a function whose return type is a type predicate: `value is T`. It must actually verify the claim. A guard that returns `true` for values it does not check is worse than `as`, because callers trust it by name.

```typescript
interface Circle { kind: "circle"; radius: number }
interface Rect { kind: "rect"; width: number; height: number }
type Shape = Circle | Rect;

function isCircle(shape: Shape): shape is Circle {
    return shape.kind === "circle";
}
```

Narrowing tools, best first:

1. Discriminated union switch or `if` on the discriminant. The compiler narrows automatically.
2. The `in` operator: `"radius" in shape`.
3. `typeof` for primitives and `instanceof` for classes.
4. A user-defined guard for shapes the above cannot express.
5. `as`, only after validation.

```typescript
function area(shape: Shape): number {
    if ("radius" in shape) return Math.PI * shape.radius ** 2;
    return shape.width * shape.height;
}
```

Assertion functions (`asserts value is T`) throw on failure and narrow the caller's scope. Use them at boundaries where the invalid case is fatal.

## Branded types

A branded type is a primitive with a phantom property that makes it structurally distinct. Use one when two values share a primitive type but must not be interchangeable.

```typescript
type Brand<T, B extends string> = T & { readonly __brand: B };
type UserId = Brand<string, "UserId">;
type OrderId = Brand<string, "OrderId">;

function parseUserId(input: string): UserId {
    if (!isUuid(input)) throw new Error(`Invalid user id: ${input}`);
    return input as UserId; // earned after validation
}
```

Validate and cast once at the boundary. Inside the program, the type is trusted. Use the same `readonly __brand` shape everywhere so brands compose.

## Constructive modeling

Make illegal states unrepresentable by construction instead of guarding against them at runtime.

Non-empty collections via a variadic tuple:

```typescript
type NonEmpty<T> = [T, ...T[]];

function pickWinner(entries: NonEmpty<string>): string {
    return entries[Math.floor(Math.random() * entries.length)];
}

const isNonEmpty = <T>(arr: T[]): arr is NonEmpty<T> => arr.length > 0;
```

Even-length collections as pairs:

```typescript
type Pairs<T> = [T, T][];
```

Time ranges as start plus a non-negative duration, so a negative range cannot be written:

```typescript
type TimeRange = { start: Date; durationMs: number };
```

Keep the loose type when every operation on it is total. `sum(xs: number[])` is fine because `[]` sums to 0. Strengthen to `NonEmpty<T>` only where the loose type forces a `!`, an `arr[0] as T`, or a "should never happen" throw.

## Exhaustiveness

In the default arm of a switch over a discriminated union, assign the value to a `never` local. Adding a variant then fails the build at every unhandled switch.

```typescript
function area(shape: Shape): number {
    switch (shape.kind) {
        case "circle":
            return Math.PI * shape.radius ** 2;
        case "rect":
            return shape.width * shape.height;
        default: {
            const _exhaustive: never = shape;
            return _exhaustive;
        }
    }
}
```

Use the return style in a value-returning switch, and a void style (`void _exhaustive;`) in a statement switch.

## Common mistakes

| Mistake | Fix |
| --- | --- |
| `any` in a public signature | Use a generic `<T>` or `unknown` |
| `as` to force a shape | Validate, then narrow, then cast if still needed |
| `object` as a parameter type | Use `Record<string, unknown>` or a concrete interface |
| `enum` for a closed set of strings | Use a string literal union |
| Re-declaring a shape that a schema or generated type already defines | Derive with `Pick`, `Omit`, or `z.infer` |
| `Function` as a callback type | Use `(arg: T) => R` |
| Optional field that is always read | Make it required or model with a discriminant |
| Deep optional chain with no fallback | Add `?? fallback` and handle the missing case |
