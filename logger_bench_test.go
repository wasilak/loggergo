package loggergo

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/wasilak/loggergo/lib/types"
)

// BenchmarkLogger_Info benchmarks basic Info logging calls
func BenchmarkLogger_Info(b *testing.B) {
	ctx := context.Background()
	var buf bytes.Buffer

	config := types.Config{
		Level:        slog.LevelInfo,
		Format:       types.LogFormatJSON,
		Output:       types.OutputConsole,
		OutputStream: &buf,
		SetAsDefault: false,
	}

	_, logger, err := Init(ctx, config)
	if err != nil {
		b.Fatalf("Failed to initialize logger: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("test message", "key", "value", "iteration", i)
	}
}

// BenchmarkLogger_Debug benchmarks Debug logging calls
func BenchmarkLogger_Debug(b *testing.B) {
	ctx := context.Background()
	var buf bytes.Buffer

	config := types.Config{
		Level:        slog.LevelDebug,
		Format:       types.LogFormatJSON,
		Output:       types.OutputConsole,
		OutputStream: &buf,
		SetAsDefault: false,
	}

	_, logger, err := Init(ctx, config)
	if err != nil {
		b.Fatalf("Failed to initialize logger: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug("test message", "key", "value", "iteration", i)
	}
}

// BenchmarkLogger_Error benchmarks Error logging calls
func BenchmarkLogger_Error(b *testing.B) {
	ctx := context.Background()
	var buf bytes.Buffer

	config := types.Config{
		Level:        slog.LevelInfo,
		Format:       types.LogFormatJSON,
		Output:       types.OutputConsole,
		OutputStream: &buf,
		SetAsDefault: false,
	}

	_, logger, err := Init(ctx, config)
	if err != nil {
		b.Fatalf("Failed to initialize logger: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Error("test message", "key", "value", "iteration", i)
	}
}

// BenchmarkLogger_ConsoleJSON benchmarks logging with console output in JSON format
func BenchmarkLogger_ConsoleJSON(b *testing.B) {
	ctx := context.Background()
	var buf bytes.Buffer

	config := types.Config{
		Level:        slog.LevelInfo,
		Format:       types.LogFormatJSON,
		Output:       types.OutputConsole,
		OutputStream: &buf,
		SetAsDefault: false,
	}

	_, logger, err := Init(ctx, config)
	if err != nil {
		b.Fatalf("Failed to initialize logger: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("test message", "key", "value")
	}
}

// BenchmarkLogger_ConsoleText benchmarks logging with console output in text format
func BenchmarkLogger_ConsoleText(b *testing.B) {
	ctx := context.Background()
	var buf bytes.Buffer

	config := types.Config{
		Level:        slog.LevelInfo,
		Format:       types.LogFormatText,
		Output:       types.OutputConsole,
		OutputStream: &buf,
		SetAsDefault: false,
	}

	_, logger, err := Init(ctx, config)
	if err != nil {
		b.Fatalf("Failed to initialize logger: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("test message", "key", "value")
	}
}

// BenchmarkLogger_DevMode benchmarks logging with development mode enabled
func BenchmarkLogger_DevMode(b *testing.B) {
	ctx := context.Background()
	var buf bytes.Buffer

	config := types.Config{
		Level:        slog.LevelInfo,
		Format:       types.LogFormatText,
		Output:       types.OutputConsole,
		OutputStream: &buf,
		SetAsDefault: false,
		DevMode:      true,
		DevFlavor:    types.DevFlavorTint,
	}

	_, logger, err := Init(ctx, config)
	if err != nil {
		b.Fatalf("Failed to initialize logger: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("test message", "key", "value")
	}
}

// BenchmarkLogger_WithContextExtraction benchmarks logging with context value extraction
func BenchmarkLogger_WithContextExtraction(b *testing.B) {
	var buf bytes.Buffer

	config := types.Config{
		Level:              slog.LevelInfo,
		Format:             types.LogFormatJSON,
		Output:             types.OutputConsole,
		OutputStream:       &buf,
		SetAsDefault:       false,
		ContextKeys:        []interface{}{"request_id", "user_id"},
		ContextKeysDefault: "unknown",
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", "req-123")
	ctx = context.WithValue(ctx, "user_id", "user-456")

	_, logger, err := Init(ctx, config)
	if err != nil {
		b.Fatalf("Failed to initialize logger: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.InfoContext(ctx, "test message", "key", "value")
	}
}

// BenchmarkLogger_WithoutContextExtraction benchmarks logging without context extraction for comparison
func BenchmarkLogger_WithoutContextExtraction(b *testing.B) {
	ctx := context.Background()
	var buf bytes.Buffer

	config := types.Config{
		Level:        slog.LevelInfo,
		Format:       types.LogFormatJSON,
		Output:       types.OutputConsole,
		OutputStream: &buf,
		SetAsDefault: false,
	}

	_, logger, err := Init(ctx, config)
	if err != nil {
		b.Fatalf("Failed to initialize logger: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.InfoContext(ctx, "test message", "key", "value")
	}
}

// BenchmarkLogger_FilteredDebug benchmarks Debug calls when level is set to Info (filtered out)
func BenchmarkLogger_FilteredDebug(b *testing.B) {
	ctx := context.Background()
	var buf bytes.Buffer

	config := types.Config{
		Level:        slog.LevelInfo, // Debug logs will be filtered
		Format:       types.LogFormatJSON,
		Output:       types.OutputConsole,
		OutputStream: &buf,
		SetAsDefault: false,
	}

	_, logger, err := Init(ctx, config)
	if err != nil {
		b.Fatalf("Failed to initialize logger: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug("test message", "key", "value", "iteration", i)
	}
}

// BenchmarkLogger_NonFilteredDebug benchmarks Debug calls when level is set to Debug (not filtered)
func BenchmarkLogger_NonFilteredDebug(b *testing.B) {
	ctx := context.Background()
	var buf bytes.Buffer

	config := types.Config{
		Level:        slog.LevelDebug, // Debug logs will be processed
		Format:       types.LogFormatJSON,
		Output:       types.OutputConsole,
		OutputStream: &buf,
		SetAsDefault: false,
	}

	_, logger, err := Init(ctx, config)
	if err != nil {
		b.Fatalf("Failed to initialize logger: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug("test message", "key", "value", "iteration", i)
	}
}

// BenchmarkLogger_FilteredInfo benchmarks Info calls when level is set to Error (filtered out)
func BenchmarkLogger_FilteredInfo(b *testing.B) {
	ctx := context.Background()
	var buf bytes.Buffer

	config := types.Config{
		Level:        slog.LevelError, // Info logs will be filtered
		Format:       types.LogFormatJSON,
		Output:       types.OutputConsole,
		OutputStream: &buf,
		SetAsDefault: false,
	}

	_, logger, err := Init(ctx, config)
	if err != nil {
		b.Fatalf("Failed to initialize logger: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("test message", "key", "value", "iteration", i)
	}
}

// BenchmarkLogger_NonFilteredInfo benchmarks Info calls when level is set to Info (not filtered)
func BenchmarkLogger_NonFilteredInfo(b *testing.B) {
	ctx := context.Background()
	var buf bytes.Buffer

	config := types.Config{
		Level:        slog.LevelInfo, // Info logs will be processed
		Format:       types.LogFormatJSON,
		Output:       types.OutputConsole,
		OutputStream: &buf,
		SetAsDefault: false,
	}

	_, logger, err := Init(ctx, config)
	if err != nil {
		b.Fatalf("Failed to initialize logger: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("test message", "key", "value", "iteration", i)
	}
}

// BenchmarkLogger_FilteredWithExpensiveOperation benchmarks filtered logs with expensive operations
// This demonstrates that expensive operations in log arguments are still evaluated even when filtered
func BenchmarkLogger_FilteredWithExpensiveOperation(b *testing.B) {
	ctx := context.Background()
	var buf bytes.Buffer

	config := types.Config{
		Level:        slog.LevelError, // Debug logs will be filtered
		Format:       types.LogFormatJSON,
		Output:       types.OutputConsole,
		OutputStream: &buf,
		SetAsDefault: false,
	}

	_, logger, err := Init(ctx, config)
	if err != nil {
		b.Fatalf("Failed to initialize logger: %v", err)
	}

	expensiveFunc := func() string {
		// Simulate expensive operation
		result := ""
		for i := 0; i < 100; i++ {
			result += "x"
		}
		return result
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug("test message", "expensive", expensiveFunc())
	}
}

// BenchmarkLogger_FilteredWithEnabledCheck benchmarks filtered logs with Enabled() check
// This demonstrates the proper way to avoid expensive operations for filtered logs
func BenchmarkLogger_FilteredWithEnabledCheck(b *testing.B) {
	ctx := context.Background()
	var buf bytes.Buffer

	config := types.Config{
		Level:        slog.LevelError, // Debug logs will be filtered
		Format:       types.LogFormatJSON,
		Output:       types.OutputConsole,
		OutputStream: &buf,
		SetAsDefault: false,
	}

	_, logger, err := Init(ctx, config)
	if err != nil {
		b.Fatalf("Failed to initialize logger: %v", err)
	}

	expensiveFunc := func() string {
		// Simulate expensive operation
		result := ""
		for i := 0; i < 100; i++ {
			result += "x"
		}
		return result
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if logger.Enabled(ctx, slog.LevelDebug) {
			logger.Debug("test message", "expensive", expensiveFunc())
		}
	}
}

// BenchmarkLogger_ConcurrentLogging benchmarks concurrent logging from multiple goroutines
func BenchmarkLogger_ConcurrentLogging(b *testing.B) {
	ctx := context.Background()
	var buf bytes.Buffer

	config := types.Config{
		Level:        slog.LevelInfo,
		Format:       types.LogFormatJSON,
		Output:       types.OutputConsole,
		OutputStream: &buf,
		SetAsDefault: false,
	}

	_, logger, err := Init(ctx, config)
	if err != nil {
		b.Fatalf("Failed to initialize logger: %v", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			logger.Info("test message", "key", "value", "iteration", i)
			i++
		}
	})
}

// BenchmarkLogger_ConcurrentLoggingWithContext benchmarks concurrent logging with context extraction
func BenchmarkLogger_ConcurrentLoggingWithContext(b *testing.B) {
	var buf bytes.Buffer

	config := types.Config{
		Level:              slog.LevelInfo,
		Format:             types.LogFormatJSON,
		Output:             types.OutputConsole,
		OutputStream:       &buf,
		SetAsDefault:       false,
		ContextKeys:        []interface{}{"request_id"},
		ContextKeysDefault: "unknown",
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", "req-123")

	_, logger, err := Init(ctx, config)
	if err != nil {
		b.Fatalf("Failed to initialize logger: %v", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			logger.InfoContext(ctx, "test message", "key", "value", "iteration", i)
			i++
		}
	})
}

// BenchmarkLogger_ConcurrentLevelChanges benchmarks concurrent logging with dynamic level changes
func BenchmarkLogger_ConcurrentLevelChanges(b *testing.B) {
	ctx := context.Background()
	var buf bytes.Buffer

	config := types.Config{
		Level:        slog.LevelInfo,
		Format:       types.LogFormatJSON,
		Output:       types.OutputConsole,
		OutputStream: &buf,
		SetAsDefault: false,
	}

	_, logger, err := Init(ctx, config)
	if err != nil {
		b.Fatalf("Failed to initialize logger: %v", err)
	}

	levelAccessor := GetLogLevelAccessor()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// Alternate between logging and changing levels
			if i%10 == 0 {
				if i%20 == 0 {
					levelAccessor.Set(slog.LevelDebug)
				} else {
					levelAccessor.Set(slog.LevelInfo)
				}
			}
			logger.Info("test message", "key", "value", "iteration", i)
			i++
		}
	})
}

// BenchmarkLogger_ConcurrentMixedLevels benchmarks concurrent logging with different log levels
func BenchmarkLogger_ConcurrentMixedLevels(b *testing.B) {
	ctx := context.Background()
	var buf bytes.Buffer

	config := types.Config{
		Level:        slog.LevelDebug,
		Format:       types.LogFormatJSON,
		Output:       types.OutputConsole,
		OutputStream: &buf,
		SetAsDefault: false,
	}

	_, logger, err := Init(ctx, config)
	if err != nil {
		b.Fatalf("Failed to initialize logger: %v", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			switch i % 4 {
			case 0:
				logger.Debug("debug message", "iteration", i)
			case 1:
				logger.Info("info message", "iteration", i)
			case 2:
				logger.Warn("warn message", "iteration", i)
			case 3:
				logger.Error("error message", "iteration", i)
			}
			i++
		}
	})
}

// BenchmarkLogger_ConcurrentWithHighContention benchmarks concurrent logging with high contention
// This uses a smaller buffer to increase lock contention
func BenchmarkLogger_ConcurrentWithHighContention(b *testing.B) {
	ctx := context.Background()
	var buf bytes.Buffer

	config := types.Config{
		Level:        slog.LevelInfo,
		Format:       types.LogFormatJSON,
		Output:       types.OutputConsole,
		OutputStream: &buf,
		SetAsDefault: false,
	}

	_, logger, err := Init(ctx, config)
	if err != nil {
		b.Fatalf("Failed to initialize logger: %v", err)
	}

	b.ResetTimer()
	b.SetParallelism(100) // Increase parallelism to create more contention
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			logger.Info("test message", "key", "value", "iteration", i)
			i++
		}
	})
}

// BenchmarkLogger_SingleGoroutine benchmarks single-threaded logging for comparison
func BenchmarkLogger_SingleGoroutine(b *testing.B) {
	ctx := context.Background()
	var buf bytes.Buffer

	config := types.Config{
		Level:        slog.LevelInfo,
		Format:       types.LogFormatJSON,
		Output:       types.OutputConsole,
		OutputStream: &buf,
		SetAsDefault: false,
	}

	_, logger, err := Init(ctx, config)
	if err != nil {
		b.Fatalf("Failed to initialize logger: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("test message", "key", "value", "iteration", i)
	}
}

// BenchmarkLogger_LockContention measures lock contention specifically in concurrent logging
// This benchmark increases parallelism to create more contention on shared resources
func BenchmarkLogger_LockContention(b *testing.B) {
	ctx := context.Background()
	var buf bytes.Buffer

	config := types.Config{
		Level:        slog.LevelInfo,
		Format:       types.LogFormatJSON,
		Output:       types.OutputConsole,
		OutputStream: &buf,
		SetAsDefault: false,
	}

	_, logger, err := Init(ctx, config)
	if err != nil {
		b.Fatalf("Failed to initialize logger: %v", err)
	}

	// Use higher parallelism to increase lock contention
	b.ResetTimer()
	b.SetParallelism(200) // Higher parallelism to stress locks
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// Perform logging that will require synchronization
			logger.Info("contention test message", "goroutine_id", i%100, "iteration", i)
			i++
		}
	})
}

// BenchmarkLogger_LockContentionWithContext measures lock contention with context extraction
// This benchmark tests the additional overhead of context value extraction under high contention
func BenchmarkLogger_LockContentionWithContext(b *testing.B) {
	var buf bytes.Buffer

	config := types.Config{
		Level:              slog.LevelInfo,
		Format:             types.LogFormatJSON,
		Output:             types.OutputConsole,
		OutputStream:       &buf,
		SetAsDefault:       false,
		ContextKeys:        []interface{}{"request_id", "user_id", "session_id"},
		ContextKeysDefault: "unknown",
	}

	// Create context with values
	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", "req-123")
	ctx = context.WithValue(ctx, "user_id", "user-456")
	ctx = context.WithValue(ctx, "session_id", "sess-789")

	_, logger, err := Init(ctx, config)
	if err != nil {
		b.Fatalf("Failed to initialize logger: %v", err)
	}

	// Use higher parallelism to increase lock contention
	b.ResetTimer()
	b.SetParallelism(150) // High parallelism for context extraction contention
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// Perform logging with context that will require value extraction and synchronization
			logger.InfoContext(ctx, "context contention test", "iteration", i)
			i++
		}
	})
}
