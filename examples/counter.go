package examples

import (
	"github.com/mziter/weft"
)

// Counter is a thread-safe counter implementation.
// This demonstrates a typical production component that needs concurrency testing.
type Counter struct {
	mu    weft.Mutex
	value int
}

// NewCounter creates a new counter starting at zero.
func NewCounter() *Counter {
	return &Counter{}
}

// Increment adds one to the counter.
func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

// IncrementWithWork simulates doing some work between read and write.
// This has a race condition - it's not properly synchronized.
// In real code, this might represent complex business logic between operations.
func (c *Counter) IncrementWithWork() {
	// Read the value (with proper locking)
	c.mu.Lock()
	temp := c.value
	c.mu.Unlock()

	// Simulate some work that takes time
	// In real code this might be validation, logging, external calls, etc.
	// The time gap creates opportunity for race conditions
	_ = temp * 2 // Some computation

	// Write back the incremented value (race condition here!)
	c.mu.Lock()
	c.value = temp + 1
	c.mu.Unlock()
}

// Value returns the current counter value.
func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// Reset sets the counter back to zero.
func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value = 0
}