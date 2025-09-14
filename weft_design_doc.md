# Weft: Deterministic Scheduler for Go

## 1. Project Goal

Provide Go developers with a **deterministic concurrency testing framework** that reliably **reproduces, explores, and shrinks** concurrency bugs that the Go race detector (`-race`) cannot detect.

The framework will:

- Offer **drop-in replacements** for a subset of standard concurrency primitives: `go` (spawn), `time.Sleep`, `time.After`, `sync.Mutex`, `sync.RWMutex`, `sync.Cond`, and channels.
- Use a **deterministic, seed-driven scheduler** to control interleavings of goroutines and blocking operations in tests.
- Support **exploration across multiple seeds** to uncover rare interleavings.
- Emit a **replayable seed and/or trace** for any failure, making bugs reproducible.

---

## 2. Scope

### In-Scope

- Test-only deterministic scheduler and concurrency primitives.
- In-kind API surface that mirrors `sync`, `time`, and channels closely.
- Deterministic replacements for goroutine spawning, mutexes, condition variables, timers, and channels.
- Multiple scheduling policies: **random (seeded)**, **round-robin**, **bounded exploration**.
- Reproducibility via **seed logging** and **trace replay**.
- Helpers that **try many seeds automatically**.
- **Deadlock** and **timeout** detection, with **trace shrinking** on failure.

### Out-of-Scope

- Transparent replacement of the Go runtime scheduler (cannot intercept raw `go`).
- Arbitrary unmodified application code (users must import and use `weft` wrappers where they want determinism).
- Syscall/cgo scheduling or OS thread interposition.

---

## 3. Functional Requirements

### 3.1 Scheduler

1. **Task Management**
   - `weft.Go(func(weft.Context))` starts a deterministic task.
   - Tasks **yield cooperatively** at blocking ops: `Sleep`, channel send/recv/close, mutex lock/unlock, cond wait/signal/broadcast.
   - The scheduler selects the **next runnable task**.

2. **Randomized Determinism (Seeds)**
   - A `uint64` seed initializes the PRNG that drives **all** nondeterministic choices.
   - **Same seed ⇒ same interleaving** and same generated test data (when integrated with property/fuzz generators via split seeds).

3. **Exploration**
   - `wefttest.Explore(t, N, buildFn)` executes **N** schedules with different seeds.
   - On the **first failure**, stop, **report the seed**, and **persist a trace** (if enabled).

4. **Replay**
   - Reproduce a failing run using `WithSeed(seed)` **or** an **explicit choice list** (`wefttest.ReplayChoices`).
   - Optional **shrunk trace** for faster, minimal reproductions.

5. **Virtual Time**
   - Provide a **logical clock** (`weft.Clock`) that advances only at timer/yield points.
   - `weft.Sleep(d)` and `weft.After(d)` fire based on logical time (no wall-clock flake).

6. **Deadlock & Timeout Detection**
   - If no tasks are runnable and unfinished tasks remain, **report deadlock**.
   - Configurable **max steps** and **logical timeout** per test.

### 3.2 Concurrency Primitives

1. **Mutex / RWMutex / Cond**
   - APIs mirror `sync.Mutex`, `sync.RWMutex`, and `sync.Cond`.
   - Correct **waiter queues**; **spurious wakeups** can be simulated.
   - `weft.NewCond(&weft.Mutex)` returns a determinstic cond bound to the given mutex.

2. **Channel**
   - `weft.MakeChan[T](cap int)` returns a deterministic channel with send/recv/close.
   - Blocking/unblocking are **scheduled deterministically** (e.g., which waiter wins).

3. **Task Context**
   - `weft.Context` is passed to `weft.Go` closures.
   - Offers `Yield()`, cancellation plumbing (optional), and debug helpers.

---

## 4. Non-Functional Requirements

- **Performance (Prod Mode)**
  - With `!detsched` (default builds), wrappers **inline** to stdlib (`sync`, `time`, native `chan`). Overhead should be negligible.

