package wefttest

import (
	"testing"

	"github.com/mziter/weft"
)

// Replay runs the build function with a specific seed for reproduction.
func Replay(t testing.TB, seed uint64, build BuildFunc) {
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

	s := weft.NewScheduler(seed)
	
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic during replay with seed %d: %v", seed, r)
		}
	}()
	
	build(s)
	s.Wait()
}

// ReplayChoices runs the build function with an explicit choice sequence.
// This is useful for replaying a minimal trace after shrinking.
func ReplayChoices(t *testing.T, choices []int, build BuildFunc) {
	t.Helper()
	
	// TODO: Implement choice-based replay when scheduler supports it
	// For now, this is a placeholder
	t.Skip("ReplayChoices not yet implemented")
}