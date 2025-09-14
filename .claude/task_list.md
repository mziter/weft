# Weft Detailed Task List

This document provides detailed technical tasks organized by milestone and component. Each task includes specific requirements, acceptance criteria, and implementation guidance.

---

## 🎯 MILESTONE 1: Core Foundation

### Task 1.1: High-Quality PRNG Implementation
**File**: `internal/prng/prng.go`
**Complexity**: Medium
**Dependencies**: None

#### Requirements
- Implement xoshiro256** or PCG64 algorithm for high-quality deterministic randomness
- Provide `Uint64()`, `Intn(n int) int`, and `Float64()` methods
- Ensure identical output across platforms for same seed
- Optimize for speed (will be called frequently)

#### Implementation Details
```go
type PRNG struct {
    state [4]uint64  // for xoshiro256**
}

func New(seed uint64) *PRNG
func (p *PRNG) Uint64() uint64
func (p *PRNG) Intn(n int) int
func (p *PRNG) Float64() float64
func (p *PRNG) Shuffle(slice []int)  // for randomizing waiter selection
```

#### Acceptance Criteria
- [ ] Same seed produces identical sequence across runs
- [ ] Same seed produces identical sequence across platforms
- [ ] Passes statistical randomness tests
- [ ] Performance: >100M operations/second
- [ ] Full test coverage with deterministic output verification

#### Files to Modify
- `internal/prng/prng.go` - Complete implementation
- `internal/prng/prng_test.go` - Create comprehensive tests

---

### Task 1.2: Task Management System
**File**: `internal/scheduler/task.go`
**Complexity**: Medium
**Dependencies**: Task 1.1

#### Requirements
- Implement task lifecycle management (Ready, Running, Blocked, Done)
- Support task spawning with proper context
- Handle task yielding and resumption
- Track task relationships and dependencies

#### Implementation Details
```go
type Task struct {
    id       int64
    state    TaskState
    fn       func(Context)
    ctx      *taskContext
    blocker  Blocker        // what this task is blocked on
    wakeup   chan struct{}  // for unblocking
}

type TaskState int
const (
    TaskReady TaskState = iota
    TaskRunning
    TaskBlocked
    TaskDone
)

type Blocker interface {
    Block(task *Task)
    Unblock(task *Task)
    String() string  // for debugging
}
```

#### Acceptance Criteria
- [ ] Tasks transition through states correctly
- [ ] Blocked tasks don't consume CPU
- [ ] Task context provides yield and cancellation
- [ ] Memory usage scales linearly with active tasks
- [ ] Race-free task state management

#### Files to Modify
- `internal/scheduler/task.go` - Complete task management
- `internal/scheduler/context.go` - Create task context implementation
- `context.go` - Update Context interface if needed

---

### Task 1.3: Core Scheduler Engine
**File**: `internal/scheduler/scheduler.go`
**Complexity**: High
**Dependencies**: Tasks 1.1, 1.2

#### Requirements
- Deterministic task scheduling using PRNG
- Cooperative multitasking with yield points
- Proper task spawning and termination
- Deadlock detection preparation (hooks for later milestone)

#### Implementation Details
```go
type Scheduler struct {
    mu          sync.Mutex
    prng        *prng.PRNG
    tasks       map[int64]*Task
    runnable    []*Task
    blocked     map[Blocker][]*Task
    nextTaskID  int64
    running     bool
    clock       *VirtualClock
}

func (s *Scheduler) Spawn(fn func(Context)) int64
func (s *Scheduler) Yield()
func (s *Scheduler) Block(blocker Blocker)
func (s *Scheduler) Unblock(task *Task)
func (s *Scheduler) selectNext() *Task  // deterministic selection
func (s *Scheduler) run() // main scheduling loop
```

#### Acceptance Criteria
- [ ] Same seed produces identical task execution order
- [ ] Tasks yield cooperatively at blocking operations
- [ ] Scheduler handles task completion properly
- [ ] No race conditions in scheduler state
- [ ] Graceful handling of panics in tasks
- [ ] Memory cleanup when tasks complete

#### Files to Modify
- `internal/scheduler/scheduler.go` - Complete scheduler implementation
- `weft.go` - Update to use real scheduler instead of stub

---