- **Usability**
  - Minimal code changes: replace `go` → `weft.Go`, `time.Sleep` → `weft.Sleep`, `time.After` → `weft.After`, `sync.*` → `weft.*`, `make(chan T, n)` → `weft.MakeChan[T](n)` where determinism is desired.
  - Provide `weftfix` **codemod** to automate safe, scoped rewrites with `--dry-run` and path globs.

- **Determinism**
  - Same seed ⇒ identical behavior across machines/OSes.
  - All nondeterministic decisions (scheduling, waiter selection, timer ordering) draw from a **single PRNG**.

- **Configurability**
  - Build tags: `//go:build detsched` enables deterministic mode; `!detsched` uses stdlib.
  - Env/flags: `WEFT_RUNS` (exploration runs), `WEFT_SEED` (pin seed), `WEFT_TRACE=1` (emit trace), `WEFT_MAX_STEPS`, `WEFT_TIMEOUT_MS`.

- **Test Integration**
  - Works with `go test`.
  - On failure, logs **seed** (and optional **trace path**).

---

## 5. Technical Details

- **PRNG**: `xoshiro256**` or `PCG64` (fast, high-quality, deterministic).
- **Trace Format (internal)**: `[{ts, taskID, op, resourceID, meta...}]`.
- **Shrinking**: Iteratively remove/merge choices until minimal failing trace remains.
- **Fairness Options**: Round-robin, random, preemption budgets.

---

## 6. Success Criteria (Things Weft Catches That `-race` Does Not)

- **Deadlocks**: incorrect `Cond.Wait` usage (`if` instead of `for`), lost wakeups, circular waits.
- **Missed/Spurious Wakeups**: predicates not rechecked, signals missed.
- **Linearizability Violations**: lost updates, stale reads, reordering under contention.
- **Livelock & Starvation**: priority inversion, unfair wake policy.
- **Fairness Issues**: starvation of readers/writers or channels.
- **Lock-Free Bugs**: ABA, protocol mistakes, double-drops, wrong atomic ordering.
- **Timer Bugs**: leaked timers, wrong resets, ordering errors.
- **Lost Data**: dropped/duplicated items without races.

---

## 7. Deliverables

1. **`weft/` Package**
   - Deterministic scheduler (`detsched` build tag).
   - Wrappers for spawn, sleep/after, mutexes, conds, channels.
   - Non-deterministic stubs for `!detsched`.

2. **`wefttest/` Helpers**
   - `Explore(t, N, buildFn)` — multi-seed exploration.
   - `Replay(t, seed, buildFn)` — seed replay.
   - `ReplayChoices(t, choices, buildFn)` — explicit trace replay.

3. **`cmd/weftfix/` Codemod**
   - Safe, reversible rewrites (`--dry-run`, `--path`).
   - Replace `time.Sleep/After`, `sync.*`, `make(chan)`, simple `go f()`.

4. **Documentation**
   - Quickstart, migration checklist, examples of bugs caught.
   - Guide on writing tests with exploration + replay.
   - FAQ on limitations.

---

## 8. Sample User Workflow

1. Write code using `weft` primitives (`weft.Go`, `weft.Sleep`, `weft.Mutex`, etc.).
2. Run tests normally: `go test ./...` (uses stdlib).  
3. Run deterministically: `go test -tags=detsched ./...`.  
4. Explore: `wefttest.Explore(t, 200, func(s *weft.Scheduler) { ... })`.  
5. On failure: log shows `seed=XYZ`.  
6. Replay with `wefttest.Replay(t, XYZ, func(s *weft.Scheduler){...})`.  
7. Debug, fix, and keep replay test as regression.

---

## 9. Future Extensions

- **Linearizability Checker** for collections.  
- **Automated Shrinking** to minimal failing traces.  
- **Chaos/Fault Injection** at yield points.  
- **Fairness Policies**: round-robin, priority-aware.  
- **CI Sharding** across seeds.  
- **Trace Export & Visualization** (JSON → timeline).  
- **Property-Based Testing Integration** (seed drives both data + schedule).  
- **Deterministic `Select` helper**.  
- **Adapters** for popular frameworks (HTTP, gRPC).  

---

**End of Document**
