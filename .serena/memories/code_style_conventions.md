# Code Style and Conventions

## Brand-Agnostic Naming
- **DO NOT** use project name in code (functions, structs, variables)
- ✅ `func Init()` not ❌ `func LoggergoInit()`
- ✅ `type Config struct` not ❌ `type LoggergoConfig struct`
- Project name OK in: documentation, package name, import path

## Go Naming Conventions
- Package names: short, lowercase, single word (no underscores)
- Exported names: MixedCaps (PascalCase)
- Unexported names: mixedCaps (camelCase)
- Avoid stuttering: `types.Config` not `types.ConfigType`

## Error Handling
- Always check errors, never ignore
- Wrap errors with context: `fmt.Errorf("failed to X: %w", err)`
- Return errors, don't panic (except truly exceptional cases)

## Concurrency
- Use `sync.RWMutex` for read-heavy workloads
- Always use `defer` to unlock mutexes
- Pass `context.Context` as first parameter
- Don't store context in structs

## Testing
- Table-driven tests for multiple cases
- Test file naming: `*_test.go`
- Property-based tests: `*_property_test.go`
- Use `t.Helper()` in test helper functions
- Maintain >80% coverage for core functionality

## Documentation
- All exported symbols need godoc comments
- Start comments with the name being documented
- Use complete sentences
- Include usage examples in comments