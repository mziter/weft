package examples

import (
	"github.com/mziter/weft"
)

// WorkerPool manages a pool of workers processing jobs concurrently.
// This demonstrates condition variables and worker coordination patterns.
type WorkerPool struct {
	mu          weft.Mutex
	cond        *weft.Cond
	jobs        []func()
	shutdown    bool
	activeJobs  int
	workerCount int
}

// NewWorkerPool creates a new worker pool with the specified number of workers.
func NewWorkerPool(workerCount int) *WorkerPool {
	wp := &WorkerPool{
		jobs:        make([]func(), 0),
		workerCount: workerCount,
	}
	wp.cond = weft.NewCond(&wp.mu)
	return wp
}

// Start begins the worker goroutines.
func (wp *WorkerPool) Start(s *weft.Scheduler) {
	for i := 0; i < wp.workerCount; i++ {
		s.Go(func(ctx weft.Context) {
			wp.worker(ctx)
		})
	}
}

// Submit adds a job to the pool.
func (wp *WorkerPool) Submit(job func()) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	wp.jobs = append(wp.jobs, job)
	wp.cond.Signal() // Wake up one worker
}

// Shutdown stops all workers after current jobs complete.
func (wp *WorkerPool) Shutdown() {
	wp.mu.Lock()
	wp.shutdown = true
	wp.cond.Broadcast() // Wake all workers
	wp.mu.Unlock()
}

// Wait blocks until all jobs are completed.
func (wp *WorkerPool) Wait() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	for wp.activeJobs > 0 || len(wp.jobs) > 0 {
		wp.cond.Wait()
	}
}

// worker is the main loop for each worker goroutine.
func (wp *WorkerPool) worker(ctx weft.Context) {
	for {
		wp.mu.Lock()

		// Wait for work or shutdown signal
		for len(wp.jobs) == 0 && !wp.shutdown {
			wp.cond.Wait()
		}

		// Check for shutdown
		if wp.shutdown && len(wp.jobs) == 0 {
			wp.mu.Unlock()
			return
		}

		// Get next job
		if len(wp.jobs) == 0 {
			wp.mu.Unlock()
			continue
		}

		job := wp.jobs[0]
		wp.jobs = wp.jobs[1:]
		wp.activeJobs++
		wp.mu.Unlock()

		// Execute job (outside of lock)
		job()

		// Mark job complete
		wp.mu.Lock()
		wp.activeJobs--
		wp.cond.Broadcast() // Signal waiters
		wp.mu.Unlock()
	}
}