//go:build !detsched

package weft

// Chan is a regular Go channel in production mode.
type Chan[T any] struct {
	ch chan T
}

// MakeChan creates a new channel with the given capacity.
func MakeChan[T any](cap int) Chan[T] {
	return Chan[T]{
		ch: make(chan T, cap),
	}
}

// Send sends a value on the channel.
func (c Chan[T]) Send(v T) {
	c.ch <- v
}

// Recv receives a value from the channel.
func (c Chan[T]) Recv() (T, bool) {
	v, ok := <-c.ch
	return v, ok
}

// TrySend attempts to send without blocking.
func (c Chan[T]) TrySend(v T) bool {
	select {
	case c.ch <- v:
		return true
	default:
		return false
	}
}

// TryRecv attempts to receive without blocking.
func (c Chan[T]) TryRecv() (T, bool) {
	select {
	case v, ok := <-c.ch:
		return v, ok
	default:
		var zero T
		return zero, false
	}
}

// Close closes the channel.
func (c Chan[T]) Close() {
	close(c.ch)
}