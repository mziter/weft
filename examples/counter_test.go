package examples

import (
	"testing"

	"github.com/mziter/weft"
	"github.com/mziter/weft/wefttest"
)

// TestCounter demonstrates testing a thread-safe counter component.
// This shows how to test production code that uses weft primitives.
func TestCounter(t *testing.T) {
	wefttest.Explore(t, 100, func(s *weft.Scheduler) {
		counter := NewCounter()

		// Spawn multiple goroutines that increment the counter
		numWorkers := 10
		for i := 0; i < numWorkers; i++ {
			s.Go(func(ctx weft.Context) {
				counter.Increment()
			})
		}

		s.Wait()

		// Verify the final value
		if counter.Value() != numWorkers {
			t.Errorf("expected counter=%d, got %d", numWorkers, counter.Value())
		}
	})
}

// TestCounterRaceCondition demonstrates how weft can expose race conditions
// in production code that appears thread-safe but has subtle bugs.
func TestCounterRaceCondition(t *testing.T) {
	wefttest.Explore(t, 200, func(s *weft.Scheduler) {
		counter := NewCounter()

		// Use the method with a race condition
		numWorkers := 5
		for i := 0; i < numWorkers; i++ {
			s.Go(func(ctx weft.Context) {
				counter.IncrementWithWork()
			})
		}

		s.Wait()

		// This should fail in many schedules due to the race condition
		// The deterministic scheduler will expose the timing-dependent bug
		final := counter.Value()
		if final != numWorkers {
			t.Logf("Race condition detected: expected %d, got %d", numWorkers, final)
			// In real tests, this would be t.Errorf() after you fix the bug
		}
	})
}