package lib

import (
	"log/slog"

	"github.com/wasilak/loggergo/lib/types"
)

// deprecated: use types.DevFlavorFromString instead
func DevFlavorFromString(name string) types.DevFlavor {
	return types.DevFlavorFromString(name)
}

// deprecated: use types.AllDevFlavors instead
func LogLevelFromString(name string) slog.Level {
	return types.LogLevelFromString(name)
}

// deprecated: use types.LogFormatFromString instead
func LogFormatFromString(name string) types.LogFormat {
	return types.LogFormatFromString(name)
}

// deprecated: use types.OutputTypeFromString instead
func OutputTypeFromString(name string) types.OutputType {
	return types.OutputTypeFromString(name)
}
