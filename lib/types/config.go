package types

import (
	"io"
	"log/slog"
)

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