### Task 1.4: Virtual Time System
**File**: `internal/scheduler/clock.go`
**Complexity**: Medium
**Dependencies**: Task 1.3

#### Requirements
- Logical time that advances only at yield points
- Timer management for Sleep() and After() operations
- Deterministic timer ordering when multiple timers expire simultaneously
- Integration with scheduler for timer-based task unblocking

#### Implementation Details
```go
type VirtualClock struct {
    mu      sync.Mutex
    now     time.Duration  // logical time
    timers  *timerHeap     // priority queue of pending timers
}

type Timer struct {
    expiry   time.Duration
    task     *Task
    callback func()
    ch       chan time.Time  // for After() operations
}

func (c *VirtualClock) Now() time.Duration
func (c *VirtualClock) Sleep(d time.Duration) // blocks current task
func (c *VirtualClock) After(d time.Duration) <-chan time.Time
func (c *VirtualClock) Advance() time.Duration  // advance to next event
func (c *VirtualClock) processExpiredTimers() []*Task
```

#### Acceptance Criteria
- [ ] Time advances deterministically
- [ ] Multiple timers with same expiry fire in deterministic order
- [ ] Sleep() blocks correctly and unblocks at right time
- [ ] After() channels receive at correct logical time
- [ ] No wall-clock dependency in deterministic mode

#### Files to Modify
- `internal/scheduler/clock.go` - Create virtual time implementation
- `weft.go` - Update Sleep() and After() to use virtual time
- `clock.go` - Complete Clock interface if needed

---

### Task 1.5: Fix Type Compatibility Issues
**Files**: Multiple
**Complexity**: Low
**Dependencies**: Tasks 1.2, 1.3

#### Requirements
- Resolve function signature mismatches between stubs and implementation
- Ensure Context interface works with both modes
- Fix channel type handling for generic channels
- Update all primitive wrappers to use correct types

#### Implementation Details
- Update `Scheduler.Spawn(fn func(Context))` signature
- Ensure `Context` interface compatibility
- Fix channel Send/Recv methods for proper typing
- Update mutex/condition variable integration

#### Acceptance Criteria
- [ ] All code compiles with `-tags=detsched`
- [ ] All code compiles without tags (production mode)
- [ ] No type assertion failures at runtime
- [ ] Examples run without compilation errors

#### Files to Modify
- `internal/scheduler/scheduler.go` - Fix Spawn signature
- `weft.go` - Update function signatures
- `channel.go` - Fix generic type handling
- All primitive files - Ensure type consistency

---

## 🔒 MILESTONE 2: Synchronization Primitives

### Task 2.1: Deterministic Mutex Implementation
**File**: `internal/scheduler/mutex.go`
**Complexity**: Medium
**Dependencies**: Milestone 1

#### Requirements
- FIFO waiter queue with deterministic ordering
- Lock ownership tracking
- Deadlock detection hooks (for milestone 3)
- Support for spurious wakeup simulation
- TryLock() implementation

#### Implementation Details
```go
type Mutex struct {
    mu       sync.Mutex  // protects internal state
    locked   bool
    owner    int64       // task ID that owns the lock
    waiters  []*Task     // FIFO queue of waiting tasks
    scheduler *Scheduler // reference for blocking/unblocking
}

func (m *Mutex) Lock()
func (m *Mutex) Unlock()
func (m *Mutex) TryLock() bool
func (m *Mutex) block(task *Task)
func (m *Mutex) unblockNext()
```

#### Acceptance Criteria
- [ ] Lock/Unlock operations are deterministic
- [ ] Waiter queue maintains FIFO ordering
- [ ] Recursive locking panics appropriately
- [ ] Unlocking unlocked mutex panics
- [ ] TryLock() never blocks and returns correct status
- [ ] Integrates properly with scheduler blocking

#### Files to Modify
- `internal/scheduler/mutex.go` - Complete implementation
- `mutex.go` - Wire up to use scheduler mutex
- `mutex_notag.go` - Ensure production mode unchanged

---

### Task 2.2: RWMutex Implementation
**File**: `internal/scheduler/rwmutex.go`
**Complexity**: Medium
**Dependencies**: Task 2.1

#### Requirements
- Reader-writer lock semantics with deterministic waiter selection
- Separate read and write waiter queues
- Configurable reader/writer priority policies
- Prevent writer starvation

