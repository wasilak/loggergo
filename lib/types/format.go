package types

import (
	"log/slog"

	"github.com/xybor-x/enum"
)

// LogFormat represents the format of the log.
type logFormat int
type LogFormat struct{ enum.SafeEnum[logFormat] }

var (
	// LogFormatJSON represents JSON format.
	LogFormatJSON = enum.NewExtended[LogFormat]("json")
	// LogFormatText represents text format.
	LogFormatText = enum.NewExtended[LogFormat]("text")
	// LogFormatOtel represents OTEL (JSON) format.
	LogFormatOtel = enum.NewExtended[LogFormat]("otel")
	_             = enum.Finalize[LogFormat]() // still required internally
)

// AllDevFlavors returns all defined DevFlavor values.
func AllLogFormats() []LogFormat {
	return enum.All[LogFormat]()
}

func LogFormatFromString(name string) LogFormat {
	if v, ok := enum.FromString[LogFormat](name); ok {
		return v
	}
	slog.Warn("Invalid log format: %q, defaulting to %s", name, LogFormatText)
	return LogFormatText
}
