package loggergo

import (
	"context"
	"log/slog"
	"sync"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/wasilak/loggergo/lib/types"
)

// **Feature: logger-library-audit-improvements, Property 4: No panics on invalid state**
// For any sequence of operations including invalid inputs, the library should never panic
// **Validates: Requirements 3.2, 3.4**
func TestProperty_NoPanicsOnInvalidState(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("Init never panics with any config", prop.ForAll(
		func(config types.Config) bool {
			// This test verifies that Init never panics, even with invalid configs
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Init panicked with config %+v: %v", config, r)
				}
			}()

			ctx := context.Background()
			_, _, err := Init(ctx, config)
			
			// We don't care if it returns an error (that's expected for invalid configs)
			// We only care that it doesn't panic
			_ = err
			
			return true
		},
		genAnyConfig(),
	))

	properties.Property("Handle never panics with any context", prop.ForAll(
		func(ctxIsNil bool, keys []string) bool {
			// This test verifies that Handle never panics, even with nil contexts or invalid keys
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Handle panicked: %v", r)
				}
			}()

			// Create a simple handler for testing
			handler := NewCustomContextAttributeHandler(
				slog.NewJSONHandler(nil, nil),
				interfaceSlice(keys),
				"default",
			)

			var ctx context.Context
			if ctxIsNil {
				ctx = nil
			} else {
				ctx = context.Background()
				// Add some random values to context
				for _, key := range keys {
					ctx = context.WithValue(ctx, key, "value")
				}
			}

			record := slog.Record{}
			record.Level = slog.LevelInfo
			record.Message = "test"
			err := handler.Handle(ctx, record)
			
			// We don't care if it returns an error
			// We only care that it doesn't panic
			_ = err
			
			return true
		},
		gen.Bool(),
		gen.SliceOf(gen.AlphaString()),
	))

	properties.TestingRun(t)
}

// genAnyConfig generates any Config instance, including invalid ones
func genAnyConfig() gopter.Gen {
	return gopter.CombineGens(
		genAnyLevel(),
		genAnyFormat(),
		gen.Bool(),
		genAnyDevFlavor(),
		gen.Bool(),
		gen.AnyString(),
		genAnyOutput(),
		gen.AnyString(),
		gen.Bool(),
	).Map(func(values []interface{}) types.Config {
		var level slog.Leveler
		if values[0] != nil {
			level = values[0].(slog.Leveler)
		}
		
		return types.Config{
			Level:              level,
			Format:             values[1].(types.LogFormat),
			DevMode:            values[2].(bool),
			DevFlavor:          values[3].(types.DevFlavor),
			OtelTracingEnabled: values[4].(bool),
			OtelLoggerName:     values[5].(string),
			Output:             values[6].(types.OutputType),
			OtelServiceName:    values[7].(string),
			SetAsDefault:       values[8].(bool),
			ContextKeys:        []interface{}{},
			ContextKeysDefault: nil,
		}
	})
}

// genAnyLevel generates random slog.Leveler values, including nil
func genAnyLevel() gopter.Gen {
	return gen.OneGenOf(
		gen.Const(nil),
		gen.Const(slog.LevelDebug),
		gen.Const(slog.LevelInfo),
		gen.Const(slog.LevelWarn),
		gen.Const(slog.LevelError),
	)
}

// genAnyFormat generates random LogFormat values, including zero value
func genAnyFormat() gopter.Gen {
	return gen.OneGenOf(
		gen.Const(types.LogFormat{}),
		gen.Const(types.LogFormatJSON),
		gen.Const(types.LogFormatText),
		gen.Const(types.LogFormatOtel),
	)
}

// genAnyDevFlavor generates random DevFlavor values, including zero value
func genAnyDevFlavor() gopter.Gen {
	return gen.OneGenOf(
		gen.Const(types.DevFlavor{}),
		gen.Const(types.DevFlavorTint),
		gen.Const(types.DevFlavorSlogor),
		gen.Const(types.DevFlavorDevslog),
	)
}

// genAnyOutput generates random OutputType values, including zero value
func genAnyOutput() gopter.Gen {
	return gen.OneGenOf(
		gen.Const(types.OutputType{}),
		gen.Const(types.OutputConsole),
		gen.Const(types.OutputOtel),
		gen.Const(types.OutputFanout),
	)
}

// interfaceSlice converts a string slice to an interface slice
func interfaceSlice(strs []string) []interface{} {
	result := make([]interface{}, len(strs))
	for i, s := range strs {
		result[i] = s
	}
	return result
}

