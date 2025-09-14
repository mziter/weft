//go:build !detsched

package weft

import "sync"

// Cond is a standard sync.Cond in production mode.
type Cond struct {
	*sync.Cond
}

// NewCond returns a new Cond with the given Locker.
func NewCond(l Locker) *Cond {
	return &Cond{
		Cond: sync.NewCond(l),
	}
}

// Locker represents types that can be locked and unlocked.
type Locker interface {
	Lock()
	Unlock()
}