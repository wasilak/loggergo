package lib

import (
	"log/slog"
	"strings"

	"github.com/wasilak/loggergo/lib/types"
)

func DevFlavorFromString(name string) types.DevFlavor {
	switch strings.ToLower(name) {
	case "tint":
		return types.DevFlavorTint
	case "slogor":
		return types.DevFlavorSlogor
	case "devslog":
		return types.DevFlavorDevslog
	default:
		return types.DevFlavorTint
	}
}

func LogLevelFromString(name string) slog.Level {
	switch strings.ToLower(name) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
