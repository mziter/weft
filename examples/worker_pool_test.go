package examples

import (
	"testing"

	"github.com/mziter/weft"
	"github.com/mziter/weft/wefttest"
)

// TestWorkerPool demonstrates testing condition variable coordination.
func TestWorkerPool(t *testing.T) {
	wefttest.Explore(t, 30, func(s *weft.Scheduler) {
		wp := NewWorkerPool(3)
		wp.Start(s)

		// Track completed jobs
		var completed []int
		var mu weft.Mutex

		// Submit jobs
		for i := 1; i <= 10; i++ {
			jobID := i
			wp.Submit(func() {
				mu.Lock()
				completed = append(completed, jobID)
				mu.Unlock()
			})
		}

		wp.Wait()
		wp.Shutdown()
		s.Wait()

		// Verify all jobs completed
		if len(completed) != 10 {
			t.Errorf("expected 10 jobs completed, got %d", len(completed))
		}

		// Verify all job IDs are present
		seen := make(map[int]bool)
		for _, id := range completed {
			seen[id] = true
		}
		for i := 1; i <= 10; i++ {
			if !seen[i] {
				t.Errorf("job %d was not completed", i)
			}
		}
	})
}