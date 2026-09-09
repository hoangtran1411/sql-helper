---
trigger: always_on
---

# Go Style Guide - SQL Helper

> **Core Rules** - For full idioms reference, see `go-idioms-reference.md`

This project is an **Excel to SQL INSERT statement converter** built with:

- **Wails v3** (`v3.0.0-beta.15`) for the Desktop GUI (Windows)
- **Go 1.27**
- **Excelize v2** for Excel file parsing & streaming
- **Internal packages** for modular business logic

---

## Code Style

- Format with `gofmt`/`goimports`. Run `golangci-lint run ./...` before commit.
- Adhere to [Effective Go](https://go.dev/doc/effective_go).
- Core logic in `internal/` packages. Wails service bindings in root package.

## Project Structure

```
sql-helper/
  main.go              - Wails v3 entry point (application.New, WebviewWindowOptions)
  app.go               - Backend application service bound to frontend
  updater.go           - Auto-update from GitHub Releases (application.Get())
  internal/
    excel/             - Excel parsing & O(1) streaming (parser.go)
    sql/               - SQL generation & formatting (generator.go, formatter.go)
  frontend/            - Vanilla HTML/CSS/JS frontend
    bindings/          - Wails v3 generated ES module bindings
  build/               - Build output
```

## Error Handling

- Wrap errors: `fmt.Errorf("context: %w", err)`
- Guard clauses for fail-fast
- Do not log and return the same error
- One responsibility per layer

## Wails v3 Integration

- Service registration: `application.NewService(appService)` in `main.go`
- Public methods on `*App` struct (PascalCase) exposed to frontend bindings
- Runtime API access via `application.Get()`:
  - Events: `application.Get().Event.Emit("eventName", data)`
  - Dialogs: `application.Get().Dialog.OpenFileWithOptions()` & `SaveFileWithOptions()`
  - Clipboard: `application.Get().Clipboard.SetText(text)`
  - Browser: `application.Get().Browser.OpenURL(url)`
  - Quit: `application.Get().Quit()`
- Frontend bindings: Generated via `wails3 generate bindings`, imported as ES modules:
  - `import * as App from "./bindings/github.com/hoangtran1411/sql-helper/app.js"`
  - `import { Events } from "@wailsio/runtime"`
- Return structs with `json` tags for JSON serialization

## Streaming & Performance Architecture

- O(1) Memory Streaming: Use `excel.IterateSheet()` row iterator with `bufio.Writer` for direct file export without loading entire dataset into RAM
- UI Preview: Limit preview data to 100 rows (`excel.GetPreview()`) to keep frontend DOM lightweight and responsive
- Batching: Support customizable SQL batch sizes (e.g. 1000 rows per `INSERT INTO`) to optimize database ingestion

## Testing & Linting

- Table-driven tests with `t.Run`
- Target 70% coverage gate specifically on `internal/` packages (`./internal/...`)
- Windows PowerShell compatibility: Always separate `-coverprofile` with space, not `=`, e.g. `go test -coverprofile coverage.out ./...`
- Commands: `make test`, `make lint`, `make check`, `make coverage`

---

## AI Agent Rules (Critical)

### Enforcement

- Prefer clarity over cleverness
- Prefer idiomatic Go over Java/C#/JS patterns
- If unsure, follow Effective Go first

### Context Accuracy

- Documentation links ≠ guarantees of correctness
- For external APIs: prefer explicit function signatures in context
- State assumptions when context is missing

### Library Version Awareness

- Check `go.mod` for actual versions before suggesting APIs (Wails v3, Excelize v2, Go 1.27)
- LLMs hallucinate APIs for newer features not in training data
- Note Wails v3 API differences from v2: `application.Get()` replaces context-based `runtime.*` calls

### Context Engineering

- Right context at right time, not all docs at once
- Reference existing codebase patterns first
- State missing context rather than guessing

---

## Quick Reference Links

- [Effective Go](https://go.dev/doc/effective_go)
- [Wails v3 Docs](https://v3.wails.io)
- [Wails GitHub](https://github.com/wailsapp/wails)
- [Excelize](https://github.com/xuri/excelize)
- [golangci-lint](https://github.com/golangci/golangci-lint)

> **Full Reference:** See `go-idioms-reference.md` for detailed idioms, code examples, and best practices.
