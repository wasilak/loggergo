package loggergo

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
