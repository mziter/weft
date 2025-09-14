package weft

import "time"

// Clock provides virtual time for deterministic testing.
type Clock interface {
	// Now returns the current virtual time.
	Now() time.Time

	// Advance advances the virtual time by the given duration.
	Advance(d time.Duration)
}