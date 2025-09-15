//go:build !detsched

package wefttest

// isDeterministicModeAvailable returns false when compiled without -tags=detsched.
// In this mode, tests will skip with a helpful message about enabling deterministic testing.
func isDeterministicModeAvailable() bool {
	return false
}