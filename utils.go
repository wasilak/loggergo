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

func (o OutputType) FromString(name string) OutputType {
	switch name {
	case "console":
		return OutputConsole
	case "otel":
		return OutputOtel
	case "fanout":
		return OutputFanout
	default:
		return OutputConsole
	}
}

// LogFormat represents the format of the log.
type LogFormat int

const (
	// LogFormatText represents text format.
	LogFormatText LogFormat = iota
	// LogFormatJSON represents JSON format.
	LogFormatJSON
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

func (f LogFormat) FromString(name string) LogFormat {
	switch name {
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

// DevFlavor represents the flavor of the development environment.
type DevFlavor int

const (
	// DevFlavorTint represents the "tint" development flavor.
	DevFlavorTint DevFlavor = iota
	// DevFlavorSlogor represents the "slogor" development flavor.
	DevFlavorSlogor
	// DevFlavorDevslog represents the production "devslog" flavor.
	DevFlavorDevslog
)

func (f DevFlavor) String() string {
	switch f {
	case DevFlavorTint:
		return "tint"
	case DevFlavorSlogor:
		return "slogor"
	case DevFlavorDevslog:
		return "devslog"
	default:
		return "unknown"
	}
}

func (f DevFlavor) FromString(name string) DevFlavor {
	switch name {
	case "tint":
		return DevFlavorTint
	case "slogor":
		return DevFlavorSlogor
	case "devslog":
		return DevFlavorDevslog
	default:
		return DevFlavorTint
	}
}

func LogLevelFromString(name string) slog.Level {
	switch name {
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
