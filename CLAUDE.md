# Weft Project Context

## Project Overview
Weft is a deterministic concurrency testing framework for Go that reliably reproduces, explores, and shrinks concurrency bugs that the Go race detector cannot detect.

## Core Objectives
- Provide drop-in replacements for Go's concurrency primitives
- Enable deterministic, seed-driven scheduling for reproducible test runs
- Support exploration across multiple schedules to uncover rare bugs
- Maintain zero overhead in production through build tags

## Architecture

### Build Modes
- **Deterministic Mode** (`-tags=detsched`): Uses internal scheduler for controlled execution
- **Production Mode** (default): Transparent passthrough to standard library with zero overhead

### Key Components
1. **Core Package (`weft/`)**: Main API with deterministic/production implementations
2. **Test Helpers (`wefttest/`)**: Exploration and replay utilities
3. **Internal Scheduler (`internal/scheduler/`)**: Deterministic scheduling engine
4. **Codemod Tool (`cmd/weftfix/`)**: Automated code transformation tool

### Supported Primitives
- `weft.Go()` - Deterministic goroutine spawning
- `weft.Mutex` / `weft.RWMutex` - Mutual exclusion locks
- `weft.Cond` - Condition variables
- `weft.MakeChan[T]()` - Deterministic channels
- `weft.Sleep()` / `weft.After()` - Virtual time operations
- `weft.Clock` - Logical clock interface

## Development Guidelines

### Commit Message Convention
**IMPORTANT**: All commits MUST follow the Conventional Commits specification:

```
<type>(<scope>): <subject>

<body>

<footer>
```

#### Types
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, missing semicolons, etc.)
- `refactor`: Code refactoring without changing functionality
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `build`: Build system or dependency changes
- `ci`: CI/CD configuration changes
- `chore`: Other changes that don't modify src or test files

#### Examples
```
feat(scheduler): add deterministic task scheduling

Implement core scheduler with seed-based PRNG for reproducible
task selection across goroutines.

Closes #12
```

```
fix(mutex): prevent deadlock in recursive lock attempts

Add ownership tracking to detect and panic on recursive
mutex lock attempts from the same goroutine.
```

```
docs(readme): update installation instructions
```

### Testing Strategy
1. All new features must include comprehensive tests
2. Use `wefttest.Explore()` to test across multiple schedules
3. Document any failing seeds for regression testing
4. Ensure both deterministic and production modes are tested

### Implementation Plan
**IMPORTANT**: Follow the detailed implementation roadmap in `.claude/implementation_roadmap.md`

The project follows a 5-milestone development plan:
1. **Milestone 1**: Core Foundation (PRNG, scheduler, virtual time)
2. **Milestone 2**: Synchronization Primitives (mutex, channels, condition variables)
3. **Milestone 3**: Advanced Detection (deadlock detection, timeouts, tracing)
4. **Milestone 4**: Testing & Production Readiness (performance, documentation)
5. **Milestone 5**: Extended Features (shrinking, policies, weftfix tool)

**Current Status**: Milestone 1 - Task 1.1 (PRNG Implementation)

### Task Management Process
1. **Before starting work**: Review current milestone in `.claude/milestone_tracker.md`
2. **During development**: Follow detailed tasks in `.claude/task_list.md`
3. **After completing tasks**: Update milestone tracker with progress
4. **Update this section**: Reflect current milestone/task status

### Implementation Status
- ✅ Basic project structure
- ✅ Build tag separation
- ✅ API stubs for all primitives
- ✅ Comprehensive examples with realistic production code
- 🔴 **Milestone 1**: Core Foundation (0/5 tasks complete)
  - ⏳ Task 1.1: PRNG Implementation (Next up)
  - ⏳ Task 1.2: Task Management System
  - ⏳ Task 1.3: Core Scheduler Engine
  - ⏳ Task 1.4: Virtual Time System
  - ⏳ Task 1.5: Fix Type Compatibility Issues

## Design Principles
1. **API Compatibility**: Mirror standard library APIs exactly
2. **Zero Production Overhead**: Build tags ensure no performance impact
3. **Deterministic by Design**: Same seed → same execution
4. **Progressive Enhancement**: Start simple, add features incrementally
5. **User-Friendly**: Clear error messages and debugging support

## Future Enhancements
- Linearizability checker for concurrent data structures
- Automated trace shrinking to minimal reproductions
- Chaos injection at yield points
- Multiple scheduling policies (round-robin, priority-based)
- Integration with property-based testing frameworks
- Trace visualization tools

## Bug Classes Detected
Weft aims to catch concurrency bugs that `-race` cannot:
- Deadlocks and circular waits
- Lost updates and stale reads
- Linearizability violations
- Starvation and livelock
- Protocol violations in lock-free algorithms
- Timing-dependent bugs

## References
- **Implementation Roadmap**: `.claude/implementation_roadmap.md` - High-level milestones and strategy
- **Task List**: `.claude/task_list.md` - Detailed technical tasks with acceptance criteria
- **Milestone Tracker**: `.claude/milestone_tracker.md` - Progress tracking and validation
- **Design Document**: `.claude/weft_design_doc.md` - Original technical specification
- **Examples**: `examples/` - Realistic production code examples
- **Issue Tracker**: https://github.com/mziter/weft/issues

## Contributing
1. Follow conventional commit messages
2. **Follow the implementation roadmap** - Check current milestone/task before starting work
3. **Update progress tracking** - Mark tasks complete in milestone tracker
4. Add tests for new features with >90% coverage target
5. Update documentation as needed
6. Ensure CI passes before merging
7. **Update CLAUDE.md status section** when milestones/tasks change