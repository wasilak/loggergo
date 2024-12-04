package types

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

func OutputTypeFromString(name string) OutputType {
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
