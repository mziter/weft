# Weft Examples

This directory contains realistic examples demonstrating how to use Weft to test concurrent Go code. Each example shows production-quality components and their corresponding deterministic tests.

## Overview

The examples are designed to show:
- **Real production code** that uses Weft primitives
- **Testing strategies** for different concurrency patterns
- **Bug detection** that standard testing might miss
- **Migration patterns** from standard library to Weft

## Production Components

### `counter.go` - Thread-Safe Counter
A simple counter with proper synchronization and a deliberately buggy variant.

**Components:**
- `Counter`: Thread-safe counter with mutex protection
- `Increment()`: Correctly synchronized increment
- `IncrementWithWork()`: **Buggy** - has race condition between read and write

**Key Learning**: Shows how Weft exposes race conditions that occur when locks are released between read-modify-write operations.

### `bank.go` - Bank Account Transfers
Classic example for testing deadlock scenarios with multiple locks.

**Components:**
- `BankAccount`: Thread-safe account with reader-writer mutex
- `Transfer()`: **Buggy** - can deadlock due to lock ordering
- `SafeTransfer()`: Fixed version using consistent lock ordering

**Key Learning**: Demonstrates deadlock detection and prevention through proper lock ordering.

### `queue.go` - Producer-Consumer Pattern
Channel-based queue with producer-consumer coordination.

**Components:**
- `Queue[T]`: Generic thread-safe queue using channels
- `ProducerConsumer`: Demonstrates typical producer-consumer patterns

**Key Learning**: Shows how to test channel-based synchronization and coordination patterns.

### `worker_pool.go` - Worker Pool Coordination
Worker pool implementation using condition variables for job coordination.

**Components:**
- `WorkerPool`: Manages multiple workers with condition variable signaling
- Demonstrates job submission, worker coordination, and shutdown patterns

**Key Learning**: Tests condition variable usage, spurious wakeups, and worker synchronization.

## Test Files

### `counter_test.go`
**Purpose**: Demonstrates basic concurrent testing and race condition detection.

- `TestCounter`: Tests correct synchronization with properly locked increment
- `TestCounterRaceCondition`: Exposes race conditions in buggy increment method

**What it shows**: How Weft's deterministic scheduler can reliably expose race conditions that would be intermittent in production.

### `bank_test.go`
**Purpose**: Deadlock detection and prevention testing.

- `TestBankAccountTransfer`: Shows deadlock scenarios (currently skipped - needs scheduler implementation)
- `TestSafeBankAccountTransfer`: Tests deadlock-free implementation

**What it shows**: How to test for deadlocks and verify that fixes actually work across different execution orders.

### `queue_test.go`
**Purpose**: Channel synchronization testing.

- `TestProducerConsumer`: Tests producer-consumer coordination with channels

**What it shows**: How to verify channel-based synchronization works correctly across different interleavings.

### `worker_pool_test.go`
**Purpose**: Condition variable and worker coordination testing.

- `TestWorkerPool`: Tests job distribution and completion tracking

**What it shows**: How to test complex coordination patterns with multiple workers and condition variables.

### `replay_test.go`
**Purpose**: Bug reproduction and debugging.

- `TestReplaySpecificSeed`: Shows how to reproduce specific failing scenarios

**What it shows**: How to use deterministic replay for debugging once you find a failing seed.

## Running the Examples

### Two-Tier Testing Approach
```bash
cd examples

# Run basic tests (fast feedback, skips deterministic tests with helpful message)
go test -v

# Run with deterministic concurrency testing (comprehensive bug detection)
go test -tags=detsched -v
```

**How it works:**
- **Without `-tags=detsched`**: Tests using `wefttest.Explore()` skip with guidance to use the tag
- **With `-tags=detsched`**: Full deterministic scheduling explores multiple execution orders
- **Current Status**: Deterministic tests fail because scheduler is not yet implemented

### User-Friendly Detection
The examples demonstrate the detection approach:
- **Helpful skip messages** when deterministic mode isn't available
- **Clear guidance** on how to enable comprehensive testing
- **Zero performance impact** in production builds

## Expected Behavior (Once Scheduler is Complete)

### Working Tests
- `TestSafeBankAccountTransfer`: Should pass - demonstrates correct implementation
- `TestProducerConsumer`: Should pass - shows proper channel coordination
- `TestWorkerPool`: Should pass - demonstrates condition variable patterns

### Bug-Detecting Tests
- `TestCounterRaceCondition`: Should detect race conditions across multiple schedules
- `TestBankAccountTransfer`: Should detect deadlocks (once deadlock detection is implemented)

### Replay Tests
- `TestReplaySpecificSeed`: Should produce identical results every run with the same seed

## Migration Strategy

To convert existing concurrent code to use Weft:

1. **Replace imports**: `sync` → `weft`, `time` → `weft`
2. **Replace primitives**: `go func(){}` → `weft.Go(func(ctx weft.Context){})`
3. **Update channels**: `make(chan T, n)` → `weft.MakeChan[T](n)`
4. **Add exploration**: Wrap tests with `wefttest.Explore()`
5. **Run tests with deterministic mode**: `go test -tags=detsched ./...` - comprehensive concurrency testing!

## Common Patterns

### Basic Exploration
```go
wefttest.Explore(t, 100, func(s *weft.Scheduler) {
    // Your concurrent code using weft primitives
    s.Wait() // Wait for all spawned goroutines
})
```

### Bug Reproduction
```go
const failingSeed = 12345
wefttest.Replay(t, failingSeed, func(s *weft.Scheduler) {
    // Same code that failed - will behave identically
})
```

### Race Detection Pattern
```go
// Instead of expecting success, log discrepancies
if actual != expected {
    t.Logf("Potential race detected: expected %d, got %d", expected, actual)
}
```

## Implementation Status

✅ **Production code**: Realistic components with actual concurrency bugs
✅ **Test structure**: Proper separation of concerns and realistic usage patterns
⚠️ **Scheduler**: Stub implementation - tests fail until completed
⚠️ **Deadlock detection**: Not implemented - related tests are skipped

The examples serve as both documentation and a test suite for the Weft framework itself.