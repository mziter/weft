package wefttest

import (
	"strings"
	"testing"

	"github.com/mziter/weft"
)

// TestExploreSkipsWithoutDetschedTag verifies that Explore skips gracefully
// when run without the detsched build tag.
func TestExploreSkipsWithoutDetschedTag(t *testing.T) {
	// Create a mock test to capture skip behavior
	mockT := newMockTestingT(t)

	// Try to run Explore
	Explore(mockT, 10, func(s *weft.Scheduler) {
		// This should not be executed if detection works
		t.Error("Build function should not run without detsched tag")
	})

	// Check behavior based on build tags
	if isDeterministicModeAvailable() {
		// With -tags=detsched, test should run (or panic since scheduler isn't implemented)
		if mockT.skipped {
			t.Error("Explore should not skip when deterministic mode is available")
		}
	} else {
		// Without -tags=detsched, test should skip
		if !mockT.skipped {
			t.Error("Explore should skip when deterministic mode is not available")
		}

		// Verify the skip message is helpful
		if !strings.Contains(mockT.skipMessage, "go test -tags=detsched") {
			t.Errorf("Skip message should contain usage instructions, got: %s", mockT.skipMessage)
		}

		if !strings.Contains(mockT.skipMessage, "Deterministic concurrency testing not available") {
			t.Errorf("Skip message should explain why test was skipped, got: %s", mockT.skipMessage)
		}
	}
}

// TestExploreWithSeedsSkipsWithoutDetschedTag verifies that ExploreWithSeeds
// skips gracefully when run without the detsched build tag.
func TestExploreWithSeedsSkipsWithoutDetschedTag(t *testing.T) {
	mockT := newMockTestingT(t)

	seeds := []uint64{123, 456, 789}
	ExploreWithSeeds(mockT, seeds, func(s *weft.Scheduler) {
		t.Error("Build function should not run without detsched tag")
	})

	if isDeterministicModeAvailable() {
		if mockT.skipped {
			t.Error("ExploreWithSeeds should not skip when deterministic mode is available")
		}
	} else {
		if !mockT.skipped {
			t.Error("ExploreWithSeeds should skip when deterministic mode is not available")
		}

		// Verify consistent skip message
		if !strings.Contains(mockT.skipMessage, "go test -tags=detsched") {
			t.Errorf("Skip message should contain usage instructions, got: %s", mockT.skipMessage)
		}
	}
}

// mockTestingT is a mock implementation of testing.T for testing skip behavior.
type mockTestingT struct {
	*testing.T  // Embed to satisfy interface
	skipped     bool
	skipMessage string
	failed      bool
	failMessage string
	logs        []string
	subtests    []string
}

func newMockTestingT(t *testing.T) *mockTestingT {
	return &mockTestingT{T: t}
}

func (m *mockTestingT) Helper() {}

func (m *mockTestingT) Skipf(format string, args ...interface{}) {
	m.skipped = true
	m.skipMessage = strings.TrimSpace(format) // Remove leading/trailing whitespace
}

func (m *mockTestingT) Skip(args ...interface{}) {
	m.skipped = true
	if len(args) > 0 {
		m.skipMessage = args[0].(string)
	}
}

func (m *mockTestingT) Fatalf(format string, args ...interface{}) {
	m.failed = true
	m.failMessage = format
	// Don't actually fail the test
}

func (m *mockTestingT) Errorf(format string, args ...interface{}) {
	m.failed = true
	m.failMessage = format
}

func (m *mockTestingT) Run(name string, f func(*testing.T)) bool {
	m.subtests = append(m.subtests, name)
	// For mock purposes, we don't actually run the subtest
	// In real implementation with -tags=detsched, this would run
	return !m.skipped
}

func (m *mockTestingT) Log(args ...interface{}) {
	m.logs = append(m.logs, args[0].(string))
}

func (m *mockTestingT) Logf(format string, args ...interface{}) {
	m.logs = append(m.logs, format)
}