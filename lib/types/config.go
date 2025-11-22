package types

import (
	"fmt"
	"io"
	"log/slog"
)

// ValidationError represents configuration validation failures.
type ValidationError struct {
	Errors []FieldError
}

// Error implements the error interface for ValidationError.
func (e *ValidationError) Error() string {
	if len(e.Errors) == 0 {
		return "configuration validation failed"
	}
	if len(e.Errors) == 1 {
		return fmt.Sprintf("configuration validation failed: %s", e.Errors[0].Error())
	}
	return fmt.Sprintf("configuration validation failed with %d errors: %s (and %d more)", 
		len(e.Errors), e.Errors[0].Error(), len(e.Errors)-1)
}

// FieldError represents a single field validation failure.
type FieldError struct {
	Field  string
	Value  interface{}
	Reason string
}

// Error implements the error interface for FieldError.
func (e *FieldError) Error() string {
	if e.Value != nil {
		return fmt.Sprintf("field %q with value %v: %s", e.Field, e.Value, e.Reason)
	}
	return fmt.Sprintf("field %q: %s", e.Field, e.Reason)
}

// Config represents the configuration options for the LoggerGo logger.
type Config struct {
	Level              slog.Leveler  `json:"level"`                // Level specifies the log level. Valid values are any of the slog.Level constants (e.g., slog.LevelInfo, slog.LevelError). Default is slog.LevelInfo.
	Format             LogFormat     `json:"format"`               // Format specifies the log format. Valid values are loggergo.LogFormatText, loggergo.LogFormatJSON, and loggergo.LogFormatOtel. Default is loggergo.LogFormatJSON.
	DevMode            bool          `json:"dev_mode"`             // Dev indicates whether the logger is running in development mode.
	DevFlavor          DevFlavor     `json:"dev_flavor"`           // DevFlavor specifies the development flavor. Valid values are loggergo.DevFlavorTint and loggergo.DevFlavorSlogor. Default is loggergo.DevFlavorTint.
	OutputStream       io.Writer     `json:"output_stream"`        // OutputStream specifies the output stream for the logger. Valid values are "stdout" (default) and "stderr".
	OtelTracingEnabled bool          `json:"otel_enabled"`         // OtelTracingEnabled specifies whether OpenTelemetry support is enabled. Default is true.
	OtelLoggerName     string        `json:"otel_logger_name"`     // OtelLoggerName specifies the name of the logger for OpenTelemetry.
	Output             OutputType    `json:"output"`               // Output specifies the type of output for the logger. Valid values are loggergo.OutputConsole, loggergo.OutputOtel, and loggergo.OutputFanout. Default is loggergo.OutputConsole.
	OtelServiceName    string        `json:"otel_service_name"`    // OtelServiceName specifies the service name for OpenTelemetry.
	SetAsDefault       bool          `json:"set_as_default"`       // SetAsDefault specifies whether the logger should be set as the default logger.
	ContextKeys        []interface{} `json:"context_keys"`         // ContextKeys specifies the keys to be added to log from context.
	ContextKeysDefault interface{}   `json:"context_keys_default"` // ContextKeysDefault specifies the default value for the context keys if not found in the context.
}

// Validate checks if the configuration is valid and returns an error if not.
// It validates required fields, field conflicts, and mode-specific requirements.
func (c *Config) Validate() error {
	var fieldErrors []FieldError

	// Validate level
	if c.Level == nil {
		fieldErrors = append(fieldErrors, FieldError{
			Field:  "Level",
			Value:  nil,
			Reason: "cannot be nil",
		})
	}

	// Validate output mode
	if c.Output == (OutputType{}) {
		fieldErrors = append(fieldErrors, FieldError{
			Field:  "Output",
			Value:  c.Output,
			Reason: "must be specified (Console, OTEL, or Fanout)",
		})
	}

	// Validate OTEL-specific fields - only check for missing required fields
	if c.Output.String() == OutputOtel.String() || c.Output.String() == OutputFanout.String() {
		if c.OtelLoggerName == "" {
			fieldErrors = append(fieldErrors, FieldError{
				Field:  "OtelLoggerName",
				Value:  c.OtelLoggerName,
				Reason: "required when Output is OTEL or Fanout",
			})
		}
		if c.OtelServiceName == "" {
			fieldErrors = append(fieldErrors, FieldError{
				Field:  "OtelServiceName",
				Value:  c.OtelServiceName,
				Reason: "required when Output is OTEL or Fanout",
			})
		}
	}

	// Validate context keys
	if c.ContextKeysDefault != nil && len(c.ContextKeys) == 0 {
		fieldErrors = append(fieldErrors, FieldError{
			Field:  "ContextKeysDefault",
			Value:  c.ContextKeysDefault,
			Reason: "cannot be set without defining ContextKeys",
		})
	}

	if len(fieldErrors) > 0 {
		return &ValidationError{Errors: fieldErrors}
	}

	return nil
}
