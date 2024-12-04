package types

import "strings"

// LogFormat represents the format of the log.
type LogFormat int

const (
	// LogFormatJSON represents text format.
	LogFormatJSON LogFormat = iota
	// LogFormatText represents JSON format.
	LogFormatText
	// LogFormatOtel represents OTEL (JSON) format.
	LogFormatOtel
)

func (f LogFormat) String() string {
	switch f {
	case LogFormatText:
		return "text"
	case LogFormatJSON:
		return "json"
	case LogFormatOtel:
		return "otel"
	default:
		return "unknown"
	}
}

func LogFormatFromString(name string) LogFormat {
	switch strings.ToLower(name) {
	case "text":
		return LogFormatText
	case "json":
		return LogFormatJSON
	case "otel":
		return LogFormatOtel
	default:
		return LogFormatText
	}
}
