// Package loggergo provides functionality for configuring and setting up different logging modes in Go applications.
// It includes support for OpenTelemetry format, JSON format, and plain format with different flavors.
// The package also supports enabling OpenTelemetry tracing for the logs.
package loggergo

import (
	"context"
	"fmt"
	"os"
	"strings"

	"log/slog"

	slogmulti "github.com/samber/slog-multi"

	"dario.cat/mergo"
)

// Config represents the configuration options for the LoggerGo logger.
type Config struct {
	Level              string     `json:"level"`             // Level specifies the log level. Valid values are "debug", "info", "warn", and "error".
	Format             string     `json:"format"`            // Format specifies the log format. Valid values are "plain" (default), "otel" and "json".
	DevMode            bool       `json:"dev_mode"`          // Dev indicates whether the logger is running in development mode.
	DevFlavor          string     `json:"dev_flavor"`        // DevFlavor specifies the development flavor. Valid values are "tint" (default), slogor and "devslog".
	OutputStream       *os.File   `json:"output_stream"`     // OutputStream specifies the output stream for the logger. Valid values are "stdout" (default) and "stderr".
	OtelTracingEnabled bool       `json:"otel_enabled"`      // OtelTracingEnabled specifies whether OpenTelemetry support is enabled. Default is true.
	OtelLoggerName     string     `json:"otel_logger_name"`  // OtelLoggerName specifies the name of the logger for OpenTelemetry.
	Output             OutputType `json:"output"`            // Output specifies where logs will be sent to. Valid values are "console", "otel" and "fanout" (default) - which is a combination of "console" and "otel".
	OtelServiceName    string     `json:"otel_service_name"` // OtelServiceName specifies the service name for OpenTelemetry.
}

// The line `var defaultConfig = Config{ Level: "info", Format: "plain", Dev: false }` is
// initializing a variable named `defaultConfig` with a default configuration for the logger. It sets
// the `Level` property to "info", indicating that the logger should record log messages with a
// severity level of "info" or higher. The `Format` property is set to "plain", specifying that the log
// messages should be formatted in a plain text format. The `Dev` property is set to `false`,
// indicating that the logger is not running in development mode.
var defaultConfig = Config{
	Level:              "info",
	Format:             "plain",
	DevMode:            false,
	DevFlavor:          "tint",
	OutputStream:       os.Stdout,
	OtelTracingEnabled: true,
	OtelLoggerName:     "my/pkg/name",
	Output:             OutputConsole,
	OtelServiceName:    "my-service",
}

// The LoggerInit function initializes a logger with the provided configuration and additional
// attributes.
func LoggerInit(ctx context.Context, config Config, additionalAttrs ...any) (*slog.Logger, error) {
	var defaultHandler slog.Handler

	err := mergo.Merge(&defaultConfig, config, mergo.WithOverride)
	if err != nil {
		return nil, err
	}

	// normalize the log level, mode and format
	defaultConfig.Level = strings.ToLower(defaultConfig.Level)
	defaultConfig.Format = strings.ToLower(defaultConfig.Format)

	logLevel := setupLogLevel()

	opts := slog.HandlerOptions{
		Level:     logLevel,
		AddSource: logLevel == slog.LevelDebug,
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

	// The code `slog.SetDefault(logger)` is setting the default logger to the newly created logger.
	slog.SetDefault(slog.New(defaultHandler))

	// The code `for _, v := range additionalAttrs { slog.SetDefault(slog.Default().With(v)) }` is
	// iterating over the `additionalAttrs` slice and calling the `With` method on the default logger for
	// each element in the slice.
	for _, v := range additionalAttrs {
		slog.Default().With(v)
	}

	return slog.Default(), nil
}
