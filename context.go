package weft

// Context provides control over a deterministic task.
type Context interface {
	// Yield voluntarily yields control to the scheduler.
	Yield()

	// Done returns a channel that's closed when the context is cancelled.
	Done() <-chan struct{}
}