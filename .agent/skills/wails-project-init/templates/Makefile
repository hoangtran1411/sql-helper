.PHONY: dev build test lint coverage clean help bindings check fmt fmt-check

# Default target
help:
	@echo "{{PROJECT_NAME}} - Available commands:"
	@echo "  make dev          - Run in development mode (wails3 dev)"
	@echo "  make bindings     - Generate frontend bindings (wails3 generate bindings)"
	@echo "  make build        - Build production binary (wails3 build)"
	@echo "  make test         - Run all tests"
	@echo "  make lint         - Run golangci-lint"
	@echo "  make coverage     - Run tests with coverage report"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make deps         - Download dependencies"

# Development mode
dev:
	wails3 dev

# Generate Wails v3 frontend bindings
bindings:
	wails3 generate bindings

# Build production binary
build:
	wails3 build

# Build for Windows
build-windows:
	wails3 build -platform windows/amd64

# Build for macOS
build-darwin:
	wails3 build -platform darwin/universal

# Build for Linux
build-linux:
	wails3 build -platform linux/amd64

# Run all tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run ./...

# Run tests with coverage (space instead of = for Windows PowerShell compatibility)
coverage:
	go test -v -coverprofile coverage.out ./...
	go tool cover -func=coverage.out
	@echo ""
	@echo "To view HTML report, run: go tool cover -html=coverage.out"

# Generate coverage HTML report
coverage-html: coverage
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	rm -rf build/
	rm -f coverage.out coverage.html internal.out

# Download dependencies
deps:
	go mod download
	go mod tidy

# Install development tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/wailsapp/wails/v3/cmd/wails3@latest

# Format code
fmt:
	go fmt ./...
	goimports -w .

# Check if code is properly formatted
fmt-check:
	@test -z "$$(gofmt -l .)" || (echo "Code is not formatted. Run 'make fmt'" && exit 1)

# Run all checks (lint + test)
check: lint test
	@echo "All checks passed!"
