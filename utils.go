package loggergo

import "log/slog"

// OutputType represents the type of output for the logger.
type OutputType int

const (
	// OutputConsole represents console output.
	OutputConsole OutputType = iota
	// OutputOtel represents otel output.
	OutputOtel
	// OutputFanout represents both console and otel output.
	OutputFanout
)

func (o OutputType) String() string {
	switch o {
	case OutputConsole:
		return "console"
	case OutputOtel:
		return "otel"
	case OutputFanout:
		return "fanout"
	default:
		return "unknown"
	}
}

// setupLogLevel sets up the log level based on the value of defaultConfig.Level.
// It evaluates the value of defaultConfig.Level and assigns a corresponding slog.Level value to the logLevel variable.
// If the value of defaultConfig.Level is not recognized, it defaults to slog.LevelInfo.
func setupLogLevel() slog.Leveler {
	var logLevel slog.Leveler

	// The `switch` statement is used to evaluate the value of `defaultConfig.Level` and assign a corresponding
	// `slog.Level` value to the `logLevel` variable.
	switch defaultConfig.Level {
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

	return logLevel
}