#### Implementation Details
```go
type RWMutex struct {
    mu          sync.Mutex
    readers     int         // number of active readers
    writer      int64       // task ID of active writer (0 if none)
    readWaiters []*Task
    writeWaiters []*Task
    scheduler   *Scheduler
}

func (rw *RWMutex) RLock()
func (rw *RWMutex) RUnlock()
func (rw *RWMutex) Lock()
func (rw *RWMutex) Unlock()
func (rw *RWMutex) selectNextWaiter() *Task  // deterministic selection
```

#### Acceptance Criteria
- [ ] Multiple readers can hold lock simultaneously
- [ ] Writer excludes all readers and other writers
- [ ] Deterministic waiter selection prevents starvation
- [ ] Correct panic behavior for misuse
- [ ] Performance: minimal overhead for uncontended case

#### Files to Modify
- `internal/scheduler/rwmutex.go` - Create implementation
- `mutex.go` - Add RWMutex integration

---

### Task 2.3: Condition Variable Implementation
**File**: `internal/scheduler/cond.go`
**Complexity**: High
**Dependencies**: Task 2.1

#### Requirements
- Wait() releases lock and blocks atomically
- Signal() wakes exactly one waiter deterministically
- Broadcast() wakes all waiters
- Spurious wakeup simulation capability
- Integration with mutex for atomic lock release/reacquisition

#### Implementation Details
```go
type Cond struct {
    mu       sync.Mutex
    lock     Locker      // associated mutex
    waiters  []*Task
    scheduler *Scheduler
}

func NewCond(l Locker) *Cond
func (c *Cond) Wait()        // must hold lock
func (c *Cond) Signal()      // wake one waiter
func (c *Cond) Broadcast()   // wake all waiters
func (c *Cond) selectWaiter() *Task  // deterministic selection
```

#### Acceptance Criteria
- [ ] Wait() atomically releases and reacquires lock
- [ ] Signal() wakes exactly one waiter when waiters exist
- [ ] Broadcast() wakes all waiters
- [ ] Spurious wakeups can be simulated deterministically
- [ ] Works correctly with both Mutex and RWMutex
- [ ] Proper panic behavior when lock not held

#### Files to Modify
- `internal/scheduler/cond.go` - Complete implementation
- `cond.go` - Wire up to scheduler cond
- `cond_notag.go` - Ensure production mode works

---

### Task 2.4: Deterministic Channel Implementation
**File**: `internal/scheduler/channel.go`
**Complexity**: High
**Dependencies**: Milestone 1

#### Requirements
- Buffered and unbuffered channel semantics
- Deterministic sender/receiver selection when multiple are ready
- Channel closing semantics with proper panic behavior
- Select-like operations preparation (for future milestone)
- Generic type support

#### Implementation Details
```go
type Chan[T any] struct {
    mu        sync.Mutex
    buf       []T         // circular buffer
    cap       int         // buffer capacity
    head      int         // buffer head
    tail      int         // buffer tail
    closed    bool
    senders   []*channelOp[T]  // blocked senders
    receivers []*channelOp[T]  // blocked receivers
    scheduler *Scheduler
}

type channelOp[T any] struct {
    task  *Task
    value T
    ok    *bool  // for receive operations
}

func MakeChan[T any](cap int) *Chan[T]
func (c *Chan[T]) Send(v T)
func (c *Chan[T]) Recv() (T, bool)
func (c *Chan[T]) TrySend(v T) bool
func (c *Chan[T]) TryRecv() (T, bool)
func (c *Chan[T]) Close()
```

#### Acceptance Criteria
- [ ] Unbuffered channels block until matched sender/receiver
- [ ] Buffered channels block only when full/empty
- [ ] Closed channel semantics match standard library exactly
- [ ] Deterministic selection when multiple senders/receivers ready
- [ ] TrySend/TryRecv never block
- [ ] Proper panic behavior for send on closed channel

#### Files to Modify
- `internal/scheduler/channel.go` - Complete implementation
- `channel.go` - Wire up to scheduler channels
- `channel_notag.go` - Ensure production mode works

---

## 🕵️ MILESTONE 3: Advanced Detection & Debugging

### Task 3.1: Deadlock Detection Algorithm
**File**: `internal/scheduler/deadlock.go`
**Complexity**: High
**Dependencies**: Milestone 2

