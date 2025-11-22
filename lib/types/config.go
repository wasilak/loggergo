package types

import (
	"fmt"
	"io"
	"log/slog"
)

// ValidationError represents configuration validation failures.
// It contains a slice of FieldError instances, each describing a specific validation failure.
//
// Example:
//
//	err := config.Validate()
//	if err != nil {
//	    var valErr *ValidationError
//	    if errors.As(err, &valErr) {
//	        for _, fieldErr := range valErr.Errors {
//	            fmt.Printf("Field %s: %s\n", fieldErr.Field, fieldErr.Reason)
//	        }
//	    }
//	}
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
// It contains the field name, its value, and the reason for the validation failure.
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

// InitError represents an error during logger initialization.
// It wraps the underlying error with context about which stage failed.
//
// Stages can be:
//   - "validation": Configuration validation failed
//   - "handler_creation": Handler creation failed
//   - "otel_setup": OpenTelemetry setup failed
//   - "panic_recovery": A panic occurred during initialization
//
// Example:
//
//	_, _, err := loggergo.Init(ctx, config)
//	if err != nil {
//	    var initErr *InitError
//	    if errors.As(err, &initErr) {
//	        fmt.Printf("Initialization failed at %s: %v\n", initErr.Stage, initErr.Cause)
//	    }
//	}
type InitError struct {
	Stage  string // Stage indicates which initialization stage failed (e.g., "validation", "handler_creation", "otel_setup")
	Cause  error  // Cause is the underlying error that caused the initialization to fail
	Config Config // Config is a sanitized copy of the configuration (sensitive data should be removed)
}

// Error implements the error interface for InitError.
func (e *InitError) Error() string {
	if e.Stage == "" {
		return fmt.Sprintf("logger initialization failed: %v", e.Cause)
	}
	return fmt.Sprintf("logger initialization failed at %s: %v", e.Stage, e.Cause)
}

// Unwrap returns the underlying cause error, allowing errors.Is and errors.As to work.
func (e *InitError) Unwrap() error {
	return e.Cause
}

// Config represents the configuration options for the logger.
//
// All fields have sensible defaults, so a zero-value Config can be used after setting
// the required Level field.
//
// Example (minimal):
//
//	config := loggergo.Config{
//	    Level: slog.LevelInfo,
//	}
//
// Example (full):
//
//	config := loggergo.Config{
//	    Level:              slog.LevelDebug,
//	    Format:             loggergo.LogFormatJSON,
//	    Output:             loggergo.OutputConsole,
//	    DevMode:            true,
//	    DevFlavor:          loggergo.DevFlavorTint,
//	    OutputStream:       os.Stdout,
//	    SetAsDefault:       true,
//	    ContextKeys:        []interface{}{"request_id", "user_id"},
//	    ContextKeysDefault: "unknown",
//	}
//
// OTEL Example:
//
//	config := loggergo.Config{
//	    Level:              slog.LevelInfo,
//	    Output:             loggergo.OutputOtel,
//	    OtelLoggerName:     "myapp/logger",
//	    OtelServiceName:    "myapp",
//	    OtelTracingEnabled: true,
//	}
type Config struct {
	Level              slog.Leveler  `json:"level"`                // Level specifies the log level. Valid values are any of the slog.Level constants (e.g., slog.LevelInfo, slog.LevelError). Default: slog.LevelInfo.
	Format             LogFormat     `json:"format"`               // Format specifies the log format. Valid values are loggergo.LogFormatText, loggergo.LogFormatJSON, and loggergo.LogFormatOtel. Default: loggergo.LogFormatJSON.
	DevMode            bool          `json:"dev_mode"`             // DevMode indicates whether the logger is running in development mode. Default: false. WARNING: When using MergeConfig, false will override true. To preserve a true value, explicitly set DevMode to true in the override config.
	DevFlavor          DevFlavor     `json:"dev_flavor"`           // DevFlavor specifies the development flavor. Valid values are loggergo.DevFlavorTint, loggergo.DevFlavorSlogor, and loggergo.DevFlavorDevslog. Default: loggergo.DevFlavorTint.
	OutputStream       io.Writer     `json:"output_stream"`        // OutputStream specifies the output stream for the logger. Default: os.Stdout.
	OtelTracingEnabled bool          `json:"otel_enabled"`         // OtelTracingEnabled specifies whether OpenTelemetry support is enabled. Default: true. WARNING: When using MergeConfig, false will override true. To preserve a true value, explicitly set OtelTracingEnabled to true in the override config.
	OtelLoggerName     string        `json:"otel_logger_name"`     // OtelLoggerName specifies the name of the logger for OpenTelemetry. Default: "my/pkg/name". Required when Output is OutputOtel or OutputFanout.
	Output             OutputType    `json:"output"`               // Output specifies the type of output for the logger. Valid values are loggergo.OutputConsole, loggergo.OutputOtel, and loggergo.OutputFanout. Default: loggergo.OutputConsole.
	OtelServiceName    string        `json:"otel_service_name"`    // OtelServiceName specifies the service name for OpenTelemetry. Default: "my-service". Required when Output is OutputOtel or OutputFanout.
	SetAsDefault       bool          `json:"set_as_default"`       // SetAsDefault specifies whether the logger should be set as the default logger. Default: true. WARNING: When using MergeConfig, false will override true. To preserve a true value, explicitly set SetAsDefault to true in the override config.
	ContextKeys        []interface{} `json:"context_keys"`         // ContextKeys specifies the keys to be added to log from context. Default: empty slice.
	ContextKeysDefault interface{}   `json:"context_keys_default"` // ContextKeysDefault specifies the default value for the context keys if not found in the context. Default: nil.
}

// Validate checks if the configuration is valid and returns an error if not.
//
// It validates:
//   - Required fields (Level, Output)
//   - Mode-specific requirements (OTEL fields when using OTEL or Fanout output)
//   - Field conflicts (ContextKeysDefault without ContextKeys)
//
// Returns:
//   - nil if the configuration is valid
//   - *ValidationError containing all validation failures
//
// Example:
//
//	config := loggergo.Config{Level: slog.LevelInfo}
//	if err := config.Validate(); err != nil {
//	    log.Fatalf("Invalid configuration: %v", err)
//	}
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
