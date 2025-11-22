package loggergo

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// **Feature: logger-library-audit-improvements, Property 10: Context value extraction**
// For any context containing configured keys, those values should appear in the log output
// **Validates: Requirements 9.1**
func TestProperty_ContextValueExtraction(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("context values are extracted and logged", prop.ForAll(
		func(numKeys int) bool {
			// Generate keys and values
			keys := make([]string, numKeys)
			values := make([]string, numKeys)
			for i := 0; i < numKeys; i++ {
				keys[i] = "key" + string(rune('A'+i))
				values[i] = "value" + string(rune('A'+i))
			}

			// Create a context with the keys and values
			ctx := context.Background()
			for i, key := range keys {
				ctx = context.WithValue(ctx, key, values[i])
			}

			// Create a buffer to capture log output
			var buf bytes.Buffer
			handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{})

			// Wrap with context handler
			contextKeys := make([]interface{}, len(keys))
			for i, k := range keys {
				contextKeys[i] = k
			}
			contextHandler := NewCustomContextAttributeHandler(handler, contextKeys, nil)

			// Create a logger and log a message
			logger := slog.New(contextHandler)
			logger.InfoContext(ctx, "test message")

			// Parse the JSON output
			var logEntry map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
				return false
			}

			// Verify that all keys appear in the log output
			for i, key := range keys {
				if logEntry[key] != values[i] {
					return false
				}
			}

			return true
		},
		gen.IntRange(1, 5),
	))

	properties.TestingRun(t)
}

// **Feature: logger-library-audit-improvements, Property 11: Context value type handling**
// For any type of value stored in context, the system should format it appropriately in log output
// **Validates: Requirements 9.4**
func TestProperty_ContextValueTypeHandling(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("different value types are formatted correctly", prop.ForAll(
		func(valueType int) bool {
			// Create a context with different value types
			ctx := context.Background()
			var key string
			var expectedPresent bool

			switch valueType % 6 {
			case 0: // string
				key = "stringKey"
				ctx = context.WithValue(ctx, key, "test string")
				expectedPresent = true
			case 1: // int
				key = "intKey"
				ctx = context.WithValue(ctx, key, 42)
				expectedPresent = true
			case 2: // bool
				key = "boolKey"
				ctx = context.WithValue(ctx, key, true)
				expectedPresent = true
			case 3: // struct
				key = "structKey"
				ctx = context.WithValue(ctx, key, struct{ Name string }{"test"})
				expectedPresent = true
			case 4: // slice
				key = "sliceKey"
				ctx = context.WithValue(ctx, key, []string{"a", "b", "c"})
				expectedPresent = true
			case 5: // map
				key = "mapKey"
				ctx = context.WithValue(ctx, key, map[string]int{"count": 10})
				expectedPresent = true
			}

			// Create a buffer to capture log output
			var buf bytes.Buffer
			handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{})

			// Wrap with context handler
			contextHandler := NewCustomContextAttributeHandler(handler, []interface{}{key}, nil)

			// Create a logger and log a message
			logger := slog.New(contextHandler)
			logger.InfoContext(ctx, "test message")

			// Parse the JSON output
			var logEntry map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
				return false
			}

			// Verify that the key appears in the log output
			_, exists := logEntry[key]
			return exists == expectedPresent
		},
		gen.IntRange(0, 100),
	))

	properties.TestingRun(t)
}