#### Requirements
- Detect when no tasks can make progress
- Identify circular wait chains
- Generate helpful error messages with task stack traces
- Configurable timeout before declaring deadlock

#### Implementation Details
```go
type DeadlockDetector struct {
    scheduler *Scheduler
}

func (d *DeadlockDetector) Check() *DeadlockInfo
func (d *DeadlockDetector) findCycles() [][]int64
func (d *DeadlockDetector) buildWaitGraph() map[int64][]int64

type DeadlockInfo struct {
    Tasks  []int64
    Cycle  []int64  // if circular wait detected
    Resources []string  // names of locked resources
    Message string
}
```

#### Acceptance Criteria
- [ ] Detects simple circular waits (A waits for B, B waits for A)
- [ ] Detects complex multi-resource deadlocks
- [ ] Provides clear error messages with task information
- [ ] Does not false positive on temporary blocking
- [ ] Performance: detection runs in reasonable time

#### Files to Modify
- `internal/scheduler/deadlock.go` - Create detector
- `internal/scheduler/scheduler.go` - Integrate deadlock detection
- `wefttest/explore.go` - Handle deadlock errors

---

### Task 3.2: Timeout Detection and Handling
**File**: `internal/scheduler/timeout.go`
**Complexity**: Medium
**Dependencies**: Task 1.4 (Virtual Time)

#### Requirements
- Configurable step count limits
- Configurable logical time limits
- Graceful test termination on timeout
- Clear timeout error messages

#### Implementation Details
```go
type TimeoutDetector struct {
    maxSteps    int64
    maxTime     time.Duration
    currentStep int64
    startTime   time.Duration
}

func (t *TimeoutDetector) CheckStep() error
func (t *TimeoutDetector) CheckTime(current time.Duration) error
func (t *TimeoutDetector) Reset()
```

#### Acceptance Criteria
- [ ] Step limits prevent infinite loops
- [ ] Time limits prevent long-running tests
- [ ] Configurable via environment variables
- [ ] Clear error messages indicate timeout cause
- [ ] Minimal performance overhead

#### Files to Modify
- `internal/scheduler/timeout.go` - Create implementation
- `internal/scheduler/scheduler.go` - Integrate timeout checking
- `wefttest/explore.go` - Handle timeout configuration

---

### Task 3.3: Trace Recording System
**File**: `internal/scheduler/trace.go`
**Complexity**: High
**Dependencies**: Milestone 2

#### Requirements
- Record all scheduling decisions and operations
- Compact binary format for efficiency
- Support for trace filtering and analysis
- Integration with replay system

#### Implementation Details
```go
type TraceEvent struct {
    Timestamp  int64
    TaskID     int64
    Operation  string
    ResourceID string
    Metadata   map[string]interface{}
}

type TraceRecorder struct {
    events []TraceEvent
    enabled bool
    filter  func(TraceEvent) bool
}

func (t *TraceRecorder) Record(event TraceEvent)
func (t *TraceRecorder) Save(filename string) error
func (t *TraceRecorder) Load(filename string) error
func (t *TraceRecorder) Replay(scheduler *Scheduler) error
```

#### Acceptance Criteria
- [ ] Records all relevant scheduling events
- [ ] Trace files are compact and portable
- [ ] Replay produces identical execution
- [ ] Configurable event filtering for performance
- [ ] Integration with wefttest utilities

#### Files to Modify
- `internal/scheduler/trace.go` - Create recorder
- `internal/scheduler/scheduler.go` - Add trace recording
- `wefttest/replay.go` - Implement trace replay

---

### Task 3.4: Environment Configuration
**Files**: Multiple
**Complexity**: Low
**Dependencies**: Tasks 3.1, 3.2, 3.3

#### Requirements
- `WEFT_RUNS` - number of exploration runs
- `WEFT_SEED` - pin specific seed
- `WEFT_TRACE` - enable trace recording
- `WEFT_MAX_STEPS` - step count limit
- `WEFT_TIMEOUT_MS` - logical time limit
- `WEFT_DEADLOCK_TIMEOUT` - deadlock detection timeout

