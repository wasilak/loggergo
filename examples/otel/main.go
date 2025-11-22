package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/wasilak/loggergo"
)

func main() {
	ctx := context.Background()

	// Configure logger with OpenTelemetry output
	config := loggergo.Config{
		Level:              slog.LevelInfo,
		Output:             loggergo.Types.OutputOtel,
		OtelLoggerName:     "github.com/wasilak/loggergo/examples/otel",
		OtelServiceName:    "loggergo-otel-example",
		OtelTracingEnabled: true,
	}

	ctx, logger, err := loggergo.Init(ctx, config)
	if err != nil {
		slog.Error("Failed to initialize logger", "error", err)
		os.Exit(1)
	}

	// Log some messages
	logger.Info("OTEL logger initialized successfully")
	logger.Info("This log will include trace_id and span_id when available")

	// Simulate some work
	logger.Debug("Processing request", "request_id", "req-123")
	time.Sleep(100 * time.Millisecond)

	logger.Info("Request processed successfully",
		"request_id", "req-123",
		"duration_ms", 100,
		"status", "success",
	)

	// Log an error
	logger.Error("Example error log",
		"error", "something went wrong",
		"code", 500,
	)

	logger.Info("OTEL example completed")

	// Note: In a real application, you would:
	// 1. Configure OTEL_EXPORTER_OTLP_ENDPOINT environment variable
	// 2. Set up proper trace context propagation
	// 3. Use context with trace information from incoming requests
	// 4. Ensure proper shutdown to flush logs
}
