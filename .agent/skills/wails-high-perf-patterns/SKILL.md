---
name: Wails High Performance Patterns
description: Reusable patterns for building high-performance Go/Wails apps, focusing on Concurrency, Memory Efficiency, and Benchmarking.
---

# Wails High Performance Patterns

This skill encapsulates the "secret sauce" for high-performance Wails applications, specifically designed for data-intensive tasks like Image Processing and Batch Excel operations.

## When to Use
- **Migrating Legacy Apps**: Replacing slow C#/.NET apps with Go.
- **Heavy Processing**: Handling thousands of files or large datasets.
- **Responsive UI**: Ensuring the frontend remains smooth while the backend crunches numbers.

## Core Patterns

### 1. Worker Pool (Concurrency)
Instead of spawning execution logic linearly or spawning unlimited goroutines (which causes thrashing), use a Worker Pool.
- **Concept**: Fixed number of workers (e.g., CPU count * 2) consuming from a Job Channel.
- **Benefit**: Controlled resource usage, maximum CPU throughput.
- **UI Integration**: Send progress events from the "Collector" phase, not individual workers, to avoid flooding the frontend event loop.

### 2. Stream Processing (Memory)
Avoid loading entire files into memory.
- **Excel**: Use `rows.Next()` (Iterator) instead of `GetRows()` (Load all).
- **Images**: Decode only `Config` headers first to check dimensions before loading full pixel data if possible.
- **Benefit**: Keeps RAM usage flat (O(1)) regardless of input size (O(n)).

### 3. Benchmarking (Verification)
Performance is a feature. Verify it with Go Benchmarks.
- **Command**: `go test -bench=. -benchmem`
- **Output**: Shows `ns/op` (speed) and `B/op` (allocations).

## Usage

### Implement Worker Pool
Copy `templates/worker_pool.go`. Modify `Job` struct and `processJob` methods.

### Implement Streaming
Copy `templates/excel_iterator.go` for Excelize v2 usage patterns.

### Verify Performance
Copy `templates/perf_test.go`. Run benchmarks to ensure you are faster than the legacy system.
