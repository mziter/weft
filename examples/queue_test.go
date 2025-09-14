package examples

import (
	"testing"

	"github.com/mziter/weft"
	"github.com/mziter/weft/wefttest"
)

// TestProducerConsumer demonstrates testing channel-based synchronization.
func TestProducerConsumer(t *testing.T) {
	wefttest.Explore(t, 50, func(s *weft.Scheduler) {
		pc := NewProducerConsumer(5)
		items := []int{1, 2, 3, 4, 5}

		// Producer goroutine
		s.Go(func(ctx weft.Context) {
			pc.Produce(items)
			pc.queue.Close()
		})

		// Consumer goroutine
		var consumed []int
		s.Go(func(ctx weft.Context) {
			consumed = pc.Consume()
		})

		s.Wait()

		// Verify all items were consumed (order may vary)
		if len(consumed) != len(items) {
			t.Errorf("expected %d items consumed, got %d", len(items), len(consumed))
		}

		// Verify sum is correct (independent of order)
		expectedSum, actualSum := 0, 0
		for _, v := range items {
			expectedSum += v
		}
		for _, v := range consumed {
			actualSum += v
		}
		if actualSum != expectedSum {
			t.Errorf("expected sum=%d, got sum=%d", expectedSum, actualSum)
		}
	})
}