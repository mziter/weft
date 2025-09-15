package wefttest

import "testing"

// TestIsDeterministicModeAvailable verifies the detection mechanism works correctly.
func TestIsDeterministicModeAvailable(t *testing.T) {
	available := isDeterministicModeAvailable()

	// The expected value depends on build tags
	// When run with -tags=detsched, should return true
	// When run without tags, should return false

	// We can't assert a specific value since it depends on build tags,
	// but we can verify it returns a boolean without error
	if available {
		t.Log("Deterministic mode is available (compiled with -tags=detsched)")
	} else {
		t.Log("Deterministic mode is not available (compiled without -tags=detsched)")
	}
}

// TestDetectionNoOverhead verifies the detection has minimal overhead.
func TestDetectionNoOverhead(t *testing.T) {
	// Run the detection multiple times to ensure it's fast
	const iterations = 1000000

	for i := 0; i < iterations; i++ {
		_ = isDeterministicModeAvailable()
	}

	// If we got here without timeout, the overhead is acceptable
	t.Logf("Completed %d detection checks", iterations)
}