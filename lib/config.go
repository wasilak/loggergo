package lib

import (
	"log/slog"
	"os"

	"github.com/wasilak/loggergo/lib/types"
)

var libConfig types.Config

func InitConfig() {
	libConfig = types.Config{
		Level:              slog.LevelInfo,
		Format:             types.LogFormatJSON,
		DevMode:            false,
		DevFlavor:          types.DevFlavorTint,
		OutputStream:       os.Stdout,
		OtelTracingEnabled: true,
		OtelLoggerName:     "my/pkg/name",
		Output:             types.OutputConsole,
		OtelServiceName:    "my-service",
		SetAsDefault:       true,
		ContextKeys:        []interface{}{},
		ContextKeysDefault: nil,
	}
}

func GetConfig() *types.Config {
	return &libConfig
}

func MergeConfig(override types.Config) types.Config {
	libConfig := GetConfig()

	if libConfig.Format != (types.LogFormat{}) {
		libConfig.Format = override.Format
	}
	if libConfig.DevFlavor != (types.DevFlavor{}) {
		libConfig.DevFlavor = override.DevFlavor
	}
	if override.Output != (types.OutputType{}) {
		libConfig.Output = override.Output
	}
	if override.Level != nil {
		libConfig.Level = override.Level
	}
	if override.OutputStream != nil {
		libConfig.OutputStream = override.OutputStream
	}
	if override.OtelLoggerName != "" {
		libConfig.OtelLoggerName = override.OtelLoggerName
	}
	if override.OtelServiceName != "" {
		libConfig.OtelServiceName = override.OtelServiceName
	}
	if override.DevMode && override.DevMode != libConfig.DevMode {
		libConfig.DevMode = override.DevMode
	}
	if override.OtelTracingEnabled && override.OtelTracingEnabled != libConfig.OtelTracingEnabled {
		libConfig.OtelTracingEnabled = override.OtelTracingEnabled
	}
	if override.SetAsDefault && override.SetAsDefault != libConfig.SetAsDefault {
		libConfig.SetAsDefault = override.SetAsDefault
	}
	if len(override.ContextKeys) > 0 {
		libConfig.ContextKeys = override.ContextKeys
	}
	if override.ContextKeysDefault != nil {
		libConfig.ContextKeysDefault = override.ContextKeysDefault
	}

	return *libConfig
}
