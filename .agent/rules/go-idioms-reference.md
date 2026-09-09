# Go Idioms Reference - SQL Helper

> **Full Reference Document** - Contains detailed idioms, code examples, and best practices.  
> For compact core rules, see `go-style-guide.md`

---

## Naming Conventions (Idiomatic Go)

- Use short but meaningful names, scoped by context:
  - `r`, `w`, `ctx`, `db`, `tx`, `cfg` are acceptable in small scopes.
  - Avoid `data`, `info`, `obj`, `temp`, `value` unless unavoidable.

- Prefer noun-based names for structs, verb-based names for functions:
  - `Parser.Parse()`, `Generator.Generate()`

- Boolean names should read naturally:
  - `isValid`, `hasHeader`, `enableStreaming`

- Avoid stuttering:
  - ❌ `sql.SQLGenerator`
  - ✅ `sql.Generator`

---

## Function Design

- Prefer small functions (≤ 40 lines).
- One function = one responsibility.
- Avoid flags that change behavior dramatically:

```go
// bad
func Parse(file string, strict bool)

// good
func ParseStrict(file string)
func ParseLenient(file string)
```

- Return early (guard clauses):

```go
if err != nil {
    return nil, err
}
```

---

## Error Handling Idioms

- Never ignore errors explicitly:

```go
_ = f.Close() // ❌ unless justified in comment
```

- Wrap errors only at package boundaries:

```go
return nil, fmt.Errorf("excel parse failed: %w", err)
```

- Do not wrap errors multiple times in the same layer.

- Prefer `errors.Is` / `errors.As` for comparisons:

```go
if errors.Is(err, excelize.ErrCellName) { ... }
```

- Avoid sentinel errors unless necessary; prefer typed errors:

```go
type ErrInvalidHeader struct {
    Column string
}
```

---

## Package Design & Boundaries

- `internal` packages must be:
  - UI-agnostic
  - Framework-agnostic (no Wails imports)

- Each package should expose minimal API surface:

```go
// good
func Parse(ctx context.Context, path string) (*Result, error)

// avoid exposing helpers
```

- Avoid circular dependencies at all cost.
- If two packages depend on each other → redesign.

---

## Context Usage Idioms

- `context.Context` must:
  - Be the first parameter
  - Never be stored in struct fields

- Do not pass nil context:

```go
ctx := context.Background()
```

- Respect cancellation in loops:

```go
select {
case <-ctx.Done():
    return ctx.Err()
default:
}
```

---

## Concurrency Patterns

- Prefer worker pool over unbounded goroutines.
- Always define ownership of goroutines:
  - Who starts?
  - Who stops?

- Use `errgroup.Group` for concurrent tasks with error propagation.

- Channels should have clear direction:

```go
func worker(in <-chan Job, out chan<- Result)
```

- Avoid closing channels you did not create.

---

## Struct & Interface Idioms

- Accept interfaces, return concrete types:

```go
func NewParser(r io.Reader) *Parser
```

- Interfaces should be small (1–3 methods):

```go
type RowReader interface {
    Next() ([]string, error)
}
```

- Do not define interfaces prematurely.

---

## Zero Value Philosophy

- Design structs so zero value is usable:

```go
var g Generator // should work
```

- Avoid constructors unless needed for invariants.

- Prefer empty slices over nil slices for JSON output:

```go
rows := make([]Row, 0)
```

---

## Slice & Map Best Practices

- Pre-allocate when size is known:

```go
rows := make([]Row, 0, estimatedRows)
```

- Check map existence properly:

```go
v, ok := m[key]
```

- Do not modify slices while ranging over them.

---

## Testing Idioms

- Test behavior, not implementation.
- Table-driven tests with descriptive names:

```go
name: "numeric column without quotes"
```

- Avoid `t.Fatal` inside loops.
- Use `cmp.Diff` or `reflect.DeepEqual` consistently.
- Tests must not depend on execution order.

---

## Logging (If Used)

- Do not log inside core business logic.
- Log at boundaries (UI / CLI / App layer).
- Logs must be structured and actionable.

---

## Comments & Documentation

- Comments explain **why**, not **what**.
- Avoid redundant comments:

```go
i++ // increment i ❌
```

- Exported comments must start with identifier name.
- Use TODO with owner & reason:

```go
// TODO(tuakmouo): support multi-sheet merge
```

---

## SQL-Specific Idioms (Project Relevant)

- Never concatenate SQL identifiers blindly.
- Keep SQL generation deterministic (same input → same output).
- Separate:
  - data normalization
  - SQL formatting
- Avoid hidden mutation during formatting.

---

## Defensive Programming

- Validate inputs at package boundary.
- Never trust Excel data types.
- Fail fast on schema mismatch.
- Prefer explicit errors over silent correction.

---

## Build & Tooling Practices

- `go.mod` must be tidy:

```bash
go mod tidy
```

- CI must fail on:
  - lint
  - test
  - formatting

- Avoid build tags unless justified.

---

## Reference Links

### Official Go Documentation

- Reference: <https://go.dev/doc>
- Primary source for Go syntax, tooling, modules.

### Effective Go

- Reference: <https://go.dev/doc/effective_go>
- Idiomatic Go practices. All `internal/` packages must comply.

### Go Modules

- Reference: <https://go.dev/ref/mod>
- Dependency management. Avoid unnecessary `replace` directives.

### Go Testing

- Reference: <https://go.dev/doc/testing>
- Standard patterns for unit tests, benchmarks, coverage.

### Go Context

- Reference: <https://pkg.go.dev/context>
- Mandatory for cancellation and timeouts in I/O operations.

### Go Error Handling

- Reference: <https://go.dev/blog/error-handling-and-go>
- Errors are values. Wrap errors; avoid panic in business logic.

### Go Concurrency

- Reference: <https://go.dev/doc/effective_go#concurrency>
- Use goroutines and channels deliberately.

### Go Standard Library

- Reference: <https://pkg.go.dev/std>
- Prefer stdlib before third-party dependencies.

### Wails Desktop App

- Reference: <https://v3.wails.io>
- GitHub: <https://github.com/wailsapp/wails>
- Follow Wails v3 service patterns (`application.NewService`, `application.Get().*`) for Go-to-frontend communication.

### Excelize

- Reference: <https://github.com/xuri/excelize>
- Use Excelize v2 for all Excel operations.

### Linting

- Reference: <https://github.com/golangci/golangci-lint>
- Run before committing. Fix all issues to pass CI.
