# CLI Architecture Patterns

## 1. Decoupled Routing via `set_defaults`
When using `argparse` subparsers, avoid complex `if/elif` blocks in `main()`. Instead, map each sub-command directly to its handler function.

```python
# In main()
proc_parser = subparsers.add_parser("process")
proc_parser.set_defaults(func=handle_process)

# Execution is then a single line
args = parser.parse_args()
args.func(args)
```

## 2. Lazy Loading for Responsiveness
Move "heavy" imports (e.g., machine learning models, large libraries) inside the specific functions that require them. This ensures that metadata-only commands (like `--help` or version checks) remain instantaneous.

## 3. Systematic Failure via `sys.exit()`
In deterministic CLI tools, use `sys.exit(1)` immediately upon encountering a state that violates data integrity or schemas. This ensures that shell command chains (e.g., `cmd1 && cmd2`) stop execution correctly.

## 6. Absolute Package-Root Imports
Always use absolute imports referencing the package name (e.g., `from mypkg.module import ...`) even within the `src/` directory. This ensures consistent behavior across entry points and scripts.

## 7. Avoid Nested identical Quotes in f-strings
Python f-strings have limits on nesting identical quote types. For complex escaping (e.g., SQL queries), use `.format()` or perform the manipulation in a separate variable to improve readability and avoid syntax errors.