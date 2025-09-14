package scheduler

import "sync"

// Cond is a deterministic condition variable.
type Cond struct {
	cond *sync.Cond
	// TODO: Add deterministic scheduling
}

// NewCond creates a new deterministic condition variable.
func NewCond(l interface{}) *Cond {
	if locker, ok := l.(sync.Locker); ok {
		return &Cond{
			cond: sync.NewCond(locker),
		}
	}
	// For weft types, we need to extract the underlying sync.Locker
	// This is a stub implementation
	return &Cond{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

// Wait waits for the condition.
func (c *Cond) Wait() {
	c.cond.Wait()
}

// Signal wakes one waiter.
func (c *Cond) Signal() {
	c.cond.Signal()
}

// Broadcast wakes all waiters.
func (c *Cond) Broadcast() {
	c.cond.Broadcast()
}