package examples

import (
	"testing"

	"github.com/mziter/weft"
	"github.com/mziter/weft/wefttest"
)

// TestReplaySpecificSeed demonstrates reproducing a specific bug.
func TestReplaySpecificSeed(t *testing.T) {
	// If you found a bug with a specific seed, you can replay it
	const bugSeed = 42

	wefttest.Replay(t, bugSeed, func(s *weft.Scheduler) {
		counter := NewCounter()

		// This will have the same deterministic behavior every time
		for i := 0; i < 3; i++ {
			s.Go(func(ctx weft.Context) {
				counter.IncrementWithWork()
			})
		}

		s.Wait()

		// With seed 42, this might consistently fail or pass
		// allowing you to debug the specific interleaving
		t.Logf("With seed %d, counter value is %d", bugSeed, counter.Value())
	})
}