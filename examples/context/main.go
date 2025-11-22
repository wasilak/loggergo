package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/wasilak/loggergo"
)

// Define custom context keys
type contextKey string

const (
	requestIDKey contextKey = "request_id"
	userIDKey    contextKey = "user_id"
	sessionIDKey contextKey = "session_id"
)

func main() {
	ctx := context.Background()

	// Configure logger with context extraction
	config := loggergo.Config{
		Level:  slog.LevelInfo,
		Format: loggergo.Types.LogFormatJSON,
		Output: loggergo.Types.OutputConsole,
		// Specify which context keys to extract
		ContextKeys: []interface{}{requestIDKey, userIDKey, sessionIDKey},
		// Default value for missing keys
		ContextKeysDefault: "unknown",
	}

	ctx, logger, err := loggergo.Init(ctx, config)
	if err != nil {
		slog.Error("Failed to initialize logger", "error", err)
		os.Exit(1)
	}

	logger.Info("Logger initialized with context extraction")

	// Example 1: Log without context values
	logger.InfoContext(ctx, "No context values set yet")
	// Output will include: request_id: "unknown", user_id: "unknown", session_id: "unknown"

	// Example 2: Add request ID to context
	ctx = context.WithValue(ctx, requestIDKey, "req-12345")
	logger.InfoContext(ctx, "Processing request")
	// Output will include: request_id: "req-12345", user_id: "unknown", session_id: "unknown"

	// Example 3: Add user ID to context
	ctx = context.WithValue(ctx, userIDKey, "user-789")
	logger.InfoContext(ctx, "User authenticated")
	// Output will include: request_id: "req-12345", user_id: "user-789", session_id: "unknown"

	// Example 4: Add session ID to context
	ctx = context.WithValue(ctx, sessionIDKey, "sess-abc")
	logger.InfoContext(ctx, "Session created")
	// Output will include: request_id: "req-12345", user_id: "user-789", session_id: "sess-abc"

	// Example 5: Simulate a function call with context
	processUserRequest(ctx, logger)

	logger.Info("Context example completed")
}

// processUserRequest demonstrates how context values propagate through function calls
func processUserRequest(ctx context.Context, logger *slog.Logger) {
	// All context values are automatically included in logs
	logger.InfoContext(ctx, "Processing user request in nested function")

	// Simulate some work
	logger.DebugContext(ctx, "Validating user permissions")
	logger.InfoContext(ctx, "User request processed successfully", "result", "ok")
}
