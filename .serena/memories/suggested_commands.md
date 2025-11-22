# Suggested Commands for LoggerGo Development

## Build & Verification
```bash
# Build all packages
go build ./...

# Run static analysis
go vet ./...

# Format code
go fmt ./...
```

## Testing
```bash
# Run all tests
go test ./...

# Run tests with race detector
go test -race ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestName ./...
```

## Dependencies
```bash
# Download dependencies
go mod download

# Tidy dependencies
go mod tidy

# Verify dependencies
go mod verify
```

## Examples
```bash
# Run simple example
cd examples/simple && go run main.go
```

## Task Completion Checklist
After implementing any task:
1. Run `go build ./...` - must succeed without errors/warnings
2. Run `go vet ./...` - must pass without issues
3. Run `go test ./...` - all tests must pass
4. Run `go test -race ./...` - no race conditions
5. Verify examples still work
6. Commit changes with descriptive message