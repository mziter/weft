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
go get github.com/mziter/weft
```

## Quick Start

### Step 1: Write Your Production Code

First, create your concurrent component using Weft primitives:

```go
// counter.go - Production code using Weft primitives
package myapp

import "github.com/mziter/weft"

// Counter is a thread-safe counter for your application
type Counter struct {
    mu    weft.Mutex
    value int
}

func NewCounter() *Counter {
    return &Counter{}
}

// Increment safely increments the counter
func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

// Value safely returns the current counter value
func (c *Counter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}
```

### Step 2: Write Deterministic Tests

Create tests that automatically use deterministic concurrency testing:

```go
// counter_test.go - Comprehensive concurrency tests
package myapp

import (
    "sync"
    "testing"
    "github.com/mziter/weft"
    "github.com/mziter/weft/wefttest"
)

func TestCounterBasic(t *testing.T) {
    counter := NewCounter()

    // Test basic functionality
    counter.Increment()
    if counter.Value() != 1 {
        t.Errorf("expected 1, got %d", counter.Value())
    }
}

func TestCounterConcurrency(t *testing.T) {
    // wefttest.Explore runs deterministically when -tags=detsched is used
    // If the tag is missing, it will skip with a helpful message
    wefttest.Explore(t, 100, func(s *weft.Scheduler) {
        counter := NewCounter()

        // Spawn 10 concurrent workers using deterministic scheduling
        numWorkers := 10
        for i := 0; i < numWorkers; i++ {
            s.Go(func(ctx weft.Context) {
                // Each worker increments the counter once
                counter.Increment()
            })
        }

        // Wait for all workers to complete
        s.Wait()

        // Verify the final result - this will catch race conditions
        // that standard testing might miss
        if counter.Value() != numWorkers {
            t.Errorf("expected counter=%d, got %d", numWorkers, counter.Value())
        }
    })
}

func TestCounterReproduceBug(t *testing.T) {
    // If wefttest.Explore found a bug with seed 12345, replay it exactly
    const failingSeed = 12345

    wefttest.Replay(t, failingSeed, func(s *weft.Scheduler) {
        counter := NewCounter()

        for i := 0; i < 5; i++ {
            s.Go(func(ctx weft.Context) {
                counter.Increment()
            })
        }

        s.Wait()

        // This will consistently reproduce the same behavior every time
        // allowing you to debug the specific problematic interleaving
        t.Logf("With seed %d, final counter value: %d", failingSeed, counter.Value())
    })
}

func TestCounterConcurrentStandard(t *testing.T) {
    counter := NewCounter()

    // Optional: You can still write standard concurrent tests
    // These run fast but may not catch timing-dependent bugs
    var wg sync.WaitGroup
    numWorkers := 10

    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Increment()
        }()
    }

    wg.Wait()

    if counter.Value() != numWorkers {
        t.Errorf("expected %d, got %d", numWorkers, counter.Value())
    }
}
```

### Step 3: Run Your Tests

You have two modes of testing:

```bash
# Run basic tests (fast feedback, basic functionality)
go test ./...

# Run with deterministic concurrency testing (thorough, catches subtle bugs)
go test -tags=detsched ./...
```

**What happens:**
- **Without `-tags=detsched`**: Tests with `wefttest.Explore()` skip with a helpful message
- **With `-tags=detsched`**: Comprehensive deterministic concurrency testing runs

### Advanced Usage

```bash
# Explore more schedules to increase bug-finding probability
WEFT_RUNS=500 go test -tags=detsched ./...

# Reproduce a specific failure found during exploration
WEFT_SEED=12345 go test -tags=detsched ./...

# Run with verbose output to see test details
go test -tags=detsched -v ./...

# Run with race detection for additional safety
go test -tags=detsched -v -race ./...
```

### How It Works

- **Your production code** uses Weft primitives (`weft.Mutex`, etc.)
- **In production builds**: Weft primitives automatically use standard library (zero overhead)
- **In tests without `-tags=detsched`**: Tests with `wefttest.Explore()` skip with helpful guidance
- **In tests with `-tags=detsched`**: Full deterministic concurrency testing with scheduler exploration

### CI/CD Integration

Set up your CI pipeline to run both test modes:

```yaml
# .github/workflows/test.yml
- name: Run basic tests
  run: go test -v ./...

- name: Run deterministic concurrency tests
  run: go test -tags=detsched -v ./...

# Optional: Extended exploration for thorough testing
- name: Run extended concurrency tests
  run: WEFT_RUNS=1000 go test -tags=detsched -v ./...
```

### Production Deployment

Your application runs normally in production with zero overhead:

```bash
# Production builds automatically use standard library
go build ./...
go run ./...
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