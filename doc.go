// Package weft provides a deterministic concurrency testing framework for Go.
//
// Weft offers drop-in replacements for Go's concurrency primitives that enable
// deterministic testing. By controlling the scheduling of goroutines and blocking
// operations, Weft can reproduce rare race conditions, explore different
// interleavings systematically, and detect issues like deadlocks and
// linearizability violations.
//
// Basic usage:
//
//	var mu weft.Mutex
//	counter := 0
//
//	for i := 0; i < 10; i++ {
//		weft.Go(func(ctx weft.Context) {
//			mu.Lock()
//			counter++
//			mu.Unlock()
//		})
//	}
//
// The framework uses build tags to switch between deterministic and standard
// implementations. Use -tags=detsched for deterministic mode.
package weft