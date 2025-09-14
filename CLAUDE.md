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

### Implementation Status
- ✅ Basic project structure
- ✅ Build tag separation
- ✅ API stubs for all primitives
- ⚠️ Scheduler implementation (stub only)
- ⚠️ Deterministic PRNG (interface only)
- ⚠️ Virtual time implementation
- ⚠️ Deadlock detection
- ⚠️ Trace recording and replay
- ⚠️ Shrinking algorithm
- ⚠️ Codemod tool implementation

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
- Design Document: `.claude/weft_design_doc.md`
- Examples: `examples/`
- Issue Tracker: https://github.com/mziter/weft/issues

## Contributing
1. Follow conventional commit messages
2. Add tests for new features
3. Update documentation as needed
4. Ensure CI passes before merging