#### Implementation Details
```go
type Config struct {
    Runs            int
    Seed            uint64
    TraceEnabled    bool
    MaxSteps        int64
    TimeoutMs       int64
    DeadlockTimeout time.Duration
}

func LoadConfig() Config
func (c Config) Apply(scheduler *Scheduler)
```

#### Acceptance Criteria
- [ ] All environment variables work correctly
- [ ] Sensible defaults when variables not set
- [ ] Configuration validation and error messages
- [ ] Documentation for all options

#### Files to Modify
- `config.go` - Create configuration system
- `wefttest/explore.go` - Use configuration
- `README.md` - Document environment variables

---

## 🧪 MILESTONE 4: Testing & Production Readiness

### Task 4.1: Comprehensive Test Suite
**Files**: Multiple `*_test.go`
**Complexity**: High
**Dependencies**: Milestone 3

#### Requirements
- Unit tests for all components with >90% coverage
- Integration tests for full workflows
- Performance benchmarks comparing modes
- Regression tests for known issues

#### Test Categories
1. **Unit Tests**: Each component in isolation
2. **Integration Tests**: Full scheduler workflows
3. **Determinism Tests**: Same seed = same result
4. **Performance Tests**: Overhead measurements
5. **Edge Case Tests**: Error conditions, boundary cases

#### Acceptance Criteria
- [ ] >90% test coverage across all packages
- [ ] All tests pass in both deterministic and production modes
- [ ] Performance tests show <5% overhead in production mode
- [ ] Regression tests prevent known issues
- [ ] Tests run quickly in CI/CD

#### Files to Modify
- Create comprehensive `*_test.go` files for all packages
- Add benchmark tests
- Set up test coverage reporting

---

### Task 4.2: Performance Optimization
**Files**: Multiple
**Complexity**: Medium
**Dependencies**: Task 4.1

#### Requirements
- Profile deterministic mode for hotspots
- Optimize scheduler data structures
- Minimize memory allocations
- Ensure production mode has zero overhead

#### Optimization Targets
1. **Scheduler**: Fast task selection algorithms
2. **PRNG**: Optimized random number generation
3. **Primitives**: Minimal locking overhead
4. **Memory**: Efficient task and resource management

#### Acceptance Criteria
- [ ] Production mode has <1% overhead vs standard library
- [ ] Deterministic mode runs at reasonable speed for testing
- [ ] Memory usage scales linearly with active tasks
- [ ] No memory leaks in long-running tests

#### Files to Modify
- All implementation files - Apply optimizations
- Add benchmarking and profiling tools

---

### Task 4.3: Documentation and Examples
**Files**: Multiple documentation files
**Complexity**: Medium
**Dependencies**: Working implementation

#### Requirements
- Complete API documentation with examples
- Migration guide from standard library
- Troubleshooting guide for common issues
- Best practices for deterministic testing

#### Documentation Structure
1. **API Reference**: Complete function documentation
2. **User Guide**: How to use Weft effectively
3. **Migration Guide**: Converting existing code
4. **Examples**: Real-world usage patterns
5. **Troubleshooting**: Common issues and solutions

#### Acceptance Criteria
- [ ] All public APIs have complete documentation
- [ ] Examples demonstrate key features
- [ ] Migration guide enables successful adoption
- [ ] Troubleshooting guide addresses common issues

#### Files to Modify
- `docs/` directory - Create comprehensive documentation
- Update `README.md` with complete information
- Enhance `examples/README.md`

---

### Task 4.4: Error Messages and Debugging
**Files**: Multiple
**Complexity**: Medium
**Dependencies**: Working implementation

#### Requirements
- Clear, actionable error messages for common mistakes
- Debug utilities for analyzing failures
- Integration with Go's testing framework
- Helpful output for deadlock and timeout scenarios

#### Debug Features
1. **Error Messages**: Clear explanations of failures
2. **Stack Traces**: Show where problems occur
3. **State Inspection**: Debug scheduler and task state
4. **Trace Analysis**: Tools for understanding execution

#### Acceptance Criteria
- [ ] Error messages guide users toward solutions
- [ ] Debug output helps identify issues quickly
- [ ] Integration with `go test -v` provides useful information
- [ ] Failure scenarios include sufficient context

#### Files to Modify
- All implementation files - Improve error messages
- Add debug utilities and formatting

---

## 🚀 MILESTONE 5: Extended Features & Tools

