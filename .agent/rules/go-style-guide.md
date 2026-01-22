---
trigger: always_on
---

# Go Style Guide - SQL Helper

This project is an **Excel to SQL INSERT statement converter** built with:
- **Wails v2** for the Desktop GUI (Windows)
- **Excelize v2** for Excel file parsing
- **Internal packages** for modular business logic

---

## Code Style

- Ensure all Go code is formatted using gofmt or goimports. Run golangci-lint run ./... before committing.
- Adhere to Effective Go and Go Code Review Comments.
- Organize core processing logic in internal/ packages (excel/, sql/).
- Keep Wails-specific code (app.go, updater.go, main.go) in the root package.

## Project Structure

sql-helper/
  main.go              - Wails entry point
  app.go               - Wails app bindings (Go methods exposed to JS)
  updater.go           - Auto-update functionality from GitHub Releases
  internal/
    excel/             - Excel file parsing logic
      parser.go          - Excel file reading using Excelize
      parser_test.go     - Unit tests for parser
    sql/               - SQL generation logic
      generator.go       - SQL INSERT statement generator
      generator_test.go  - Unit tests for generator
      formatter.go       - Value formatting for SQL
      formatter_test.go  - Unit tests for formatter
  frontend/            - Wails frontend (HTML/CSS/JS)
  build/               - Build assets and output binaries
  wails.json           - Wails project configuration

## Error Handling

- Always wrap errors using %w: fmt.Errorf("context: %w", err). Critical for tracing Excel I/O errors.
- Implement fail fast logic using guard clauses to minimize indentation.
- Handle close errors in defer statements: use defer f.Close() pattern for excelize.File.

## Context and Concurrency

- Functions performing I/O or long-running operations MUST accept context.Context as the first argument.
- Use sync.WaitGroup for coordinating concurrent operations if needed.

## Wails Integration

- All Wails-bound methods must be on the *App struct and be exported (PascalCase).
- Use runtime.EventsEmit() for progress updates to frontend.
- Return structs with json tags for frontend consumption (e.g., ExcelResult, SheetData).
- Use runtime.OpenFileDialog() and runtime.SaveFileDialog() for native file selection.
- Use runtime.ClipboardSetText() for clipboard operations.

## Excel Processing (Excelize)

- Use excelize.OpenFile() to read existing Excel files.
- Use f.Rows() iterator for memory-efficient reading of large files.
- Always call defer f.Close() after opening an Excel file.
- Parse headers from first row, data rows from subsequent rows.

## SQL Generation

- Generate VALUES clause for SQL INSERT statements.
- Support number columns detection to avoid quoting numeric values.
- Use proper escaping for string values (single quotes).
- Implement Find and Replace functionality for data transformation.

## Documentation

- Every exported function, variable, and type must have clear documentation comments.
- Document edge cases and design decisions in comments.

## Testing

- Prioritize Table-driven tests combined with t.Run for comprehensive test coverage.
- Target minimum 70% code coverage (CI enforced for internal/ packages).
- Use go test ./... -v to run all tests.
- Use make test for convenient test execution.

## Linting (golangci-lint)

Recommended linters:
- errcheck, gosimple, govet, ineffassign, staticcheck, unused
- gofmt, goimports, misspell, gocritic, gosec

Excluded patterns:
- frontend/ and build/ directories should be excluded

## Efficiency and Tone

- Avoid greetings, apologies, or meta-commentary; focus strictly on code and execution logs.
- Provide code as minimal diffs/blocks whenever possible.

---

## Reference and Resource Mapping

### Wails Desktop App
- Reference: https://github.com/wailsapp/wails
- Guideline: Follow Wails v2 patterns for Go-to-frontend binding. Use event system for real-time progress updates.

### Excel Processing
- Reference: https://github.com/xuri/excelize
- Guideline: Use Excelize v2 for all Excel operations. Use streaming API for large files.

### Linting
- Reference: https://github.com/golangci/golangci-lint
- Guideline: Run golangci-lint run ./... before committing. Fix all issues to pass CI.