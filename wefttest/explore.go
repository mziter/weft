package wefttest

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/yourusername/weft"
)

// BuildFunc is a function that builds a test scenario using a scheduler.
type BuildFunc func(*weft.Scheduler)

// Explore runs the build function with multiple different schedules.
func Explore(t *testing.T, runs int, build BuildFunc) {
	t.Helper()

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	for i := 0; i < runs; i++ {
		seed := rng.Uint64()
		t.Run(fmt.Sprintf("seed_%d", seed), func(t *testing.T) {
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
	}
}

// ExploreWithSeeds runs the build function with specific seeds.
func ExploreWithSeeds(t *testing.T, seeds []uint64, build BuildFunc) {
	t.Helper()

	for _, seed := range seeds {
		t.Run(fmt.Sprintf("seed_%d", seed), func(t *testing.T) {
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
	}
}