// **Feature: logger-library-audit-improvements, Property 5: Concurrent logging safety**
// For any number of goroutines logging simultaneously, all log calls should complete without data races
// **Validates: Requirements 4.1**
func TestProperty_ConcurrentLoggingSafety(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("concurrent logging is safe", prop.ForAll(
		func(numGoroutines int, numIterations int) bool {
			// Limit to reasonable ranges
			if numGoroutines < 1 || numGoroutines > 50 {
				return true
			}
			if numIterations < 1 || numIterations > 100 {
				return true
			}

			// Initialize logger
			ctx := context.Background()
			config := types.Config{
				Level:        slog.LevelInfo,
				Format:       types.LogFormatJSON,
				Output:       types.OutputConsole,
				SetAsDefault: false,
			}

			_, logger, err := Init(ctx, config)
			if err != nil {
				t.Errorf("Init failed: %v", err)
				return false
			}

			// Track panics
			panicked := false
			var panicMu sync.Mutex

			// Launch concurrent goroutines that log
			done := make(chan bool, numGoroutines)
			for i := 0; i < numGoroutines; i++ {
				go func(id int) {
					defer func() {
						if r := recover(); r != nil {
							panicMu.Lock()
							panicked = true
							panicMu.Unlock()
							t.Errorf("Goroutine %d panicked: %v", id, r)
						}
						done <- true
					}()

					for j := 0; j < numIterations; j++ {
						logger.Info("test message", "goroutine", id, "iteration", j)
					}
				}(i)
			}

			// Wait for all goroutines
			for i := 0; i < numGoroutines; i++ {
				<-done
			}

			return !panicked
		},
		gen.IntRange(1, 50),
		gen.IntRange(1, 100),
	))

	properties.TestingRun(t)
}

// **Feature: logger-library-audit-improvements, Property 6: Atomic level changes**
// For any sequence of concurrent log level changes and log calls, the level change should be applied atomically
// **Validates: Requirements 4.2**
func TestProperty_AtomicLevelChanges(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("level changes are atomic", prop.ForAll(
		func(numGoroutines int, numIterations int) bool {
			// Limit to reasonable ranges
			if numGoroutines < 1 || numGoroutines > 20 {
				return true
			}
			if numIterations < 1 || numIterations > 50 {
				return true
			}

			// Initialize logger
			ctx := context.Background()
			config := types.Config{
				Level:        slog.LevelInfo,
				Format:       types.LogFormatJSON,
				Output:       types.OutputConsole,
				SetAsDefault: false,
			}

			_, logger, err := Init(ctx, config)
			if err != nil {
				t.Errorf("Init failed: %v", err)
				return false
			}

			levelVar := GetLogLevelAccessor()

			// Track panics
			panicked := false
			var panicMu sync.Mutex

			levels := []slog.Level{
				slog.LevelDebug,
				slog.LevelInfo,
				slog.LevelWarn,
				slog.LevelError,
			}

			// Launch concurrent goroutines that change levels and log
			done := make(chan bool, numGoroutines)
			for i := 0; i < numGoroutines; i++ {
				go func(id int) {
					defer func() {
						if r := recover(); r != nil {
							panicMu.Lock()
							panicked = true
							panicMu.Unlock()
							t.Errorf("Goroutine %d panicked: %v", id, r)
						}
						done <- true
					}()

					for j := 0; j < numIterations; j++ {
						// Change level
						newLevel := levels[j%len(levels)]
						levelVar.Set(newLevel)

						// Read level and log at that level
						currentLevel := levelVar.Level()
						logger.Log(ctx, currentLevel, "test message", "goroutine", id, "iteration", j)
					}
				}(i)
			}

			// Wait for all goroutines
			for i := 0; i < numGoroutines; i++ {
				<-done
			}

			// Verify we can still read the level (it should be one of the valid levels)
			finalLevel := levelVar.Level()
			validLevel := false
			for _, level := range levels {
				if finalLevel == level {
					validLevel = true
					break
				}
			}

			return !panicked && validLevel
		},
		gen.IntRange(1, 20),
		gen.IntRange(1, 50),
	))

	properties.TestingRun(t)
}

// **Feature: logger-library-audit-improvements, Property 8: Concurrent initialization safety**
// For any number of concurrent Init() calls, the system should handle them safely without races or panics
// **Validates: Requirements 4.4**
func TestProperty_ConcurrentInitializationSafety(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("concurrent Init calls are safe", prop.ForAll(
		func(numGoroutines int) bool {
			// Limit to reasonable range
			if numGoroutines < 1 || numGoroutines > 20 {
				return true
			}

			// Track panics and errors
			panicked := false
			var panicMu sync.Mutex
			var errorCount int
			var errorMu sync.Mutex

			// Launch concurrent Init calls
			done := make(chan bool, numGoroutines)
			for i := 0; i < numGoroutines; i++ {
				go func(id int) {
					defer func() {
						if r := recover(); r != nil {
							panicMu.Lock()
							panicked = true
							panicMu.Unlock()
							t.Errorf("Goroutine %d panicked during Init: %v", id, r)
						}
						done <- true
					}()

					ctx := context.Background()
					config := types.Config{
						Level:        slog.LevelInfo,
						Format:       types.LogFormatJSON,
						Output:       types.OutputConsole,
						SetAsDefault: false,
					}

					_, logger, err := Init(ctx, config)
					if err != nil {
						errorMu.Lock()
						errorCount++
						errorMu.Unlock()
					}

					// Try to use the logger if initialization succeeded
					if logger != nil {
						logger.Info("test message", "goroutine", id)
					}
				}(i)
			}

			// Wait for all goroutines
			for i := 0; i < numGoroutines; i++ {
				<-done
			}

			// The test passes if no panics occurred
			// Some errors are acceptable (e.g., if Init has internal state that prevents concurrent calls)
			// but panics are not acceptable
			return !panicked
		},
		gen.IntRange(1, 20),
	))

	properties.TestingRun(t)
}
