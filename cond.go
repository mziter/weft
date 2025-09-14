//go:build detsched

package weft

import (
	"github.com/yourusername/weft/internal/scheduler"
)

// Cond implements a condition variable for deterministic testing.
type Cond struct {
	cond *scheduler.Cond
}

// NewCond returns a new Cond with the given Locker.
func NewCond(l Locker) *Cond {
	return &Cond{
		cond: scheduler.NewCond(l),
	}
}

// Wait atomically unlocks the Locker and waits to be signaled.
func (c *Cond) Wait() {
	c.cond.Wait()
}

// Signal wakes one goroutine waiting on the condition variable.
func (c *Cond) Signal() {
	c.cond.Signal()
}

// Broadcast wakes all goroutines waiting on the condition variable.
func (c *Cond) Broadcast() {
	c.cond.Broadcast()
}

// Locker represents types that can be locked and unlocked.
type Locker interface {
	Lock()
	Unlock()
}