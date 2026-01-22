package patterns

import (
	"testing"
)

// To run this: go test -bench=. -benchmem

// BenchmarkSingleThread benchmarks the slow synchronous approach
func BenchmarkSingleThread(b *testing.B) {
	// Setup generic data
	jobs := make([]Job, 1000)
	for i := range jobs {
		jobs[i] = Job{ID: i, Path: "dummy"}
	}

	b.ResetTimer() // Start tracking from here
	for i := 0; i < b.N; i++ {
		// Run entire sync process
		for _, j := range jobs {
			process(j)
		}
	}
}

// BenchmarkWorkerPool benchmarks the concurrent approach
func BenchmarkWorkerPool(b *testing.B) {
	jobs := make([]string, 1000)
	for i := range jobs {
		jobs[i] = "dummy"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Run async process
		RunExample(jobs)
	}
}

// Ensure zero allocations for critical paths
func BenchmarkCriticalFunction(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// Call your hot path function
		// optimizedFunction()
	}
}
