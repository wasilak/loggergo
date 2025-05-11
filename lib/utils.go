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

func MergeConfig(def, override types.Config) types.Config {
	// Initialize with default values
	result := types.Config{
		Level:              def.Level,
		Format:             types.LogFormatJSON, // Initialize with default
		DevMode:            def.DevMode,
		DevFlavor:          types.DevFlavorTint, // Initialize with default
		OutputStream:       def.OutputStream,
		OtelTracingEnabled: def.OtelTracingEnabled,
		OtelLoggerName:     def.OtelLoggerName,
		Output:             types.OutputConsole,
		OtelServiceName:    def.OtelServiceName,
		SetAsDefault:       def.SetAsDefault,
		ContextKeys:        def.ContextKeys,
		ContextKeysDefault: def.ContextKeysDefault,
	}

	// If defaults are set, use them
	if override.Format.String() != "" {
		result.Format = override.Format
	}
	if override.DevFlavor.String() != "" {
		result.DevFlavor = override.DevFlavor
	}
	if override.Output.String() != "" {
		result.Output = override.Output
	}

	// Only override non-zero values from override config
	if override.Level != nil {
		result.Level = override.Level
	}
	if override.OutputStream != nil {
		result.OutputStream = override.OutputStream
	}
	if override.OtelLoggerName != "" {
		result.OtelLoggerName = override.OtelLoggerName
	}
	if override.OtelServiceName != "" {
		result.OtelServiceName = override.OtelServiceName
	}
	if override.DevMode {
		result.DevMode = override.DevMode
	}
	if override.OtelTracingEnabled {
		result.OtelTracingEnabled = override.OtelTracingEnabled
	}
	if override.SetAsDefault {
		result.SetAsDefault = override.SetAsDefault
	}
	if len(override.ContextKeys) > 0 {
		result.ContextKeys = override.ContextKeys
	}
	if override.ContextKeysDefault != nil {
		result.ContextKeysDefault = override.ContextKeysDefault
	}

	return result
}
