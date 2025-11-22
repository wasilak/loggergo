package types

import (
	"log/slog"
	"strings"
	"testing"
)

// TestConfig_Validate_NilLevel tests validation when Level is nil
func TestConfig_Validate_NilLevel(t *testing.T) {
	config := Config{
		Level:  nil, // Invalid
		Output: OutputConsole,
	}

	err := config.Validate()
	if err == nil {
		t.Fatal("Expected validation error for nil Level, but got nil")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "Level") {
		t.Errorf("Expected error message to mention 'Level', but got: %v", err)
	}
}

// TestConfig_Validate_MissingOutput tests validation when Output is zero value
func TestConfig_Validate_MissingOutput(t *testing.T) {
	config := Config{
		Level:  slog.LevelInfo,
		Output: OutputType{}, // Invalid - zero value
	}

	err := config.Validate()
	if err == nil {
		t.Fatal("Expected validation error for missing Output, but got nil")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "Output") {
		t.Errorf("Expected error message to mention 'Output', but got: %v", err)
	}
}

// TestConfig_Validate_EmptyOtelLoggerName tests validation when OtelLoggerName is empty with OTEL mode
func TestConfig_Validate_EmptyOtelLoggerName(t *testing.T) {
	config := Config{
		Level:           slog.LevelInfo,
		Output:          OutputOtel,
		OtelLoggerName:  "", // Invalid for OTEL mode
		OtelServiceName: "test-service",
	}

	err := config.Validate()
	if err == nil {
		t.Fatal("Expected validation error for empty OtelLoggerName with OTEL output, but got nil")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "OtelLoggerName") {
		t.Errorf("Expected error message to mention 'OtelLoggerName', but got: %v", err)
	}
}

// TestConfig_Validate_EmptyOtelServiceName tests validation when OtelServiceName is empty with OTEL mode
func TestConfig_Validate_EmptyOtelServiceName(t *testing.T) {
	config := Config{
		Level:           slog.LevelInfo,
		Output:          OutputOtel,
		OtelLoggerName:  "test-logger",
		OtelServiceName: "", // Invalid for OTEL mode
	}

	err := config.Validate()
	if err == nil {
		t.Fatal("Expected validation error for empty OtelServiceName with OTEL output, but got nil")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "OtelServiceName") {
		t.Errorf("Expected error message to mention 'OtelServiceName', but got: %v", err)
	}
}

// TestConfig_Validate_EmptyOtelFieldsWithFanout tests validation when OTEL fields are empty with Fanout mode
func TestConfig_Validate_EmptyOtelFieldsWithFanout(t *testing.T) {
	config := Config{
		Level:           slog.LevelInfo,
		Output:          OutputFanout,
		OtelLoggerName:  "", // Invalid for Fanout mode
		OtelServiceName: "", // Invalid for Fanout mode
	}

	err := config.Validate()
	if err == nil {
		t.Fatal("Expected validation error for empty OTEL fields with Fanout output, but got nil")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "OtelLoggerName") && !strings.Contains(errMsg, "OtelServiceName") {
		t.Errorf("Expected error message to mention OTEL fields, but got: %v", err)
	}
}

// TestConfig_Validate_ContextKeysDefaultWithoutKeys tests validation when ContextKeysDefault is set without ContextKeys
func TestConfig_Validate_ContextKeysDefaultWithoutKeys(t *testing.T) {
	config := Config{
		Level:              slog.LevelInfo,
		Output:             OutputConsole,
		ContextKeysDefault: "default-value", // Invalid without ContextKeys
		ContextKeys:        []interface{}{}, // Empty
	}

	err := config.Validate()
	if err == nil {
		t.Fatal("Expected validation error for ContextKeysDefault without ContextKeys, but got nil")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "ContextKeysDefault") {
		t.Errorf("Expected error message to mention 'ContextKeysDefault', but got: %v", err)
	}
}

// TestConfig_Validate_ContextKeysDefaultWithNilKeys tests validation when ContextKeysDefault is set with nil ContextKeys
func TestConfig_Validate_ContextKeysDefaultWithNilKeys(t *testing.T) {
	config := Config{
		Level:              slog.LevelInfo,
		Output:             OutputConsole,
		ContextKeysDefault: "default-value", // Invalid without ContextKeys
		ContextKeys:        nil,             // Nil
	}

	err := config.Validate()
	if err == nil {
		t.Fatal("Expected validation error for ContextKeysDefault with nil ContextKeys, but got nil")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "ContextKeysDefault") {
		t.Errorf("Expected error message to mention 'ContextKeysDefault', but got: %v", err)
	}
}

