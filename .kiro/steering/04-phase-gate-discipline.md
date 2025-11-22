---
title: Incremental Development Discipline
inclusion: always
---

# Incremental Development Discipline

## Library Development Principles

This library follows **strict incremental development** with mandatory quality gates. This discipline ensures each improvement is fully tested and working before moving to the next.

## Development Rules

### 1. Sequential Task Execution
- **NEVER** start Task N+1 until Task N is complete and verified
- **NEVER** implement features from future tasks
- **NEVER** add "nice to have" features not in current task spec

### 2. Task Completion Requirements

Every task must meet:

1. ✅ **Implementation complete** - All code changes made as specified
2. ✅ **Build succeeds** - `go build ./...` completes without errors or warnings
3. ✅ **Vet passes** - `go vet ./...` reports no issues
4. ✅ **Tests pass** - `go test ./...` passes all tests
5. ✅ **Race detector clean** - `go test -race ./...` finds no races
6. ✅ **Examples work** - Example code runs without errors
7. ✅ **Documentation updated** - Godoc comments and README reflect changes

### 3. Task Failure Protocol

If ANY requirement fails:

1. **STOP** - Do not proceed to next task
2. **FIX** - Address the failing requirement
3. **RE-TEST** - Verify the fix works
4. **DOCUMENT** - Note what failed and how it was fixed
5. **RETRY** - Re-run full verification checklist

### 4. No Scope Creep

During task implementation:

- ❌ "While I'm here, let me add..."
- ❌ "This would be better if..."
- ❌ "Let me refactor this to be more flexible..."
- ✅ "Does this task spec require this? No? Then don't do it."

## Quality Gates

### Build Quality
```bash
# Must pass without errors or warnings
go build ./...
go vet ./...
```

### Test Quality
```bash
# All tests must pass
go test ./...

# No race conditions
go test -race ./...

# Coverage should be maintained or improved
go test -cover ./...
```

### Example Quality
```bash
# Examples must run successfully
cd examples/simple
go run main.go
```

## Testing Discipline

### Test-After-Implementation
- Implement the feature first
- Write tests to verify it works
- Fix any issues found by tests
- Ensure all tests pass before marking complete

### Test Coverage
- Maintain >80% coverage for core functionality
- Test happy paths and error cases
- Include edge cases (nil, empty, invalid inputs)
- Use property-based tests for universal properties

### Test Organization
```go
// Table-driven tests for multiple cases
func TestConfig_Validate(t *testing.T) {
    tests := []struct {
        name    string
        config  Config
        wantErr bool
    }{
        {"valid config", validConfig, false},
        {"invalid level", invalidConfig, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.config.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Git Discipline

### Commit After Each Task
```bash
# After task completion and user approval
git add <relevant files>
git commit -m "feat: add config validation

Task: 2.1 Implement Config.Validate() method"
```

### Commit Message Format
- Type: `feat:`, `fix:`, `refactor:`, `test:`, `docs:`
- Brief description (50 chars or less)
- Blank line
- Task reference

## Communication

### When Completing Task
Agent should say:
```
Task N implementation complete.

Verification:
✅ Build passes
✅ Vet passes  
✅ Tests pass
✅ Race detector clean
✅ Examples work

Ready for your review.
```

### When Task Passes Review
User approves, then agent:
```
Task N complete ✅

Committing changes...
Next: Task N+1 - [Description]

Ready to proceed?
```

### When Verification Fails
```
Task N verification FAILED ❌

Failed check: [description]
Issue: [what went wrong]

Fixing now...
```

## Discipline Checklist

Before starting any task, ask:

- [ ] Is this task in the current spec?
- [ ] Has the previous task been completed and committed?
- [ ] Am I adding features not in the spec?
- [ ] Am I over-engineering the solution?
- [ ] Is this the simplest approach that works?

If any answer is wrong, **STOP** and reconsider.

## Library-Specific Considerations

### API Stability
- Don't break existing public APIs
- Use deprecation for API evolution
- Add new functions rather than changing signatures
- Maintain backward compatibility

### Documentation
- Every exported symbol needs godoc
- README must have clear examples
- Document all configuration options
- Include troubleshooting guidance

### Dependencies
- Minimize external dependencies
- Pin dependency versions
- Prefer standard library when possible
- Document why each dependency is needed

## Remember

**The goal is a reliable, well-tested library that users can depend on.**

Build incrementally. Test thoroughly. Document clearly. Repeat.
