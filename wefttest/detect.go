//go:build detsched

package wefttest

// isDeterministicModeAvailable returns true when compiled with -tags=detsched.
// This enables the deterministic scheduler for comprehensive concurrency testing.
func isDeterministicModeAvailable() bool {
	return true
}