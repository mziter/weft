package examples

import (
	"testing"

	"github.com/mziter/weft"
	"github.com/mziter/weft/wefttest"
)

// TestConcurrentCounter demonstrates basic usage of weft for testing
// a simple concurrent counter.
func TestConcurrentCounter(t *testing.T) {
	wefttest.Explore(t, 100, func(s *weft.Scheduler) {
		var mu weft.Mutex
		counter := 0
		done := make(chan bool, 10)

		// Spawn 10 goroutines that increment the counter
		for i := 0; i < 10; i++ {
			s.Go(func(ctx weft.Context) {
				mu.Lock()
				temp := counter
				ctx.Yield() // Yield to potentially expose race conditions
				counter = temp + 1
				mu.Unlock()
				done <- true
			})
		}

		// Wait for all goroutines
		for i := 0; i < 10; i++ {
			<-done
		}

		// Check the final value
		if counter != 10 {
			t.Errorf("expected counter=10, got %d", counter)
		}
	})
}

// TestDeadlock demonstrates deadlock detection.
func TestDeadlock(t *testing.T) {
	t.Skip("Deadlock detection not yet implemented")
	
	wefttest.Explore(t, 10, func(s *weft.Scheduler) {
		var mu1, mu2 weft.Mutex

		// Classic deadlock scenario
		s.Go(func(ctx weft.Context) {
			mu1.Lock()
			ctx.Yield()
			mu2.Lock()
			mu2.Unlock()
			mu1.Unlock()
		})

		s.Go(func(ctx weft.Context) {
			mu2.Lock()
			ctx.Yield()
			mu1.Lock()
			mu1.Unlock()
			mu2.Unlock()
		})
	})
}

// TestChannel demonstrates deterministic channel operations.
func TestChannel(t *testing.T) {
	wefttest.Explore(t, 50, func(s *weft.Scheduler) {
		ch := weft.MakeChan[int](1)
		result := 0

		// Producer
		s.Go(func(ctx weft.Context) {
			for i := 1; i <= 5; i++ {
				ch.Send(i)
			}
			ch.Close()
		})

		// Consumer
		s.Go(func(ctx weft.Context) {
			for {
				if v, ok := ch.Recv(); ok {
					result += v
				} else {
					break
				}
			}
		})

		s.Wait()

		// Sum of 1+2+3+4+5 = 15
		if result != 15 {
			t.Errorf("expected result=15, got %d", result)
		}
	})
}

// TestConditionVariable demonstrates condition variable usage.
func TestConditionVariable(t *testing.T) {
	wefttest.Explore(t, 50, func(s *weft.Scheduler) {
		var mu weft.Mutex
		cond := weft.NewCond(&mu)
		ready := false
		processed := false

		// Waiter
		s.Go(func(ctx weft.Context) {
			mu.Lock()
			for !ready {
				cond.Wait()
			}
			processed = true
			mu.Unlock()
		})

		// Signaler
		s.Go(func(ctx weft.Context) {
			s.Sleep(10) // Small delay
			mu.Lock()
			ready = true
			cond.Signal()
			mu.Unlock()
		})

		s.Wait()

		if !processed {
			t.Error("condition variable signal was not received")
		}
	})
}

// TestReplay demonstrates how to replay a specific failing seed.
func TestReplay(t *testing.T) {
	// If you found a bug with seed 12345, you can replay it:
	const failingSeed = 12345
	
	wefttest.Replay(t, failingSeed, func(s *weft.Scheduler) {
		// Your test code here
		var counter int
		for i := 0; i < 5; i++ {
			s.Go(func(ctx weft.Context) {
				counter++
			})
		}
		s.Wait()
		
		// This test intentionally has a race condition
		// In deterministic mode, it will consistently fail or pass
		// depending on the seed
	})
}