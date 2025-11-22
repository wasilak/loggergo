package loggergo

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
)

// TestContextHandler_NilContext tests that nil context is handled gracefully
func TestContextHandler_NilContext(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{})
	contextHandler := NewCustomContextAttributeHandler(handler, []interface{}{"key1"}, "default")

	logger := slog.New(contextHandler)
	
	// Pass nil context - should not panic
	logger.Info("test message")

	// Verify log was written
	if buf.Len() == 0 {
		t.Error("Expected log output, got none")
	}

	// Parse the JSON output
	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	// Verify default value was used
	if logEntry["key1"] != "default" {
		t.Errorf("Expected default value 'default', got %v", logEntry["key1"])
	}
}

// TestContextHandler_MissingContextKeys tests handling of missing context keys
func TestContextHandler_MissingContextKeys(t *testing.T) {
	tests := []struct {
		name         string
		contextKeys  []interface{}
		defaultValue interface{}
		wantKey      string
		wantValue    interface{}
		wantPresent  bool
	}{
		{
			name:         "missing key with default",
			contextKeys:  []interface{}{"missingKey"},
			defaultValue: "default",
			wantKey:      "missingKey",
			wantValue:    "default",
			wantPresent:  true,
		},
		{
			name:         "missing key without default",
			contextKeys:  []interface{}{"missingKey"},
			defaultValue: nil,
			wantKey:      "missingKey",
			wantPresent:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{})
			contextHandler := NewCustomContextAttributeHandler(handler, tt.contextKeys, tt.defaultValue)

			logger := slog.New(contextHandler)
			ctx := context.Background()
			logger.InfoContext(ctx, "test message")

			// Parse the JSON output
			var logEntry map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
				t.Fatalf("Failed to parse log output: %v", err)
			}

			// Check if key is present
			val, exists := logEntry[tt.wantKey]
			if exists != tt.wantPresent {
				t.Errorf("Expected key presence %v, got %v", tt.wantPresent, exists)
			}

			if tt.wantPresent && val != tt.wantValue {
				t.Errorf("Expected value %v, got %v", tt.wantValue, val)
			}
		})
	}
}

// TestContextHandler_NilContextValues tests handling of context keys with nil values
func TestContextHandler_NilContextValues(t *testing.T) {
	tests := []struct {
		name         string
		contextValue interface{}
		defaultValue interface{}
		wantValue    interface{}
		wantPresent  bool
	}{
		{
			name:         "nil value with default",
			contextValue: nil,
			defaultValue: "default",
			wantValue:    "default",
			wantPresent:  true,
		},
		{
			name:         "nil value without default",
			contextValue: nil,
			defaultValue: nil,
			wantPresent:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{})
			contextHandler := NewCustomContextAttributeHandler(handler, []interface{}{"key1"}, tt.defaultValue)

			logger := slog.New(contextHandler)
			ctx := context.WithValue(context.Background(), "key1", tt.contextValue)
			logger.InfoContext(ctx, "test message")

			// Parse the JSON output
			var logEntry map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
				t.Fatalf("Failed to parse log output: %v", err)
			}

			// Check if key is present
			val, exists := logEntry["key1"]
			if exists != tt.wantPresent {
				t.Errorf("Expected key presence %v, got %v", tt.wantPresent, exists)
			}

			if tt.wantPresent && val != tt.wantValue {
				t.Errorf("Expected value %v, got %v", tt.wantValue, val)
			}
		})
	}
}

// TestContextHandler_VariousValueTypes tests handling of various value types
func TestContextHandler_VariousValueTypes(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value interface{}
	}{
		{
			name:  "string value",
			key:   "stringKey",
			value: "test string",
		},
		{
			name:  "int value",
			key:   "intKey",
			value: 42,
		},
		{
			name:  "bool value",
			key:   "boolKey",
			value: true,
		},
		{
			name:  "struct value",
			key:   "structKey",
			value: struct{ Name string }{"test"},
		},
		{
			name:  "slice value",
			key:   "sliceKey",
			value: []string{"a", "b", "c"},
		},
		{
			name:  "map value",
			key:   "mapKey",
			value: map[string]int{"count": 10},
		},
		{
			name:  "pointer value",
			key:   "ptrKey",
			value: func() *int { i := 100; return &i }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{})
			contextHandler := NewCustomContextAttributeHandler(handler, []interface{}{tt.key}, nil)

			logger := slog.New(contextHandler)
			ctx := context.WithValue(context.Background(), tt.key, tt.value)
			logger.InfoContext(ctx, "test message")

			// Parse the JSON output
			var logEntry map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
				t.Fatalf("Failed to parse log output: %v", err)
			}

			// Verify the key exists in the output
			if _, exists := logEntry[tt.key]; !exists {
				t.Errorf("Expected key %s to be present in log output", tt.key)
			}
		})
	}
}

// TestContextHandler_MultipleKeys tests handling of multiple context keys
func TestContextHandler_MultipleKeys(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{})
	
	keys := []interface{}{"key1", "key2", "key3"}
	contextHandler := NewCustomContextAttributeHandler(handler, keys, "default")

	logger := slog.New(contextHandler)
	
	// Create context with some keys present, some missing
	ctx := context.Background()
	ctx = context.WithValue(ctx, "key1", "value1")
	ctx = context.WithValue(ctx, "key3", "value3")
	// key2 is missing

	logger.InfoContext(ctx, "test message")

	// Parse the JSON output
	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	// Verify key1 has correct value
	if logEntry["key1"] != "value1" {
		t.Errorf("Expected key1='value1', got %v", logEntry["key1"])
	}

	// Verify key2 has default value
	if logEntry["key2"] != "default" {
		t.Errorf("Expected key2='default', got %v", logEntry["key2"])
	}

	// Verify key3 has correct value
	if logEntry["key3"] != "value3" {
		t.Errorf("Expected key3='value3', got %v", logEntry["key3"])
	}
}

// TestContextHandler_EmptyKeys tests handling when no keys are configured
func TestContextHandler_EmptyKeys(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{})
	contextHandler := NewCustomContextAttributeHandler(handler, []interface{}{}, nil)

	logger := slog.New(contextHandler)
	ctx := context.WithValue(context.Background(), "someKey", "someValue")
	logger.InfoContext(ctx, "test message")

	// Parse the JSON output
	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	// Verify no extra keys were added
	if _, exists := logEntry["someKey"]; exists {
		t.Error("Expected no context keys to be extracted when none are configured")
	}
}
