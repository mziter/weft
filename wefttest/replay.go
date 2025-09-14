package wefttest

import (
	"testing"

	"github.com/mziter/weft"
)

// Replay runs the build function with a specific seed for reproduction.
func Replay(t *testing.T, seed uint64, build BuildFunc) {
	t.Helper()
	
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