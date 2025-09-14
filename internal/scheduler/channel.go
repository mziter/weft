package scheduler

// Chan is a deterministic channel.
type Chan[T any] struct {
	ch   chan T
	// TODO: Add deterministic scheduling
}

// MakeChan creates a new deterministic channel.
func MakeChan[T any](cap int) *Chan[T] {
	return &Chan[T]{
		ch: make(chan T, cap),
	}
}

// Send sends a value.
func (c *Chan[T]) Send(v T) {
	c.ch <- v
}

// Recv receives a value.
func (c *Chan[T]) Recv() (T, bool) {
	v, ok := <-c.ch
	return v, ok
}

// TrySend tries to send without blocking.
func (c *Chan[T]) TrySend(v T) bool {
	select {
	case c.ch <- v:
		return true
	default:
		return false
	}
}

// TryRecv tries to receive without blocking.
func (c *Chan[T]) TryRecv() (T, bool) {
	select {
	case v, ok := <-c.ch:
		return v, ok
	default:
		var zero T
		return zero, false
	}
}

// Close closes the channel.
func (c *Chan[T]) Close() {
	close(c.ch)
}