// Package loggergo provides a lightweight, customizable logging library for Go applications.
//
// LoggerGo is built on top of Go's standard log/slog package and provides:
//   - Multiple output modes (console, OpenTelemetry, fanout)
//   - Multiple log formats (JSON, text, OTEL)
//   - Development mode with different flavors (tint, slogor, devslog)
//   - Context-aware logging with automatic value extraction
//   - OpenTelemetry integration with trace/span ID injection
//   - Dynamic log level adjustment
//   - Thread-safe configuration management
//   - Comprehensive error handling with graceful degradation
//
// Quick Start:
//
//	import (
//	    "context"
//	    "log/slog"
//	    "github.com/wasilak/loggergo"
//	)
//
//	func main() {
//	    ctx := context.Background()
//	    config := loggergo.Config{
//	        Level:  slog.LevelInfo,
//	        Format: loggergo.LogFormatJSON,
//	        Output: loggergo.OutputConsole,
//	    }
//	    ctx, logger, err := loggergo.Init(ctx, config)
//	    if err != nil {
//	        panic(err)
//	    }
//	    logger.Info("Hello, LoggerGo!")
//	}
//
// For more examples, see the examples/ directory in the repository.
package loggergo

import (
	"context"
	"fmt"
	"os"

	"log/slog"

	slogmulti "github.com/samber/slog-multi"

	"github.com/wasilak/loggergo/lib"
	"github.com/wasilak/loggergo/lib/modes"
	"github.com/wasilak/loggergo/lib/types"
)

var logLevel = new(slog.LevelVar)

// Config represents the configuration options for the logger.
// It is an alias for types.Config and is exported for external usage.
//
// See types.Config for detailed field documentation and usage examples.
type Config = types.Config

