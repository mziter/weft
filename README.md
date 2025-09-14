# Weft

A deterministic concurrency testing framework for Go that reliably reproduces, explores, and shrinks concurrency bugs that the Go race detector cannot detect.

## Overview

Weft provides drop-in replacements for Go's concurrency primitives that enable deterministic testing. By controlling the scheduling of goroutines and blocking operations, Weft can:

- Reproduce rare race conditions and deadlocks consistently
- Explore different interleavings systematically
- Shrink failing traces to minimal reproductions
- Detect issues like deadlocks, livelocks, and linearizability violations

## Features

- **Deterministic Scheduling**: Seed-driven scheduler ensures reproducible test runs
- **Drop-in Replacements**: Mirrors standard library APIs for easy adoption
- **Multiple Exploration Strategies**: Random, round-robin, and bounded exploration
- **Virtual Time**: Logical clock eliminates timing-related flakiness
- **Zero Overhead in Production**: Build tags enable seamless switching between deterministic and standard implementations
- **Comprehensive Primitives**: Support for goroutines, mutexes, channels, timers, and condition variables

## Installation

```bash
go get github.com/yourusername/weft
```

## Quick Start

Replace standard library imports with Weft equivalents where you want deterministic testing:

```go
package mypackage

import (
    "testing"
    "github.com/yourusername/weft"
    "github.com/yourusername/weft/wefttest"
)

func TestConcurrentCounter(t *testing.T) {
    // Explore 100 different schedules
    wefttest.Explore(t, 100, func(s *weft.Scheduler) {
        var mu weft.Mutex
        counter := 0

        for i := 0; i < 10; i++ {
            weft.Go(func(ctx weft.Context) {
                mu.Lock()
                counter++
                mu.Unlock()
            })
        }

        s.Wait() // Wait for all goroutines

        if counter != 10 {
            t.Errorf("expected counter=10, got %d", counter)
        }
    })
}
```

## Usage

### Development/Testing
```bash
# Run with deterministic scheduler
go test -tags=detsched ./...

# Explore multiple schedules
WEFT_RUNS=1000 go test -tags=detsched ./...

# Replay a specific failure
WEFT_SEED=12345 go test -tags=detsched ./...
```

### Production
```bash
# Run normally (uses standard library)
go test ./...
go build ./...
```

## API Reference

### Core Primitives

- `weft.Go(func(Context))` - Spawn a deterministic goroutine
- `weft.Sleep(duration)` - Deterministic sleep
- `weft.After(duration)` - Deterministic timer
- `weft.Mutex` / `weft.RWMutex` - Deterministic mutexes
- `weft.NewCond(*Mutex)` - Deterministic condition variable
- `weft.MakeChan[T](capacity)` - Deterministic channel

### Testing Helpers

- `wefttest.Explore(t, runs, buildFn)` - Explore multiple schedules
- `wefttest.Replay(t, seed, buildFn)` - Replay specific seed
- `wefttest.ReplayChoices(t, choices, buildFn)` - Replay trace

## What Weft Catches

Weft detects concurrency bugs that the race detector cannot:

- **Deadlocks**: Circular waits, incorrect condition variable usage
- **Lost Updates**: Non-atomic read-modify-write sequences
- **Linearizability Violations**: Inconsistent operation ordering
- **Starvation**: Unfair scheduling, priority inversion
- **Protocol Violations**: Incorrect synchronization patterns

## Documentation

- [Getting Started Guide](docs/getting-started.md)
- [API Documentation](docs/api.md)
- [Examples](examples/)
- [Migration Guide](docs/migration.md)

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

Inspired by deterministic testing frameworks like Loom (Rust) and similar tools in other ecosystems.