---
title: Naming Conventions and Brand Agnosticism
inclusion: always
---

# Naming Conventions and Brand Agnosticism

## Brand-Agnostic Code

Whatever project name we settle on, code MUST be written in a way that is **name-agnostic**.

### ✅ Allowed Uses of Project Name

- **Documentation**: README, specs, comments, user-facing docs
- **Package name**: `package loggergo`, `module github.com/wasilak/loggergo`
- **Repository name**: GitHub repo name
- **Import path**: `github.com/wasilak/loggergo`

### ❌ NOT Allowed in Code

- **Function names**: ❌ `func LoggergoInit()` → ✅ `func Init()`
- **Struct names**: ❌ `type LoggergoConfig struct` → ✅ `type Config struct`
- **Variable names**: ❌ `loggergoConfig` → ✅ `config`
- **Interface names**: ❌ `LoggergoHandler` → ✅ `Handler`
- **Method names**: ❌ `InitLoggergo()` → ✅ `Init()`
- **Constants**: ❌ `LoggergoVersion` → ✅ `Version`

## Go Naming Conventions

Follow standard Go naming conventions:

### Package Names
- Short, lowercase, single word
- No underscores or mixedCaps
- Examples: `types`, `modes`, `outputs`

### Exported Names
- Use MixedCaps (PascalCase) for exported names
- Examples: `Config`, `Handler`, `LogFormat`

### Unexported Names
- Use mixedCaps (camelCase) for unexported names
- Examples: `initConfig`, `mergeConfig`, `setupHandler`

### Interface Names
- Single-method interfaces: verb + "er" suffix
- Examples: `Reader`, `Writer`, `Handler`
- Multi-method interfaces: descriptive noun
- Examples: `Logger`, `Formatter`, `Validator`

### Avoid Stuttering
- ❌ `types.ConfigType` → ✅ `types.Config`
- ❌ `handler.HandlerOptions` → ✅ `handler.Options`
- ❌ `logger.LoggerInit()` → ✅ `logger.Init()`

## Examples

### ✅ Good - Brand Agnostic

```go
package loggergo

type Config struct {
    Level  slog.Leveler
    Format LogFormat
}

func Init(ctx context.Context, config Config) (context.Context, *slog.Logger, error) {
    // Initialize logger
    return ctx, logger, nil
}

func (c *Config) Validate() error {
    // Validate configuration
    return nil
}
```

### ❌ Bad - Brand Specific

```go
package loggergo

type LoggergoConfig struct {
    Level  slog.Leveler
    Format LogFormat
}

func InitLoggergo(ctx context.Context, config LoggergoConfig) (context.Context, *slog.Logger, error) {
    // Initialize logger
    return ctx, logger, nil
}

func (c *LoggergoConfig) ValidateLoggergo() error {
    // Validate configuration
    return nil
}
```

## Why This Matters

1. **Reusability**: Code can be forked/reused without renaming everything
2. **Clarity**: Shorter names are easier to read and understand
3. **Go Idioms**: Follows standard Go conventions
4. **Maintainability**: Less coupling to brand name
5. **Professionalism**: Shows understanding of proper Go design

## Documentation is Different

In documentation, user-facing messages, and comments, using the project name is fine:

```go
// Config represents the configuration options for the LoggerGo logger.
type Config struct {
    // ...
}

// Package loggergo provides functionality for configuring and setting up 
// different logging modes in Go applications.
package loggergo
```

The key is: **code structure and naming should be generic, documentation can be branded**.
