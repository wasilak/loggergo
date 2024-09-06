// Package loggergo provides functionality for configuring and setting up different logging modes in Go applications.
// It includes support for OpenTelemetry format, JSON format, and plain format with different flavors.
// The package also supports enabling OpenTelemetry tracing for the logs.
package loggergo

import (
	"context"
	"fmt"
	"io"
	"os"

	"log/slog"

	slogmulti "github.com/samber/slog-multi"

	"dario.cat/mergo"
)

// Config represents the configuration options for the LoggerGo logger.
type Config struct {
	Level              slog.Leveler `json:"level"`             // Level specifies the log level. Valid values are any of the slog.Level constants (e.g., slog.LevelInfo, slog.LevelError). Default is slog.LevelInfo.
	Format             LogFormat    `json:"format"`            // Format specifies the log format. Valid values are loggergo.LogFormatText, loggergo.LogFormatJSON, and loggergo.LogFormatOtel. Default is loggergo.LogFormatJSON.
	DevMode            bool         `json:"dev_mode"`          // Dev indicates whether the logger is running in development mode.
	DevFlavor          DevFlavor    `json:"dev_flavor"`        // DevFlavor specifies the development flavor. Valid values are loggergo.DevFlavorTint and loggergo.DevFlavorSlogor. Default is loggergo.DevFlavorTint.
	OutputStream       io.Writer    `json:"output_stream"`     // OutputStream specifies the output stream for the logger. Valid values are "stdout" (default) and "stderr".
	OtelTracingEnabled bool         `json:"otel_enabled"`      // OtelTracingEnabled specifies whether OpenTelemetry support is enabled. Default is true.
	OtelLoggerName     string       `json:"otel_logger_name"`  // OtelLoggerName specifies the name of the logger for OpenTelemetry.
	Output             OutputType   `json:"output"`            // Output specifies the type of output for the logger. Valid values are loggergo.OutputConsole, loggergo.OutputOtel, and loggergo.OutputFanout. Default is loggergo.OutputConsole.
	OtelServiceName    string       `json:"otel_service_name"` // OtelServiceName specifies the service name for OpenTelemetry.
	SetAsDefault       bool         `json:"set_as_default"`    // SetAsDefault specifies whether the logger should be set as the default logger.
}

var defaultConfig = Config{
	Level:              slog.LevelInfo,
	Format:             LogFormatJSON,
	DevMode:            false,
	DevFlavor:          DevFlavorTint,
	OutputStream:       os.Stdout,
	OtelTracingEnabled: true,
	OtelLoggerName:     "my/pkg/name",
	Output:             OutputConsole,
	OtelServiceName:    "my-service",
	SetAsDefault:       true,
}

// The LoggerInit function initializes a logger with the provided configuration and additional
// attributes.
func LoggerInit(ctx context.Context, config Config, additionalAttrs ...any) (*slog.Logger, error) {
	var defaultHandler slog.Handler

	err := mergo.Merge(&defaultConfig, config, mergo.WithOverride)
	if err != nil {
		return nil, err
	}

	opts := slog.HandlerOptions{
		Level:     defaultConfig.Level,
		AddSource: defaultConfig.Level == slog.LevelDebug,
	}

	switch defaultConfig.Output {
	case OutputConsole:
		defaultHandler, err = consoleMode(defaultConfig, opts)
		if err != nil {
			return nil, err
		}
	case OutputOtel:
		defaultHandler, err = otelMode(ctx, defaultConfig)
		if err != nil {
			return nil, err
		}
	case OutputFanout:
		consoleModeHandler, err := consoleMode(defaultConfig, opts)
		if err != nil {
			return nil, err
		}
		otelModeHandler, err := otelMode(ctx, defaultConfig)
		if err != nil {
			return nil, err
		}

		defaultHandler = slogmulti.Fanout(
			consoleModeHandler,
			otelModeHandler,
		)
	default:
		return nil, fmt.Errorf("invalid mode: %s. Valid options: [loggergo.OutputConsole, loggergo.OutputOtel, loggergo.OutputFanout] ", defaultConfig.Output)
	}

	logger := slog.New(defaultHandler)

	// The code below is iterating over the `additionalAttrs` slice and calling the `With` method on the default logger for
	// each element in the slice.
	for _, v := range additionalAttrs {
		logger.With(v)
	}

	if defaultConfig.SetAsDefault {
		// The code `slog.SetDefault(logger)` is setting the default logger to the newly created logger.
		slog.SetDefault(logger)
	}

	return logger, nil
}
