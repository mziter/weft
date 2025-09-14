//go:build detsched

package weft

import (
	"time"

	"github.com/yourusername/weft/internal/scheduler"
)

// Scheduler controls the execution of deterministic tasks.
type Scheduler struct {
	sched *scheduler.Scheduler
}

// NewScheduler creates a new deterministic scheduler with the given seed.
func NewScheduler(seed uint64) *Scheduler {
	return &Scheduler{
		sched: scheduler.New(seed),
	}
}

// Go spawns a new deterministic goroutine.
func Go(fn func(Context)) {
	defaultScheduler.Go(fn)
}

// Go spawns a new deterministic goroutine on this scheduler.
func (s *Scheduler) Go(fn func(Context)) {
	s.sched.Spawn(fn)
}

// Wait blocks until all spawned tasks complete.
func (s *Scheduler) Wait() {
	s.sched.Wait()
}

// Sleep pauses the current task for the specified duration.
func Sleep(d time.Duration) {
	defaultScheduler.Sleep(d)
}

// Sleep pauses the current task for the specified duration.
func (s *Scheduler) Sleep(d time.Duration) {
	s.sched.Sleep(d)
}

// After returns a channel that receives after the duration.
func After(d time.Duration) <-chan time.Time {
	return defaultScheduler.After(d)
}

// After returns a channel that receives after the duration.
func (s *Scheduler) After(d time.Duration) <-chan time.Time {
	return s.sched.After(d)
}

var defaultScheduler = NewScheduler(0)