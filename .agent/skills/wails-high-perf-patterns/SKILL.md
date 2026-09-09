---
name: Wails High Performance Patterns
description: Reusable patterns for building high-performance Go/Wails apps, focusing on Concurrency, Memory Efficiency, and Benchmarking.
---

# Wails High Performance Patterns

This skill encapsulates patterns for high-performance Wails applications, specifically designed for data-intensive tasks like batch Excel processing and file generation.

## When to Use
- **Heavy Processing**: Handling tens of thousands of rows or large Excel sheets.
- **Responsive UI**: Ensuring the frontend remains smooth while the Go backend crunches numbers.
- **Low Memory Footprint**: Avoiding Out-Of-Memory (OOM) errors by using streaming iterators and buffered disk I/O.

## Core Patterns

### 1. Worker Pool (Concurrency)
Instead of spawning execution logic linearly or spawning unlimited goroutines (which causes thrashing), use a Worker Pool.
- **Concept**: Fixed number of workers consuming from a buffered job channel.
- **Benefit**: Controlled resource usage, maximum CPU throughput.
- **Wails v3 UI Integration**: Send progress events from the "Collector" phase via `application.Get().Event.Emit("progress", data)`, not individual workers, to avoid flooding the frontend event loop.

### 2. Stream Processing (Memory)
Avoid loading entire files into memory.
- **Excelize Streaming**: Use `rows.Next()` iterator (`excel.IterateSheet`) instead of loading all cells with `GetRows()`.
- **Buffered Output**: Write SQL batches directly to disk using `bufio.NewWriter`.
- **Benefit**: Keeps RAM usage flat (O(1)) regardless of input size (O(n)).

### 3. Benchmarking (Verification)
Performance is a feature. Verify it with Go Benchmarks.
- **Command**: `go test -bench=. -benchmem ./...`
- **Output**: Reports `ns/op` (speed) and `B/op` (allocations).

## Usage

### Implement Worker Pool
Use `templates/worker_pool.go`. Modify `Job` struct and `process` method for your specific pipeline.

### Implement Streaming
Use `templates/excel_iterator.go` for Excelize v2 row streaming and buffered file export patterns.

### Verify Performance
Use `templates/perf_test.go` to measure and benchmark critical execution paths.