### Task 5.1: Trace Shrinking Algorithm
**File**: `internal/scheduler/shrink.go`
**Complexity**: High
**Dependencies**: Task 3.3 (Trace Recording)

#### Requirements
- Automatically reduce failing traces to minimal reproductions
- Binary search and delta debugging approaches
- Preserve failure behavior while removing irrelevant events
- Integration with wefttest utilities

#### Implementation Details
```go
type TraceShrinkker struct {
    originalTrace []TraceEvent
    testFunc     func([]TraceEvent) bool  // returns true if still fails
}

func (s *TraceShrinkker) Shrink() []TraceEvent
func (s *TraceShrinkker) binarySearch() []TraceEvent
func (s *TraceShrinkker) deltaDebug() []TraceEvent
```

#### Acceptance Criteria
- [ ] Produces minimal failing traces
- [ ] Preserves original failure behavior
- [ ] Runs in reasonable time
- [ ] Integrates with wefttest workflows

#### Files to Modify
- `internal/scheduler/shrink.go` - Create shrinking algorithm
- `wefttest/` - Add shrinking utilities

---

### Task 5.2: Multiple Scheduling Policies
**File**: `internal/scheduler/policies.go`
**Complexity**: Medium
**Dependencies**: Milestone 2

#### Requirements
- Random scheduling (current default)
- Round-robin scheduling for fairness testing
- Priority-based scheduling
- Bounded exploration with systematic coverage

#### Implementation Details
```go
type SchedulingPolicy interface {
    SelectNext(tasks []*Task, prng *prng.PRNG) *Task
    Name() string
}

type RandomPolicy struct{}
type RoundRobinPolicy struct { lastIndex int }
type PriorityPolicy struct { priorities map[int64]int }
type BoundedExplorationPolicy struct { choices []int }
```

#### Acceptance Criteria
- [ ] Multiple policies available via configuration
- [ ] Each policy has distinct testing characteristics
- [ ] Deterministic behavior within each policy
- [ ] Documentation explains when to use each policy

#### Files to Modify
- `internal/scheduler/policies.go` - Create policies
- `internal/scheduler/scheduler.go` - Integrate policy selection
- `config.go` - Add policy configuration

---

### Task 5.3: weftfix Tool Implementation
**File**: `cmd/weftfix/main.go`
**Complexity**: High
**Dependencies**: Working framework

#### Requirements
- Parse Go source code and identify concurrency primitives
- Safely replace standard library calls with weft equivalents
- Handle imports and package references
- Support dry-run mode and rollback

#### Transformation Rules
1. `go func(){}` → `weft.Go(func(ctx weft.Context){})`
2. `time.Sleep()` → `weft.Sleep()`
3. `sync.Mutex` → `weft.Mutex`
4. `make(chan T, n)` → `weft.MakeChan[T](n)`
5. Import management

#### Acceptance Criteria
- [ ] Correctly identifies and transforms concurrency code
- [ ] Preserves code behavior and semantics
- [ ] Handles edge cases and complex expressions
- [ ] Provides clear feedback on transformations
- [ ] Supports undo/rollback operations

#### Files to Modify
- `cmd/weftfix/main.go` - Complete implementation
- `cmd/weftfix/transform.go` - Transformation logic
- `cmd/weftfix/parser.go` - Go AST parsing

---

### Task 5.4: CI/CD Integration and Templates
**Files**: CI configuration templates
**Complexity**: Low
**Dependencies**: Working framework

#### Requirements
- GitHub Actions workflow templates
- Integration with popular CI systems
- Systematic exploration across different seeds
- Performance regression detection

#### Templates
1. **Basic Testing**: Run tests in both modes
2. **Exploration**: Systematic seed exploration
3. **Performance**: Benchmark regression detection
4. **Nightly**: Extended testing with many seeds

#### Acceptance Criteria
- [ ] Templates work with major CI systems
- [ ] Provide good coverage without excessive runtime
- [ ] Detect performance regressions
- [ ] Easy to customize for different projects

#### Files to Create
- `.github/workflows/weft-ci.yml` - GitHub Actions template
- `ci/templates/` - Templates for other CI systems
- Documentation for CI setup

---

This detailed task list provides concrete implementation guidance for each milestone while maintaining flexibility for iterative development and refinement.