// TestConfig_Validate_ValidConsoleConfig tests validation with valid Console configuration
func TestConfig_Validate_ValidConsoleConfig(t *testing.T) {
	config := Config{
		Level:  slog.LevelInfo,
		Output: OutputConsole,
	}

	err := config.Validate()
	if err != nil {
		t.Errorf("Expected no validation error for valid Console config, but got: %v", err)
	}
}

// TestConfig_Validate_ValidOtelConfig tests validation with valid OTEL configuration
func TestConfig_Validate_ValidOtelConfig(t *testing.T) {
	config := Config{
		Level:           slog.LevelInfo,
		Output:          OutputOtel,
		OtelLoggerName:  "test-logger",
		OtelServiceName: "test-service",
	}

	err := config.Validate()
	if err != nil {
		t.Errorf("Expected no validation error for valid OTEL config, but got: %v", err)
	}
}

// TestConfig_Validate_ValidFanoutConfig tests validation with valid Fanout configuration
func TestConfig_Validate_ValidFanoutConfig(t *testing.T) {
	config := Config{
		Level:           slog.LevelInfo,
		Output:          OutputFanout,
		OtelLoggerName:  "test-logger",
		OtelServiceName: "test-service",
	}

	err := config.Validate()
	if err != nil {
		t.Errorf("Expected no validation error for valid Fanout config, but got: %v", err)
	}
}

// TestConfig_Validate_ValidContextKeys tests validation with valid ContextKeys configuration
func TestConfig_Validate_ValidContextKeys(t *testing.T) {
	config := Config{
		Level:              slog.LevelInfo,
		Output:             OutputConsole,
		ContextKeys:        []interface{}{"key1", "key2"},
		ContextKeysDefault: "default-value",
	}

	err := config.Validate()
	if err != nil {
		t.Errorf("Expected no validation error for valid ContextKeys config, but got: %v", err)
	}
}

// TestConfig_Validate_MultipleErrors tests that validation returns multiple errors
func TestConfig_Validate_MultipleErrors(t *testing.T) {
	config := Config{
		Level:  nil,          // Error 1
		Output: OutputType{}, // Error 2
	}

	err := config.Validate()
	if err == nil {
		t.Fatal("Expected validation error for multiple invalid fields, but got nil")
	}

	errMsg := err.Error()
	// Should mention at least one error and indicate multiple errors
	if !strings.Contains(errMsg, "Level") && !strings.Contains(errMsg, "Output") {
		t.Errorf("Expected error message to mention at least one field, but got: %v", err)
	}
	// Should indicate multiple errors
	if !strings.Contains(errMsg, "2 errors") && !strings.Contains(errMsg, "and 1 more") {
		t.Errorf("Expected error message to indicate multiple errors, but got: %v", err)
	}
}

// TestValidationError_Error tests the ValidationError.Error() method
func TestValidationError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      ValidationError
		contains []string
	}{
		{
			name: "single error",
			err: ValidationError{
				Errors: []FieldError{
					{Field: "Level", Value: nil, Reason: "cannot be nil"},
				},
			},
			contains: []string{"Level", "cannot be nil"},
		},
		{
			name: "multiple errors",
			err: ValidationError{
				Errors: []FieldError{
					{Field: "Level", Value: nil, Reason: "cannot be nil"},
					{Field: "Output", Value: OutputType{}, Reason: "must be specified"},
				},
			},
			// Only check for first error and indication of multiple errors
			contains: []string{"Level", "2 errors"},
		},
		{
			name:     "no errors",
			err:      ValidationError{Errors: []FieldError{}},
			contains: []string{"validation failed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errMsg := tt.err.Error()
			for _, substr := range tt.contains {
				if !strings.Contains(errMsg, substr) {
					t.Errorf("Expected error message to contain %q, but got: %s", substr, errMsg)
				}
			}
		})
	}
}

// TestFieldError_Error tests the FieldError.Error() method
func TestFieldError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      FieldError
		contains []string
	}{
		{
			name: "with value",
			err: FieldError{
				Field:  "Level",
				Value:  "invalid",
				Reason: "must be a valid level",
			},
			contains: []string{"Level", "invalid", "must be a valid level"},
		},
		{
			name: "without value",
			err: FieldError{
				Field:  "Output",
				Value:  nil,
				Reason: "cannot be nil",
			},
			contains: []string{"Output", "cannot be nil"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errMsg := tt.err.Error()
			for _, substr := range tt.contains {
				if !strings.Contains(errMsg, substr) {
					t.Errorf("Expected error message to contain %q, but got: %s", substr, errMsg)
				}
			}
		})
	}
}
