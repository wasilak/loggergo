package loggergo

import (
	"os"
	"strings"
	"time"

	"log/slog"

	"github.com/golang-cz/devslog"
	"github.com/mattn/go-isatty"
	otelgoslog "github.com/wasilak/otelgo/slog"
	"gitlab.com/greyxor/slogor"

	"dario.cat/mergo"

	"github.com/lmittmann/tint"
)

// LoggerGoConfig represents the configuration options for the LoggerGo logger.
type LoggerGoConfig struct {
	Level        string `json:"level"`         // Level specifies the log level. Valid values are "debug", "info", "warn", and "error".
	Format       string `json:"format"`        // Format specifies the log format. Valid values are "text" and "json".
	DevMode      bool   `json:"dev_mode"`      // Dev indicates whether the logger is running in development mode.
	DevFlavor    string `json:"dev_flavor"`    // DevFlavor specifies the development flavor. Valid values are "tint" (default), slogor and "devslog".
	OutputStream string `json:"output_stream"` // OutputStream specifies the output stream. Valid values are "stdout" (default) and "stderr".
}

// The line `var defaultConfig = LoggerGoConfig{ Level: "info", Format: "plain", Dev: false }` is
// initializing a variable named `defaultConfig` with a default configuration for the logger. It sets
// the `Level` property to "info", indicating that the logger should record log messages with a
// severity level of "info" or higher. The `Format` property is set to "plain", specifying that the log
// messages should be formatted in a plain text format. The `Dev` property is set to `false`,
// indicating that the logger is not running in development mode.
var defaultConfig = LoggerGoConfig{
	Level:        "info",
	Format:       "plain",
	DevMode:      false,
	DevFlavor:    "tint",
	OutputStream: "stdout",
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
		AddSource: logLevel == slog.LevelDebug,
	}

	var defaultHandler slog.Handler

	var defaultOutputStream = os.Stdout
	if strings.ToLower(defaultConfig.OutputStream) == "stderr" {
		defaultOutputStream = os.Stderr
	}

	// The `if` statement is checking if the value of `defaultConfig.Format` is equal to "json". If it is, it
	// sets the default logger handler to a new `slog.NewJSONHandler` with the provided options. This
	// means that log messages will be formatted as JSON when written to the log output.
	if strings.ToLower(defaultConfig.Format) == "json" {
		defaultHandler = slog.NewJSONHandler(defaultOutputStream, &opts)
	} else {
		if defaultConfig.DevMode {

			if defaultConfig.DevFlavor == "slogor" {
				defaultHandler = slogor.NewHandler(defaultOutputStream, slogor.Options{
					TimeFormat: time.Stamp,
					Level:      opts.Level.Level(),
					ShowSource: opts.AddSource,
				})
			} else if defaultConfig.DevFlavor == "devslog" {
				defaultHandler = devslog.NewHandler(defaultOutputStream, &devslog.Options{
					HandlerOptions:    &opts,
					MaxSlicePrintSize: 10,
					SortKeys:          true,
				})
			} else {
				defaultHandler = tint.NewHandler(defaultOutputStream, &tint.Options{
					Level:     opts.Level,
					NoColor:   !isatty.IsTerminal(os.Stderr.Fd()),
					AddSource: opts.AddSource,
				})
			}
		} else {
			defaultHandler = slog.NewTextHandler(defaultOutputStream, &opts)
		}
	}

	slog.SetDefault(slog.New(otelgoslog.NewTracingHandler(defaultHandler)))

	// The code `for _, v := range additionalAttrs { slog.SetDefault(slog.Default().With(v)) }` is
	// iterating over the `additionalAttrs` slice and calling the `With` method on the default logger for
	// each element in the slice.
	for _, v := range additionalAttrs {
		slog.SetDefault(slog.Default().With(v))
	}

	return slog.Default(), nil
}
