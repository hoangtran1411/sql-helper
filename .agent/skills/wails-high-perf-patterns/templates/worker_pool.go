package patterns

import (
	"context"
	"fmt"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Job represents the unit of work
type Job struct {
	ID   int
	Path string
}

// Result represents the outcome of processing a job
type Result struct {
	JobID int
	Data  interface{}
	Err   error
}

// Processor manages the concurrent worker pool
type Processor struct {
	WorkerCount int
	Jobs        chan Job
	Results     chan Result
	Done        chan bool
}

// NewProcessor creates a pool with specified concurrency
func NewProcessor(workers int) *Processor {
	return &Processor{
		WorkerCount: workers,
		Jobs:        make(chan Job, 100),    // Buffered job queue
		Results:     make(chan Result, 100), // Buffered result queue
		Done:        make(chan bool),
	}
}

// Start initializes workers and returns immediately (non-blocking)
func (p *Processor) Start(ctx context.Context) {
	var wg sync.WaitGroup

	// 1. Spawn Workers
	for i := 0; i < p.WorkerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range p.Jobs {
				select {
				case <-ctx.Done():
					return
				default:
					res := process(job)
					p.Results <- res
				}
			}
		}(i)
	}

	// 2. Waiter Goroutine (Closes results when all workers finish)
	go func() {
		wg.Wait()
		close(p.Results)
		p.Done <- true
	}()
}

// Heavy processing logic
func process(j Job) Result {
	// Simulate work
	return Result{JobID: j.ID, Err: nil}
}

// Example usage in Wails v3 Application
func RunExample(items []string) {
	p := NewProcessor(10)
	ctx := context.Background()
	p.Start(ctx)

	// Producer
	go func() {
		for i, item := range items {
			p.Jobs <- Job{ID: i, Path: item}
		}
		close(p.Jobs)
	}()

	// Collector (Update UI via Wails v3 event emitter)
	total := len(items)
	completed := 0

	for res := range p.Results {
		completed++
		if res.Err != nil {
			fmt.Printf("Job %d failed: %v\n", res.JobID, res.Err)
		}

		// Emit real-time progress to frontend
		if app := application.Get(); app != nil {
			progressPercent := float64(completed) / float64(total) * 100
			app.Event.Emit("progress", progressPercent)
		}
	}
}
