package loggergo

import "log/slog"

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
