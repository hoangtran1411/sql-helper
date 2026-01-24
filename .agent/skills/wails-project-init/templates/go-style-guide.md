---
trigger: always_on
---

# Go Style Guide - {{PROJECT_NAME}}

This project is a **{{PROJECT_DESCRIPTION}}** built with:
- **Wails v2** for the Desktop GUI (Windows)
- **Internal packages** for modular business logic

---

## Code Style

- Ensure all Go code is formatted using `gofmt` or `goimports`. Run `golangci-lint run ./...` before committing.
- Adhere to [Effective Go](https://go.dev/doc/effective_go) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).
- Organize core processing logic in `internal/` packages.
- Keep Wails-specific code (`app.go`, `main.go`) in the root package.

## Project Structure

```
{{PROJECT_NAME}}/
  main.go              - Wails entry point
  app.go               - Wails app bindings (Go methods exposed to JS)
  internal/
    {{DOMAIN}}/        - Domain-specific logic
      {{DOMAIN}}.go      - Core implementation
      {{DOMAIN}}_test.go - Unit tests
  frontend/            - Wails frontend (HTML/CSS/JS)
  build/               - Build assets and output binaries
  wails.json           - Wails project configuration
```

## Error Handling

- Always wrap errors using `%w`: `fmt.Errorf("context: %w", err)`.
- Implement fail fast logic using guard clauses to minimize indentation.
- Handle close errors in defer statements appropriately.
- Do not log and return the same error.
- Choose one responsibility per layer.

## Context and Concurrency

- Functions performing I/O or long-running operations MUST accept `context.Context` as the first argument.
- Use `sync.WaitGroup` for coordinating concurrent operations if needed.
- Use `errgroup.Group` for concurrent tasks with error propagation.

## Wails Integration

- All Wails-bound methods must be on the `*App` struct and be exported (PascalCase).
- Use `runtime.EventsEmit()` for progress updates to frontend.
- Return structs with `json` tags for frontend consumption.
- Use `runtime.OpenFileDialog()` and `runtime.SaveFileDialog()` for native file selection.
- Use `runtime.ClipboardSetText()` for clipboard operations.

## Documentation

- Every exported function, variable, and type must have clear documentation comments.
- Document edge cases and design decisions in comments.

## Testing

- Prioritize Table-driven tests combined with `t.Run` for comprehensive test coverage.
- Target minimum 70% code coverage (CI enforced for `internal/` packages).
- Use `go test ./... -v` to run all tests.
- Use `make test` for convenient test execution.

## Linting (golangci-lint)

Recommended linters:
- `errcheck`, `gosimple`, `govet`, `ineffassign`, `staticcheck`, `unused`
- `gofmt`, `goimports`, `misspell`, `gocritic`, `gosec`

Excluded patterns:
- `frontend/` and `build/` directories should be excluded

## Efficiency and Tone

- Avoid greetings, apologies, or meta-commentary; focus strictly on code and execution logs.
- Provide code as minimal diffs/blocks whenever possible.

---

## Go Idioms & Professional Practices

### Naming Conventions (Idiomatic Go)

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

### Function Design

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

### Error Handling Idioms

- Never ignore errors explicitly:

```go
_ = f.Close() // ❌ unless justified in comment
```

- Wrap errors only at package boundaries:

```go
return nil, fmt.Errorf("operation failed: %w", err)
```

- Do not wrap errors multiple times in the same layer.

- Prefer `errors.Is` / `errors.As` for comparisons:

```go
if errors.Is(err, ErrNotFound) { ... }
```

- Avoid sentinel errors unless necessary; prefer typed errors:

```go
type ErrInvalidInput struct {
    Field string
}
```

### Package Design & Boundaries

- `internal` packages must be:
  - UI-agnostic
  - Framework-agnostic (no Wails imports)

- Each package should expose minimal API surface:

```go
// good
func Process(ctx context.Context, input string) (*Result, error)

// avoid exposing helpers
```

- Avoid circular dependencies at all cost.
- If two packages depend on each other → redesign.

### Context Usage Idioms

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

### Concurrency Patterns

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

### Struct & Interface Idioms

- Accept interfaces, return concrete types:

```go
func NewProcessor(r io.Reader) *Processor
```

- Interfaces should be small (1–3 methods):

```go
type Reader interface {
    Read() ([]byte, error)
}
```

- Do not define interfaces prematurely.

### Zero Value Philosophy

- Design structs so zero value is usable:

```go
var p Processor // should work
```

- Avoid constructors unless needed for invariants.

- Prefer empty slices over nil slices for JSON output:

```go
items := make([]Item, 0)
```

### Slice & Map Best Practices

- Pre-allocate when size is known:

```go
items := make([]Item, 0, estimatedSize)
```

- Check map existence properly:

```go
v, ok := m[key]
```

- Do not modify slices while ranging over them.

### Testing Idioms

- Test behavior, not implementation.
- Table-driven tests with descriptive names:

```go
name: "valid input returns expected result"
```

- Avoid `t.Fatal` inside loops.
- Use `cmp.Diff` or `reflect.DeepEqual` consistently.
- Tests must not depend on execution order.

### Logging (If Used)

- Do not log inside core business logic.
- Log at boundaries (UI / CLI / App layer).
- Logs must be structured and actionable.

### Comments & Documentation

- Comments explain **why**, not **what**.
- Avoid redundant comments:

```go
i++ // increment i ❌
```

- Exported comments must start with identifier name.
- Use TODO with owner & reason:

```go
// TODO(username): implement feature X
```

### Defensive Programming

- Validate inputs at package boundary.
- Never trust external data types.
- Fail fast on schema mismatch.
- Prefer explicit errors over silent correction.

### Build & Tooling Practices

- `go.mod` must be tidy:

```bash
go mod tidy
```

- CI must fail on:
  - lint
  - test
  - formatting

- Avoid build tags unless justified.

### AI Agent Enforcement Rule

When generating Go code for this project:

- Prefer clarity over cleverness.
- Prefer explicit code over abstractions.
- Prefer idiomatic Go over patterns from other languages (Java/C#/JS).
- If unsure, follow Effective Go first.

---

## Reference and Resource Mapping

### Official Go Documentation

- Reference: https://go.dev/doc
- Guideline: Primary source for Go syntax, tooling, modules, and release notes.

### Effective Go

- Reference: https://go.dev/doc/effective_go
- Guideline: Defines idiomatic Go practices (naming, control flow, error handling). All `internal/` packages must comply.

### Go Modules

- Reference: https://go.dev/ref/mod
- Guideline: Dependency management using `go.mod` and `go.sum`. Avoid unnecessary `replace` directives in production.

### Go Testing

- Reference: https://go.dev/doc/testing
- Guideline: Standard patterns for unit tests, table-driven tests, benchmarks, and coverage analysis.

### Go Context

- Reference: https://pkg.go.dev/context
- Guideline: Mandatory for cancellation, timeouts, and request-scoped values in I/O and concurrent operations.

### Go Error Handling

- Reference: https://go.dev/blog/error-handling-and-go
- Guideline: Errors are values. Always return and wrap errors; avoid panic in business logic.

### Go Concurrency Patterns

- Reference: https://go.dev/doc/effective_go#concurrency
- Guideline: Use goroutines and channels deliberately. Avoid shared mutable state unless properly synchronized.

### Go Standard Library

- Reference: https://pkg.go.dev/std
- Guideline: Prefer the Go standard library before introducing third-party dependencies.

### Wails Desktop App

- Reference: https://github.com/wailsapp/wails
- Guideline: Follow Wails v2 patterns for Go-to-frontend binding. Use event system for real-time progress updates.

### Linting

- Reference: https://github.com/golangci/golangci-lint
- Guideline: Run `golangci-lint run ./...` before committing. Fix all issues to pass CI.

---

## Project-Specific Sections

> Add any domain-specific coding guidelines below.

