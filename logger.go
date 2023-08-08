package loggergo

import (
	"os"
	"strings"

	"log/slog"

	otelgoslog "github.com/wasilak/otelgo/slog"
)

// The `LoggerInit` function initializes a logger with a specified log level and log format, allowing
// the user to choose between JSON or text format.
func LoggerInit(level string, logFormat string) {

	// The code block is assigning a value to the `logLevel` variable based on the value of the `level`
	// parameter passed to the `LoggerInit` function. It uses a switch statement to check the lowercase
	// value of `level` and assigns the corresponding `slog.Level` value to `logLevel`. If the value of
	// `level` is "info", it assigns `slog.LevelInfo` to `logLevel`. If the value is "error", it assigns
	// `slog.LevelError`, and so on. If the value of `level` does not match any of the cases, it assigns
	// `slog.LevelInfo` as the default value for `logLevel`.
	var logLevel slog.Leveler

	switch strings.ToLower(level) {
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

	// The code block is checking the value of the `logFormat` parameter passed to the `LoggerInit`
	// function. If the lowercase value of `logFormat` is equal to "json", it sets the default logger to a
	// new logger with a JSON log format. It does this by calling `slog.NewJSONHandler` with `os.Stderr`
	// as the output destination and `&opts` as the options. The resulting handler is then passed to
	// `otelgoslog.NewTracingHandler`, which adds tracing functionality to the logger. Finally, the
	// resulting tracing handler is passed to `slog.New`, which creates a new logger with the specified
	// handler, and `slog.SetDefault` is called to set this logger as the default logger.
	if strings.ToLower(logFormat) == "json" {
		slog.SetDefault(slog.New(otelgoslog.NewTracingHandler(slog.NewJSONHandler(os.Stderr, &opts))))
	} else {
		slog.SetDefault(slog.New(otelgoslog.NewTracingHandler(slog.NewTextHandler(os.Stderr, &opts))))
	}
}
