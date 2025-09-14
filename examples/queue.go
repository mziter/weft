package examples

import (
	"github.com/mziter/weft"
)

// Queue is a thread-safe FIFO queue implementation using channels.
// This demonstrates testing channel-based synchronization.
type Queue[T any] struct {
	items weft.Chan[T]
	done  weft.Chan[struct{}]
}

// NewQueue creates a new queue with the given capacity.
func NewQueue[T any](capacity int) *Queue[T] {
	return &Queue[T]{
		items: weft.MakeChan[T](capacity),
		done:  weft.MakeChan[struct{}](0),
	}
}

// Push adds an item to the queue.
func (q *Queue[T]) Push(item T) bool {
	return q.items.TrySend(item)
}

// Pop removes and returns an item from the queue.
func (q *Queue[T]) Pop() (T, bool) {
	return q.items.TryRecv()
}

// Close signals that no more items will be added.
func (q *Queue[T]) Close() {
	q.items.Close()
	q.done.Close()
}

// ProducerConsumer demonstrates a typical producer-consumer pattern.
type ProducerConsumer struct {
	queue    *Queue[int]
	mu       weft.Mutex
	consumed []int
}

// NewProducerConsumer creates a new producer-consumer system.
func NewProducerConsumer(queueSize int) *ProducerConsumer {
	return &ProducerConsumer{
		queue:    NewQueue[int](queueSize),
		consumed: make([]int, 0),
	}
}

// Produce adds items to the queue.
func (pc *ProducerConsumer) Produce(items []int) {
	for _, item := range items {
		for !pc.queue.Push(item) {
			// In real code, might use backoff or context cancellation
		}
	}
}

// Consume processes items from the queue.
func (pc *ProducerConsumer) Consume() []int {
	for {
		item, ok := pc.queue.Pop()
		if !ok {
			break
		}
		pc.mu.Lock()
		pc.consumed = append(pc.consumed, item)
		pc.mu.Unlock()
	}

	pc.mu.Lock()
	result := pc.consumed
	pc.mu.Unlock()
	return result
}

// GetConsumed returns all consumed items.
func (pc *ProducerConsumer) GetConsumed() []int {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	return append([]int{}, pc.consumed...)
}