// Init initializes a logger with the provided configuration and additional attributes.
//
// It validates the configuration, creates the appropriate handler based on the output mode,
// and returns the updated context, configured logger, and any error encountered.
//
// The logger is automatically set as slog.Default() if Config.SetAsDefault is true (default).
//
// Parameters:
//   - ctx: The context to use for initialization and OTEL setup
//   - config: The logger configuration (see Config for details)
//   - additionalAttrs: Optional additional attributes to add to all log entries
//
// Returns:
//   - context.Context: Updated context (may include OTEL trace context)
//   - *slog.Logger: Configured logger instance
//   - error: InitError if initialization fails, nil on success
//
// Error Handling:
//
// Init never panics. All panics are recovered and returned as InitError with stage "panic_recovery".
// If OTEL mode fails, Init gracefully degrades to console mode with a warning to stderr.
//
// Thread Safety:
//
// Init is safe to call concurrently from multiple goroutines. However, only one initialization
// should typically be performed per application.
//
// Example:
//
//	config := loggergo.Config{
//	    Level:  slog.LevelInfo,
//	    Format: loggergo.LogFormatJSON,
//	    Output: loggergo.OutputConsole,
//	}
//	ctx, logger, err := loggergo.Init(ctx, config)
//	if err != nil {
//	    log.Fatalf("Failed to initialize logger: %v", err)
//	}
//	logger.Info("Logger initialized successfully")
func Init(ctx context.Context, config types.Config, additionalAttrs ...any) (retCtx context.Context, retLogger *slog.Logger, retErr error) {
	// Panic recovery to ensure Init never panics
	defer func() {
		if r := recover(); r != nil {
			cfg := lib.GetConfig()
			retCtx = ctx
			retLogger = nil
			retErr = &types.InitError{
				Stage:  "panic_recovery",
				Cause:  fmt.Errorf("panic during initialization: %v", r),
				Config: cfg,
			}
			// Log the panic to stderr for debugging
			fmt.Fprintf(os.Stderr, "PANIC recovered in Init(): %v\n", r)
		}
	}()

	var defaultHandler slog.Handler
	var err error

	lib.InitConfig()
	lib.MergeConfig(config)

	// Validate configuration before initialization
	cfg := lib.GetConfig()
	if err := cfg.Validate(); err != nil {
		return ctx, nil, &types.InitError{
			Stage:  "validation",
			Cause:  err,
			Config: cfg,
		}
	}

	logLevel.Set(lib.GetConfig().Level.Level())

	opts := slog.HandlerOptions{
		Level:     logLevel,
		AddSource: lib.GetConfig().Level == slog.LevelDebug,
	}

	switch lib.GetConfig().Output {
	case types.OutputConsole:
		defaultHandler, err = modes.ConsoleMode(opts)
		if err != nil {
			return ctx, nil, &types.InitError{
				Stage:  "handler_creation",
				Cause:  err,
				Config: cfg,
			}
		}
	case types.OutputOtel:
		defaultHandler, ctx, err = modes.OtelMode(ctx)
		if err != nil {
			// Graceful degradation: fall back to console mode on OTEL failure
			fmt.Fprintf(os.Stderr, "WARNING: OTEL initialization failed (%v), falling back to console mode\n", err)
			defaultHandler, err = modes.ConsoleMode(opts)
			if err != nil {
				return ctx, nil, &types.InitError{
					Stage:  "handler_creation",
					Cause:  fmt.Errorf("OTEL setup failed and console fallback also failed: %w", err),
					Config: cfg,
				}
			}
		}
	case types.OutputFanout:
		consoleModeHandler, err := modes.ConsoleMode(opts)
		if err != nil {
			return ctx, nil, &types.InitError{
				Stage:  "handler_creation",
				Cause:  err,
				Config: cfg,
			}
		}
		otelModeHandler, newCtx, err := modes.OtelMode(ctx)
		if err != nil {
			// Graceful degradation: use only console mode on OTEL failure in fanout
			fmt.Fprintf(os.Stderr, "WARNING: OTEL initialization failed in fanout mode (%v), using console mode only\n", err)
			defaultHandler = consoleModeHandler
		} else {
			ctx = newCtx
			defaultHandler = slogmulti.Fanout(
				consoleModeHandler,
				otelModeHandler,
			)
		}
	default:
		return ctx, nil, &types.InitError{
			Stage:  "validation",
			Cause:  fmt.Errorf("invalid mode: %s. Valid options: [loggergo.OutputConsole, loggergo.OutputOtel, loggergo.OutputFanout]", lib.GetConfig().Output),
			Config: cfg,
		}
	}

	// The code below is creating a new CustomContextAttributeHandler with the default handler and the context keys.
	defaultHandler = NewCustomContextAttributeHandler(defaultHandler, lib.GetConfig().ContextKeys, lib.GetConfig().ContextKeysDefault)

	logger := slog.New(defaultHandler)

	// The code below is iterating over the `additionalAttrs` slice and calling the `With` method on the default logger for
	// each element in the slice.
	for _, v := range additionalAttrs {
		logger.With(v)
	}

	if lib.GetConfig().SetAsDefault {
		// The code `slog.SetDefault(logger)` is setting the default logger to the newly created logger.
		slog.SetDefault(logger)
	}

	return ctx, logger, nil
}

// GetLogLevelAccessor returns the log level accessor for dynamic level changes.
//
// Thread-Safety Guarantees:
// This function is safe for concurrent use. The returned *slog.LevelVar uses
// atomic operations internally, making it safe to call Set() and Level() methods
// from multiple goroutines simultaneously without additional synchronization.
//
// Example usage:
//
//	levelVar := loggergo.GetLogLevelAccessor()
//	levelVar.Set(slog.LevelDebug) // Thread-safe level change
func GetLogLevelAccessor() *slog.LevelVar {
	return logLevel
}

// GetConfig returns the current logger configuration.
// This function is thread-safe and returns a copy of the configuration.
func GetConfig() types.Config {
	return lib.GetConfig()
}

// Shutdown performs cleanup of all registered resources.
//
// This function should be called when the application is shutting down to ensure
// proper cleanup of resources like OTEL providers, file handles, etc.
//
// Cleanup functions are called in reverse order of registration (LIFO).
// If any cleanup function returns an error, Shutdown continues with remaining
// cleanup functions and returns a combined error at the end.
//
// Thread Safety:
//
// Shutdown is safe to call concurrently, but should typically only be called once
// during application shutdown.
//
// Example:
//
//	ctx, logger, err := loggergo.Init(ctx, config)
//	if err != nil {
//	    panic(err)
//	}
//	defer loggergo.Shutdown()
//	
//	// Use logger...
func Shutdown() error {
	return lib.Shutdown()
}
