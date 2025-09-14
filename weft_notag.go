//go:build !detsched

package weft

import (
	"time"
)

// Scheduler is a no-op in production mode.
type Scheduler struct{}

// NewScheduler returns a no-op scheduler in production mode.
func NewScheduler(seed uint64) *Scheduler {
	return &Scheduler{}
}

// Go spawns a regular goroutine in production mode.
func Go(fn func(Context)) {
	go fn(productionContext{})
}

// Go spawns a regular goroutine in production mode.
func (s *Scheduler) Go(fn func(Context)) {
	go fn(productionContext{})
}

// Wait is a no-op in production mode.
func (s *Scheduler) Wait() {
	// In production mode, there's no tracking of goroutines
}

// Sleep delegates to time.Sleep in production mode.
func Sleep(d time.Duration) {
	time.Sleep(d)
}

// Sleep delegates to time.Sleep in production mode.
func (s *Scheduler) Sleep(d time.Duration) {
	time.Sleep(d)
}

// After delegates to time.After in production mode.
func After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

// After delegates to time.After in production mode.
func (s *Scheduler) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

type productionContext struct{}

func (productionContext) Yield() {}
func (productionContext) Done() <-chan struct{} { return nil }