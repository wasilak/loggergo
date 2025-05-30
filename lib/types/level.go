package types

import (
	"fmt"
	"log/slog"
)

func LogLevelFromString(name string) slog.Level {
	var level slog.Level
	err := level.UnmarshalText([]byte(name))
	if err != nil {
		slog.Warn(fmt.Sprintf("Unknown log level: %q, defaulting to %s", name, slog.LevelInfo))
		return slog.LevelInfo
	}
	return level
}

func AllLogLevels() []slog.Level {
	return []slog.Level{
		slog.LevelDebug,
		slog.LevelInfo,
		slog.LevelWarn,
		slog.LevelError,
	}
}
