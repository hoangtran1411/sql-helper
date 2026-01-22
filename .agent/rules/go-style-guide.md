---
trigger: always_on
---

# Go Style Guide - Excel Image Importer

This project is an **Excel image insertion utility** built with:
- **Wails v2** for the Desktop GUI (Windows)
- **Excelize v2** for Excel file manipulation
- **Worker Pool** for parallel image processing

---

## Code Style

- Ensure all Go code is formatted using gofmt or goimports. Run golangci-lint run ./... before committing.
- Adhere to Effective Go and Go Code Review Comments.
- Organize core processing logic in internal/engine/ package.
- Keep Wails-specific code (app.go, updater.go, main.go) in the root package.

## Project Structure

GoExcelImageImporter/
  main.go              - Wails entry point
  app.go               - Wails app bindings (Go methods exposed to JS)
  updater.go           - Auto-update functionality from GitHub Releases
  internal/
    engine/            - Core processing logic
      processor.go       - Image-to-Excel engine with Worker Pool
      processor_test.go  - Unit tests for processor
  frontend/            - Wails frontend (HTML/CSS/JS)
  build/               - Build assets and output binaries
  docs/                - Documentation and roadmap
  wails.json           - Wails project configuration

## Error Handling

- Always wrap errors using %w: fmt.Errorf(context: %w, err). Critical for tracing Excel/Image I/O errors.
- Implement fail fast logic using guard clauses to minimize indentation.
- Handle close errors in defer statements: use defer f.Close() pattern for excelize.File.

## Context and Concurrency

- Functions performing I/O or long-running operations MUST accept context.Context as the first argument.
- Use Worker Pool pattern with chan Job and chan Result for parallel image loading.
- Use sync.WaitGroup for coordinating workers.

## Wails Integration

- All Wails-bound methods must be on the *App struct and be exported (PascalCase).
- Use runtime.EventsEmit() for progress updates to frontend.
- Return structs with json tags for frontend consumption (e.g., Config, ProcessResult).
- Use runtime.OpenFileDialog() and runtime.OpenDirectoryDialog() for native file/folder selection.

## Excel Processing (Excelize)

- Use excelize.OpenFile() to read existing Excel files.
- Use f.Rows() iterator for memory-efficient reading of large files.
- Always call defer f.Close() after opening an Excel file.

## Image Handling

- Support common image formats: .jpg, .jpeg, .png, .gif, .webp.
- Import blank image decoders for jpeg, png, and golang.org/x/image/webp.
- Use image.DecodeConfig() to get image dimensions without fully decoding.
- Calculate scale factors to fit images within cell dimensions.

## Documentation

- Every exported function, variable, and type must have clear documentation comments.
- Document edge cases and design decisions in comments.

## Testing

- Prioritize Table-driven tests combined with t.Run for comprehensive test coverage.
- Target minimum 40% code coverage (CI enforced).
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