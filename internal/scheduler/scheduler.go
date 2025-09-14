package scheduler

import (
	"math/rand"
	"sync"
	"time"
)

// Scheduler manages deterministic task execution.
type Scheduler struct {
	mu       sync.Mutex
	rng      *rand.Rand
	tasks    []*Task
	runnable []int
	current  int
	waitGroup sync.WaitGroup
}

// New creates a new scheduler with the given seed.
func New(seed uint64) *Scheduler {
	return &Scheduler{
		rng: rand.New(rand.NewSource(int64(seed))),
	}
}

// Spawn creates a new task.
func (s *Scheduler) Spawn(fn func(interface{})) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// TODO: Implement task spawning
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		fn(nil)
	}()
}

// Wait waits for all tasks to complete.
func (s *Scheduler) Wait() {
	s.waitGroup.Wait()
}

// Sleep pauses the current task.
func (s *Scheduler) Sleep(d time.Duration) {
	// TODO: Implement virtual time sleep
	time.Sleep(d / 1000) // Speed up for testing
}

// After returns a timer channel.
func (s *Scheduler) After(d time.Duration) <-chan time.Time {
	// TODO: Implement virtual time after
	return time.After(d / 1000) // Speed up for testing
}