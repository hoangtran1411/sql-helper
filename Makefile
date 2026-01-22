.PHONY: dev build test lint coverage clean help

# Default target
help:
	@echo "SQL Helper - Available commands:"
	@echo "  make dev       - Run in development mode"
	@echo "  make build     - Build production binary"
	@echo "  make test      - Run all tests"
	@echo "  make lint      - Run golangci-lint"
	@echo "  make coverage  - Run tests with coverage report"
	@echo "  make clean     - Remove build artifacts"
	@echo "  make deps      - Download dependencies"

# Development mode
dev:
	wails dev

# Build production binary
build:
	wails build

# Build for Windows
build-windows:
	wails build -platform windows/amd64

# Build for macOS
build-darwin:
	wails build -platform darwin/universal

# Build for Linux
build-linux:
	wails build -platform linux/amd64

# Run all tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run ./...

# Run tests with coverage
coverage:
	go test -v -coverprofile=coverage.out ./...
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
	rm -f coverage.out coverage.html

# Download dependencies
deps:
	go mod download
	go mod tidy

# Install development tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

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
