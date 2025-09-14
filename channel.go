//go:build detsched

package weft

import (
	"github.com/yourusername/weft/internal/scheduler"
)

// Chan is a deterministic channel.
type Chan[T any] struct {
	ch *scheduler.Chan[T]
}

// MakeChan creates a new deterministic channel with the given capacity.
func MakeChan[T any](cap int) Chan[T] {
	return Chan[T]{
		ch: scheduler.MakeChan[T](cap),
	}
}

// Send sends a value on the channel.
func (c Chan[T]) Send(v T) {
	c.ch.Send(v)
}

// Recv receives a value from the channel.
func (c Chan[T]) Recv() (T, bool) {
	return c.ch.Recv()
}

// TrySend attempts to send without blocking.
func (c Chan[T]) TrySend(v T) bool {
	return c.ch.TrySend(v)
}

// TryRecv attempts to receive without blocking.
func (c Chan[T]) TryRecv() (T, bool) {
	return c.ch.TryRecv()
}

// Close closes the channel.
func (c Chan[T]) Close() {
	c.ch.Close()
}