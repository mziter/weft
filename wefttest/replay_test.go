package wefttest

import (
	"strings"
	"testing"

	"github.com/mziter/weft"
)

// TestReplaySkipsWithoutDetschedTag verifies that Replay skips gracefully
// when run without the detsched build tag.
func TestReplaySkipsWithoutDetschedTag(t *testing.T) {
	mockT := newMockTestingT(t)

	// Try to run Replay with a specific seed
	Replay(mockT, 12345, func(s *weft.Scheduler) {
		t.Error("Build function should not run without detsched tag")
	})

	// Check behavior based on build tags
	if isDeterministicModeAvailable() {
		// With -tags=detsched, test should run (or panic since scheduler isn't implemented)
		if mockT.skipped {
			t.Error("Replay should not skip when deterministic mode is available")
		}
	} else {
		// Without -tags=detsched, test should skip
		if !mockT.skipped {
			t.Error("Replay should skip when deterministic mode is not available")
		}

		// Verify the skip message matches Explore's message
		if !strings.Contains(mockT.skipMessage, "go test -tags=detsched") {
			t.Errorf("Skip message should contain usage instructions, got: %s", mockT.skipMessage)
		}

		if !strings.Contains(mockT.skipMessage, "Deterministic concurrency testing not available") {
			t.Errorf("Skip message should explain why test was skipped, got: %s", mockT.skipMessage)
		}

		if !strings.Contains(mockT.skipMessage, "explores multiple") {
			t.Errorf("Skip message should explain benefits, got: %s", mockT.skipMessage)
		}
	}
}

// TestReplayMessageConsistency verifies that Replay and Explore use the same skip message.
func TestReplayMessageConsistency(t *testing.T) {
	if isDeterministicModeAvailable() {
		t.Skip("This test only runs without -tags=detsched")
	}

	// Get skip message from Explore
	exploreT := newMockTestingT(t)
	Explore(exploreT, 1, func(s *weft.Scheduler) {})

	// Get skip message from Replay
	replayT := newMockTestingT(t)
	Replay(replayT, 123, func(s *weft.Scheduler) {})

	// Messages should be identical
	if exploreT.skipMessage != replayT.skipMessage {
		t.Errorf("Skip messages should be identical.\nExplore: %s\nReplay: %s",
			exploreT.skipMessage, replayT.skipMessage)
	}
}