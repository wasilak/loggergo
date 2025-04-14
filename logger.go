// Package loggergo provides functionality for configuring and setting up different logging modes in Go applications.
// It includes support for OpenTelemetry format, JSON format, and plain format with different flavors.
// The package also supports enabling OpenTelemetry tracing for the logs.
package loggergo

import (
	"context"
	"fmt"
	"os"

	"log/slog"

	slogmulti "github.com/samber/slog-multi"

	"dario.cat/mergo"
	"github.com/wasilak/loggergo/lib/modes"
	"github.com/wasilak/loggergo/lib/types"
)

var defaultConfig = types.Config{
	Level:              slog.LevelInfo,
	Format:             types.LogFormatJSON,
	DevMode:            false,
	DevFlavor:          types.DevFlavorTint,
	OutputStream:       os.Stdout,
	OtelTracingEnabled: true,
	OtelLoggerName:     "my/pkg/name",
	Output:             types.OutputConsole,
	OtelServiceName:    "my-service",
	SetAsDefault:       true,
	ContextKeys:        []interface{}{},
	ContextKeysDefault: nil,
}

var logLevel = new(slog.LevelVar)

// expose the types for external usage
type Config = types.Config

// The LoggerInit function initializes a logger with the provided configuration and additional
// attributes.
func Init(ctx context.Context, config types.Config, additionalAttrs ...any) (context.Context, *slog.Logger, error) {
	var defaultHandler slog.Handler

	err := mergo.Merge(&defaultConfig, config, mergo.WithOverride)
	if err != nil {
		return ctx, nil, err
	}

	logLevel.Set(defaultConfig.Level.Level())

	opts := slog.HandlerOptions{
		Level:     logLevel,
		AddSource: defaultConfig.Level == slog.LevelDebug,
	}

	switch defaultConfig.Output {
	case types.OutputConsole:
		defaultHandler, err = modes.ConsoleMode(defaultConfig, opts)
		if err != nil {
			return ctx, nil, err
		}
	case types.OutputOtel:
		defaultHandler, ctx, err = modes.OtelMode(ctx, defaultConfig)
		if err != nil {
			return ctx, nil, err
		}
	case types.OutputFanout:
		consoleModeHandler, err := modes.ConsoleMode(defaultConfig, opts)
		if err != nil {
			return ctx, nil, err
		}
		otelModeHandler, ctx, err := modes.OtelMode(ctx, defaultConfig)
		if err != nil {
			return ctx, nil, err
		}

		defaultHandler = slogmulti.Fanout(
			consoleModeHandler,
			otelModeHandler,
		)
	default:
		return ctx, nil, fmt.Errorf("invalid mode: %s. Valid options: [loggergo.OutputConsole, loggergo.OutputOtel, loggergo.OutputFanout] ", defaultConfig.Output)
	}

	// The code below is creating a new CustomContextAttributeHandler with the default handler and the context keys.
	defaultHandler = NewCustomContextAttributeHandler(defaultHandler, defaultConfig.ContextKeys, defaultConfig.ContextKeysDefault)

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

	return ctx, logger, nil
}

func GetLogLevelAccessor() *slog.LevelVar {
	return logLevel
}
