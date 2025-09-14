package scheduler

import "sync"

// Mutex is a deterministic mutex.
type Mutex struct {
	mu sync.Mutex
	// TODO: Add deterministic scheduling
}

// NewMutex creates a new deterministic mutex.
func NewMutex() *Mutex {
	return &Mutex{}
}

// Lock locks the mutex.
func (m *Mutex) Lock() {
	m.mu.Lock()
}

// Unlock unlocks the mutex.
func (m *Mutex) Unlock() {
	m.mu.Unlock()
}

// TryLock tries to lock the mutex.
func (m *Mutex) TryLock() bool {
	return m.mu.TryLock()
}

// RWMutex is a deterministic reader/writer mutex.
type RWMutex struct {
	mu sync.RWMutex
	// TODO: Add deterministic scheduling
}

// NewRWMutex creates a new deterministic RWMutex.
func NewRWMutex() *RWMutex {
	return &RWMutex{}
}

// Lock locks for writing.
func (rw *RWMutex) Lock() {
	rw.mu.Lock()
}

// Unlock unlocks for writing.
func (rw *RWMutex) Unlock() {
	rw.mu.Unlock()
}

// RLock locks for reading.
func (rw *RWMutex) RLock() {
	rw.mu.RLock()
}

// RUnlock unlocks for reading.
func (rw *RWMutex) RUnlock() {
	rw.mu.RUnlock()
}