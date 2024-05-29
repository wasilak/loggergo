package loggergo

import (
	"os"
	"strings"
	"time"

	"log/slog"

	otelgoslog "github.com/wasilak/otelgo/slog"

	"github.com/golang-cz/devslog"

	"dario.cat/mergo"

	"gitlab.com/greyxor/slogor"
)

// The LoggerGoConfig type is a configuration struct for a logger in Go, with fields for level, format,
// and dev mode.
// @property {string} Level - The "Level" property in the LoggerGoConfig struct represents the logging
// level. It determines the severity of the log messages that will be recorded. Common levels include
// "debug", "info", "warning", "error", and "fatal".
// @property {string} Format - The `Format` property in the `LoggerGoConfig` struct represents the
// desired format for the log messages. It specifies how the log messages should be formatted when they
// are written to the log output.
// @property {bool} Dev - The `Dev` property is a boolean flag that indicates whether the logger is
// running in development mode or not. It can be used to enable or disable certain logging features or
// behaviors specific to development environments.
type LoggerGoConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
	Dev    bool   `json:"dev"`
}

// The line `var defaultConfig = LoggerGoConfig{ Level: "info", Format: "plain", Dev: false }` is
// initializing a variable named `defaultConfig` with a default configuration for the logger. It sets
// the `Level` property to "info", indicating that the logger should record log messages with a
// severity level of "info" or higher. The `Format` property is set to "plain", specifying that the log
// messages should be formatted in a plain text format. The `Dev` property is set to `false`,
// indicating that the logger is not running in development mode.
var defaultConfig = LoggerGoConfig{
	Level:  "info",
	Format: "plain",
	Dev:    false,
}

// The LoggerInit function initializes a logger with the provided configuration and additional
// attributes.
func LoggerInit(config LoggerGoConfig, additionalAttrs ...any) (*slog.Logger, error) {

	err := mergo.Merge(&defaultConfig, config, mergo.WithOverride)
	if err != nil {
		return nil, err
	}

	var logLevel slog.Leveler

	// The `switch` statement is used to evaluate the value of `defaultConfig.Level` and assign a corresponding
	// `slog.Level` value to the `logLevel` variable.
	switch strings.ToLower(defaultConfig.Level) {
	case "info":
		logLevel = slog.LevelInfo
	case "error":
		logLevel = slog.LevelError
	case "warn":
		logLevel = slog.LevelWarn
	case "debug":
		logLevel = slog.LevelDebug
	default:
		logLevel = slog.LevelInfo
	}

	opts := slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	}

	// The `if` statement is checking if the value of `defaultConfig.Format` is equal to "json". If it is, it
	// sets the default logger handler to a new `slog.NewJSONHandler` with the provided options. This
	// means that log messages will be formatted as JSON when written to the log output.
	if strings.ToLower(defaultConfig.Format) == "json" {
		slog.SetDefault(slog.New(otelgoslog.NewTracingHandler(slog.NewJSONHandler(os.Stderr, &opts))))
	} else {
		if defaultConfig.Dev {
			devOpts := &devslog.Options{
				HandlerOptions:    &opts,
				MaxSlicePrintSize: 10,
				SortKeys:          true,
			}

			slog.SetDefault(slog.New(slogor.NewHandler(os.Stderr, slogor.Options{
				TimeFormat: time.Stamp,
				Level:      slog.LevelError,
				ShowSource: false,
			})))

			slog.SetDefault(slog.New(otelgoslog.NewTracingHandler(devslog.NewHandler(os.Stderr, devOpts))))
		} else {
			slog.SetDefault(slog.New(otelgoslog.NewTracingHandler(slog.NewTextHandler(os.Stderr, &opts))))
		}
	}

	// The code `for _, v := range additionalAttrs { slog.SetDefault(slog.Default().With(v)) }` is
	// iterating over the `additionalAttrs` slice and calling the `With` method on the default logger for
	// each element in the slice.
	for _, v := range additionalAttrs {
		slog.SetDefault(slog.Default().With(v))
	}

	return slog.Default(), nil
}
