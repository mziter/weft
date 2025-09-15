package wefttest

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/mziter/weft"
)

// BuildFunc is a function that builds a test scenario using a scheduler.
type BuildFunc func(*weft.Scheduler)

// Explore runs the build function with multiple different schedules.
func Explore(t testing.TB, runs int, build BuildFunc) {
	t.Helper()

	if !isDeterministicModeAvailable() {
		t.Skipf(`
Deterministic concurrency testing not available.
For comprehensive concurrency testing that can detect race conditions,
deadlocks, and other subtle bugs, run with:

    go test -tags=detsched

This enables Weft's deterministic scheduler which explores multiple
execution orders to find bugs that standard tests might miss.`)
		return
	}

	rng := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))

	for i := 0; i < runs; i++ {
		seed := rng.Uint64()
		// Type assert to *testing.T for Run method
		if tt, ok := t.(*testing.T); ok {
			tt.Run(fmt.Sprintf("seed_%d", seed), func(t *testing.T) {
				t.Helper()
				s := weft.NewScheduler(seed)

				// Run the build function
				defer func() {
					if r := recover(); r != nil {
						t.Fatalf("panic with seed %d: %v", seed, r)
					}
				}()

				build(s)
				s.Wait()
			})
		} else {
			// Fallback for non-*testing.T types (like our mock)
			s := weft.NewScheduler(seed)

			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("panic with seed %d: %v", seed, r)
				}
			}()

			build(s)
			s.Wait()
		}
	}
}

// ExploreWithSeeds runs the build function with specific seeds.
func ExploreWithSeeds(t testing.TB, seeds []uint64, build BuildFunc) {
	t.Helper()

	if !isDeterministicModeAvailable() {
		t.Skipf(`
Deterministic concurrency testing not available.
For comprehensive concurrency testing that can detect race conditions,
deadlocks, and other subtle bugs, run with:

    go test -tags=detsched

This enables Weft's deterministic scheduler which explores multiple
execution orders to find bugs that standard tests might miss.`)
		return
	}

	for _, seed := range seeds {
		// Type assert to *testing.T for Run method
		if tt, ok := t.(*testing.T); ok {
			tt.Run(fmt.Sprintf("seed_%d", seed), func(t *testing.T) {
				t.Helper()
				s := weft.NewScheduler(seed)

				defer func() {
					if r := recover(); r != nil {
						t.Fatalf("panic with seed %d: %v", seed, r)
					}
				}()

				build(s)
				s.Wait()
			})
		} else {
			// Fallback for non-*testing.T types (like our mock)
			s := weft.NewScheduler(seed)

			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("panic with seed %d: %v", seed, r)
				}
			}()

			build(s)
			s.Wait()
		}
	}
}