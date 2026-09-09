---
name: Go Test & Makefile
description: Standard setup for Go projects with Makefile automation and Unit Testing patterns.
---

# Go Test & Makefile Skill

This skill provides a setup for Go projects, enabling automation of common tasks (build, test, lint, format) via a `Makefile` and standardized Unit Testing patterns.

## When to Use
- Adding CI/CD readiness or build automation to a Go project.
- Standardizing build, lint, and test commands across developer environments.
- Enforcing test coverage gates (e.g. 70% coverage on `internal/` packages).

## Components

### 1. Makefile
A `Makefile` supporting:
- **Build**: Compiles binaries cross-platform (`build-windows`, `build-darwin`, `build-linux`) and integrates Wails v3 commands (`wails3 dev`, `wails3 build`, `wails3 generate bindings`).
- **Test**: Runs unit tests with verbose output (`make test`).
- **Coverage**: Generates coverage profiles without `=` flag for Windows PowerShell compatibility (`go test -v -coverprofile coverage.out ./...`).
- **Lint**: Runs `golangci-lint` (`make lint`).
- **Format & Check**: Formats code with `gofmt`/`goimports` and checks style via `make fmt-check`.
- **Clean**: Removes build outputs and coverage profiles.

### 2. Test Template
A `main_test.go` template showing how to:
- Write idiomatic table-driven tests with `t.Run`.
- Setup/Teardown with `t.Cleanup()`.
- Validate error and success branches.

## Usage

### Step 1: Add Makefile
Copy `templates/Makefile` to your project root. Adjust target paths if necessary.

### Step 2: Add Test File
Copy `templates/main_test.go` to your package directory (e.g., `internal/excel` or `internal/sql`) and rename it to `<feature>_test.go`.

### Step 3: Run Commands
```bash
# Run unit tests
make test

# Generate and view coverage report
make coverage
make coverage-html

# Format and verify lint
make fmt
make check
```

## Prerequisites
- **Go**: 1.27+
- **Make**:
    - **Linux/macOS**: Pre-installed.
    - **Windows**: Install via `choco install make` or use Git Bash / MSYS2.
- **GolangCI-Lint**: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
- **Wails v3 CLI**: `go install github.com/wailsapp/wails/v3/cmd/wails3@latest`
