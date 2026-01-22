package patterns

import (
	"context"
	"fmt"
	"sync"
)

// Job represents the unit of work
type Job struct {
	ID   int
	Path string
	// Add other fields needed for processing
}

// Result represents the outcome
type Result struct {
	JobID int
	Data  interface{}
	Err   error
}

// Processor manages the pool
type Processor struct {
	WorkerCount int
	Jobs        chan Job
	Results     chan Result
	Done        chan bool
}

// NewProcessor creates a pool
func NewProcessor(workers int) *Processor {
	return &Processor{
		WorkerCount: workers,
		Jobs:        make(chan Job, 100),    // Buffer job queue
		Results:     make(chan Result, 100), // Buffer result queue
		Done:        make(chan bool),
	}
}

// Start initializes workers and returns immediately (non-blocking)
func (p *Processor) Start(ctx context.Context) {
	var wg sync.WaitGroup

	// 1. Spawn Workers
	for i := 0; i < p.WorkerCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range p.Jobs {
				// Check context cancellation
				select {
				case <-ctx.Done():
					return
				default:
					// Do Heavy Work
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
	// img, err := imaging.Open(j.Path) ...
	return Result{JobID: j.ID, Err: nil}
}

// Example usage in Wails App
func RunExample(items []string) {
	p := NewProcessor(10)
	ctx := context.Background()
	p.Start(ctx)

	// Sender
	go func() {
		for i, item := range items {
			p.Jobs <- Job{ID: i, Path: item}
		}
		close(p.Jobs)
	}()

	// Receiver (Update UI here)
	for res := range p.Results {
		if res.Err != nil {
			fmt.Println("Error:", res.Err)
		} else {
			// runtime.EventsEmit(ctx, "progress", ...)
		}
	}
}
