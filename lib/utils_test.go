package lib

import (
	"log/slog"
	"testing"

	"github.com/wasilak/loggergo/lib/types"
)

// TestDeprecatedFunctions verifies that deprecated functions still work
// and delegate to their new implementations correctly.

func TestDevFlavorFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected types.DevFlavor
	}{
		{"tint flavor", "tint", types.DevFlavorTint},
		{"slogor flavor", "slogor", types.DevFlavorSlogor},
		{"devslog flavor", "devslog", types.DevFlavorDevslog},
		{"empty string defaults to tint", "", types.DevFlavorTint},
		{"invalid flavor defaults to tint", "invalid", types.DevFlavorTint},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test deprecated function
			deprecatedResult := DevFlavorFromString(tt.input)
			// Test new function
			newResult := types.DevFlavorFromString(tt.input)

			// Verify they produce the same result
			if deprecatedResult != newResult {
				t.Errorf("DevFlavorFromString(%q) = %v, want %v (from types.DevFlavorFromString)", 
					tt.input, deprecatedResult, newResult)
			}

			// Verify expected result
			if deprecatedResult != tt.expected {
				t.Errorf("DevFlavorFromString(%q) = %v, want %v", 
					tt.input, deprecatedResult, tt.expected)
			}
		})
	}
}

func TestLogLevelFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected slog.Level
	}{
		{"debug level", "debug", slog.LevelDebug},
		{"info level", "info", slog.LevelInfo},
		{"warn level", "warn", slog.LevelWarn},
		{"error level", "error", slog.LevelError},
		{"uppercase", "INFO", slog.LevelInfo},
		{"empty string defaults to info", "", slog.LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test deprecated function
			deprecatedResult := LogLevelFromString(tt.input)
			// Test new function
			newResult := types.LogLevelFromString(tt.input)

			// Verify they produce the same result
			if deprecatedResult != newResult {
				t.Errorf("LogLevelFromString(%q) = %v, want %v (from types.LogLevelFromString)", 
					tt.input, deprecatedResult, newResult)
			}

			// Verify expected result
			if deprecatedResult != tt.expected {
				t.Errorf("LogLevelFromString(%q) = %v, want %v", 
					tt.input, deprecatedResult, tt.expected)
			}
		})
	}
}

func TestLogFormatFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected types.LogFormat
	}{
		{"json format", "json", types.LogFormatJSON},
		{"text format", "text", types.LogFormatText},
		{"otel format", "otel", types.LogFormatOtel},
		{"uppercase defaults to text", "JSON", types.LogFormatText},
		{"empty string defaults to text", "", types.LogFormatText},
		{"invalid format defaults to text", "invalid", types.LogFormatText},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test deprecated function
			deprecatedResult := LogFormatFromString(tt.input)
			// Test new function
			newResult := types.LogFormatFromString(tt.input)

			// Verify they produce the same result
			if deprecatedResult != newResult {
				t.Errorf("LogFormatFromString(%q) = %v, want %v (from types.LogFormatFromString)", 
					tt.input, deprecatedResult, newResult)
			}

			// Verify expected result
			if deprecatedResult != tt.expected {
				t.Errorf("LogFormatFromString(%q) = %v, want %v", 
					tt.input, deprecatedResult, tt.expected)
			}
		})
	}
}

func TestOutputTypeFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected types.OutputType
	}{
		{"console output", "console", types.OutputConsole},
		{"otel output", "otel", types.OutputOtel},
		{"fanout output", "fanout", types.OutputFanout},
		{"uppercase", "CONSOLE", types.OutputConsole},
		{"empty string defaults to console", "", types.OutputConsole},
		{"invalid output defaults to console", "invalid", types.OutputConsole},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test deprecated function
			deprecatedResult := OutputTypeFromString(tt.input)
			// Test new function
			newResult := types.OutputTypeFromString(tt.input)

			// Verify they produce the same result
			if deprecatedResult != newResult {
				t.Errorf("OutputTypeFromString(%q) = %v, want %v (from types.OutputTypeFromString)", 
					tt.input, deprecatedResult, newResult)
			}

			// Verify expected result
			if deprecatedResult != tt.expected {
				t.Errorf("OutputTypeFromString(%q) = %v, want %v", 
					tt.input, deprecatedResult, tt.expected)
			}
		})
	}
}
