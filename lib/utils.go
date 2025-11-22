package lib

import (
	"log/slog"

	"github.com/wasilak/loggergo/lib/types"
)

// DevFlavorFromString converts a string to a DevFlavor type.
//
// Deprecated: Use types.DevFlavorFromString instead.
// This function is maintained for backward compatibility and will be removed in v3.0.0.
//
// Migration example:
//
//	// Old code:
//	flavor := lib.DevFlavorFromString("tint")
//
//	// New code:
//	flavor := types.DevFlavorFromString("tint")
func DevFlavorFromString(name string) types.DevFlavor {
	return types.DevFlavorFromString(name)
}

// LogLevelFromString converts a string to a slog.Level.
//
// Deprecated: Use types.LogLevelFromString instead.
// This function is maintained for backward compatibility and will be removed in v3.0.0.
//
// Migration example:
//
//	// Old code:
//	level := lib.LogLevelFromString("info")
//
//	// New code:
//	level := types.LogLevelFromString("info")
func LogLevelFromString(name string) slog.Level {
	return types.LogLevelFromString(name)
}

// LogFormatFromString converts a string to a LogFormat type.
//
// Deprecated: Use types.LogFormatFromString instead.
// This function is maintained for backward compatibility and will be removed in v3.0.0.
//
// Migration example:
//
//	// Old code:
//	format := lib.LogFormatFromString("json")
//
//	// New code:
//	format := types.LogFormatFromString("json")
func LogFormatFromString(name string) types.LogFormat {
	return types.LogFormatFromString(name)
}

// OutputTypeFromString converts a string to an OutputType.
//
// Deprecated: Use types.OutputTypeFromString instead.
// This function is maintained for backward compatibility and will be removed in v3.0.0.
//
// Migration example:
//
//	// Old code:
//	output := lib.OutputTypeFromString("console")
//
//	// New code:
//	output := types.OutputTypeFromString("console")
func OutputTypeFromString(name string) types.OutputType {
	return types.OutputTypeFromString(name)
}
