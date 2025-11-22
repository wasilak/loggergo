package loggergo

import (
	"log/slog"

	"github.com/wasilak/loggergo/lib/types"
)

// Types provides access to all type constants, enums, and conversion functions.
//
// This struct exports all the type-related functionality from the internal types package,
// making it easy to access log formats, output types, dev flavors, and their conversion functions.
//
// Example usage:
//
//	// Using constants
//	config := loggergo.Config{
//	    Format:    loggergo.Types.LogFormatJSON,
//	    Output:    loggergo.Types.OutputConsole,
//	    DevFlavor: loggergo.Types.DevFlavorTint,
//	}
//
//	// Using conversion functions
//	format := loggergo.Types.LogFormatFromString("json")
//	output := loggergo.Types.OutputTypeFromString("console")
//
//	// Getting all available values
//	allFormats := loggergo.Types.AllLogFormats()
//	for _, format := range allFormats {
//	    fmt.Println(format)
//	}
var Types = struct {
	AllDevFlavors       func() []types.DevFlavor
	DevFlavorFromString func(string) types.DevFlavor
	DevFlavorTint       types.DevFlavor
	DevFlavorSlogor     types.DevFlavor
	DevFlavorDevslog    types.DevFlavor

	AllLogFormats       func() []types.LogFormat
	LogFormatFromString func(string) types.LogFormat
	LogFormatText       types.LogFormat
	LogFormatJSON       types.LogFormat
	LogFormatOtel       types.LogFormat

	AllLogLevels       func() []slog.Level
	LogLevelFromString func(string) slog.Level

	AllOutputTypes       func() []types.OutputType
	OutputTypeFromString func(string) types.OutputType
	OutputConsole        types.OutputType
	OutputOtel           types.OutputType
	OutputFanout         types.OutputType
}{
	AllDevFlavors:       types.AllDevFlavors,
	DevFlavorFromString: types.DevFlavorFromString,
	DevFlavorTint:       types.DevFlavorTint,
	DevFlavorSlogor:     types.DevFlavorSlogor,
	DevFlavorDevslog:    types.DevFlavorDevslog,

	AllLogFormats:       types.AllLogFormats,
	LogFormatFromString: types.LogFormatFromString,
	LogFormatText:       types.LogFormatText,
	LogFormatJSON:       types.LogFormatJSON,
	LogFormatOtel:       types.LogFormatOtel,

	AllLogLevels:       types.AllLogLevels,
	LogLevelFromString: types.LogLevelFromString,

	AllOutputTypes:       types.AllOutputTypes,
	OutputTypeFromString: types.OutputTypeFromString,
	OutputConsole:        types.OutputConsole,
	OutputOtel:           types.OutputOtel,
	OutputFanout:         types.OutputFanout,
}
