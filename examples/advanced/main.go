package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/wasilak/loggergo"
)

type contextKey string

const (
	requestIDKey contextKey = "request_id"
	traceIDKey   contextKey = "trace_id"
)

func main() {
	ctx := context.Background()

	// Advanced configuration with all options
	config := loggergo.Config{
		// Core settings
		Level:        slog.LevelDebug,
		Format:       loggergo.Types.LogFormatJSON,
		Output:       loggergo.Types.OutputFanout, // Send to both console and OTEL
		OutputStream: os.Stdout,
		SetAsDefault: true, // Register as slog.Default()

		// Development mode settings
		DevMode:   true,
		DevFlavor: loggergo.Types.DevFlavorTint, // Pretty colored output

		// OpenTelemetry settings
		OtelTracingEnabled: true,
		OtelLoggerName:     "github.com/wasilak/loggergo/examples/advanced",
		OtelServiceName:    "loggergo-advanced-example",

		// Context extraction
		ContextKeys:        []interface{}{requestIDKey, traceIDKey},
		ContextKeysDefault: "n/a",
	}

	ctx, logger, err := loggergo.Init(ctx, config)
	if err != nil {
		slog.Error("Failed to initialize logger", "error", err)
		os.Exit(1)
	}

	logger.Info("Advanced logger initialized",
		"output_mode", "fanout",
		"dev_mode", true,
		"otel_enabled", true,
	)

	// Demonstrate dynamic level changes
	demonstrateDynamicLevels(ctx, logger)

	// Demonstrate context-aware logging
	demonstrateContextLogging(ctx, logger)

	// Demonstrate structured logging
	demonstrateStructuredLogging(ctx, logger)

	// Demonstrate error handling
	demonstrateErrorHandling(ctx, logger)

	logger.Info("Advanced example completed")
}

func demonstrateDynamicLevels(ctx context.Context, logger *slog.Logger) {
	logger.Info("=== Dynamic Level Changes ===")

	// Get the level accessor
	levelVar := loggergo.GetLogLevelAccessor()

	// Start with Debug level
	levelVar.Set(slog.LevelDebug)
	logger.Debug("Debug message (visible)")
	logger.Info("Info message (visible)")

	// Change to Info level
	levelVar.Set(slog.LevelInfo)
	logger.Debug("Debug message (hidden)")
	logger.Info("Info message (visible)")

	// Change to Warn level
	levelVar.Set(slog.LevelWarn)
	logger.Info("Info message (hidden)")
	logger.Warn("Warn message (visible)")

	// Reset to Debug for rest of examples
	levelVar.Set(slog.LevelDebug)
}

func demonstrateContextLogging(ctx context.Context, logger *slog.Logger) {
	logger.Info("=== Context-Aware Logging ===")

	// Add request ID to context
	ctx = context.WithValue(ctx, requestIDKey, "req-advanced-001")
	logger.InfoContext(ctx, "Request received")

	// Add trace ID to context
	ctx = context.WithValue(ctx, traceIDKey, "trace-xyz-789")
	logger.InfoContext(ctx, "Processing with trace context")

	// Simulate nested function calls
	processRequest(ctx, logger)
}

func processRequest(ctx context.Context, logger *slog.Logger) {
	// Context values automatically included
	logger.InfoContext(ctx, "Processing in nested function")

	// Simulate some work
	time.Sleep(50 * time.Millisecond)

	logger.InfoContext(ctx, "Nested processing complete", "duration_ms", 50)
}

func demonstrateStructuredLogging(ctx context.Context, logger *slog.Logger) {
	logger.Info("=== Structured Logging ===")

	// Log with multiple attributes
	logger.Info("User action",
		"action", "login",
		"user_id", "user-123",
		"ip_address", "192.168.1.1",
		"timestamp", time.Now().Unix(),
		"success", true,
	)

	// Log with nested attributes using groups
	logger.LogAttrs(ctx, slog.LevelInfo, "Database query",
		slog.Group("query",
			slog.String("table", "users"),
			slog.String("operation", "SELECT"),
			slog.Int("rows", 42),
		),
		slog.Group("performance",
			slog.Duration("duration", 15*time.Millisecond),
			slog.Bool("cached", false),
		),
	)

	// Log with complex data structures
	logger.Info("API response",
		"endpoint", "/api/v1/users",
		"method", "GET",
		"status_code", 200,
		"response_time_ms", 125,
		"cache_hit", false,
	)
}

func demonstrateErrorHandling(ctx context.Context, logger *slog.Logger) {
	logger.Info("=== Error Handling ===")

	// Log errors with context
	err := simulateError()
	if err != nil {
		logger.Error("Operation failed",
			"error", err,
			"operation", "simulate_error",
			"retry_count", 3,
			"will_retry", false,
		)
	}

	// Log warnings
	logger.Warn("Resource usage high",
		"resource", "memory",
		"usage_percent", 85,
		"threshold_percent", 80,
	)

	// Log with stack trace information
	logger.Error("Critical error",
		"error", "database connection lost",
		"component", "database_pool",
		"connections_active", 0,
		"connections_max", 10,
	)
}

func simulateError() error {
	return os.ErrNotExist
}
