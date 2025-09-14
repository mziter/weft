//go:build !detsched

package weft

import "sync"

// Mutex is a standard sync.Mutex in production mode.
type Mutex struct {
	sync.Mutex
}

// TryLock tries to lock the mutex and returns true if successful.
func (m *Mutex) TryLock() bool {
	return m.Mutex.TryLock()
}

// RWMutex is a standard sync.RWMutex in production mode.
type RWMutex struct {
	sync.RWMutex
}