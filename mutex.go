//go:build detsched

package weft

import (
	"github.com/yourusername/weft/internal/scheduler"
)

// Mutex is a deterministic mutual exclusion lock.
type Mutex struct {
	mu *scheduler.Mutex
}

// Lock locks the mutex.
func (m *Mutex) Lock() {
	if m.mu == nil {
		m.mu = scheduler.NewMutex()
	}
	m.mu.Lock()
}

// Unlock unlocks the mutex.
func (m *Mutex) Unlock() {
	if m.mu == nil {
		panic("unlock of unlocked mutex")
	}
	m.mu.Unlock()
}

// TryLock tries to lock the mutex and returns true if successful.
func (m *Mutex) TryLock() bool {
	if m.mu == nil {
		m.mu = scheduler.NewMutex()
	}
	return m.mu.TryLock()
}

// RWMutex is a deterministic reader/writer mutual exclusion lock.
type RWMutex struct {
	mu *scheduler.RWMutex
}

// Lock locks the mutex for writing.
func (rw *RWMutex) Lock() {
	if rw.mu == nil {
		rw.mu = scheduler.NewRWMutex()
	}
	rw.mu.Lock()
}

// Unlock unlocks the mutex for writing.
func (rw *RWMutex) Unlock() {
	if rw.mu == nil {
		panic("unlock of unlocked mutex")
	}
	rw.mu.Unlock()
}

// RLock locks the mutex for reading.
func (rw *RWMutex) RLock() {
	if rw.mu == nil {
		rw.mu = scheduler.NewRWMutex()
	}
	rw.mu.RLock()
}

// RUnlock unlocks the mutex for reading.
func (rw *RWMutex) RUnlock() {
	if rw.mu == nil {
		panic("runlock of unlocked mutex")
	}
	rw.mu.RUnlock()
}