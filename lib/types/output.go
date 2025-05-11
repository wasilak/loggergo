package types

import (
	"fmt"
	"log/slog"

	"github.com/xybor-x/enum"
)

// OutputType represents the type of output for the logger.
type outputType int
type OutputType struct{ enum.SafeEnum[outputType] }

var (
	// OutputConsole represents console output.
	OutputConsole = enum.NewExtended[OutputType]("console")
	// OutputOtel represents otel output.
	OutputOtel = enum.NewExtended[OutputType]("otel")
	// OutputFanout represents both console and otel output.
	OutputFanout = enum.NewExtended[OutputType]("fanout")
	_            = enum.Finalize[OutputType]() // still required internally
)

func AllOutputTypes() []OutputType {
	return enum.All[OutputType]()
}

func OutputTypeFromString(name string) OutputType {
	if v, ok := enum.FromString[OutputType](name); ok {
		return v
	}
	slog.Warn(fmt.Sprintf("Unknown output type: %q, defaulting to %s", name, OutputConsole))
	return OutputConsole
}
