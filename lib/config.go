package lib

import (
	"log/slog"
	"os"
	"sync"

	"github.com/wasilak/loggergo/lib/types"
)

// configManager provides thread-safe access to the global configuration
type configManager struct {
	mu     sync.RWMutex
	config types.Config
}

// globalConfigManager is the singleton instance for configuration management
var globalConfigManager = &configManager{}

func InitConfig() {
	globalConfigManager.mu.Lock()
	defer globalConfigManager.mu.Unlock()
	
	globalConfigManager.config = types.Config{
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

// GetConfig returns a copy of the current configuration in a thread-safe manner
func GetConfig() types.Config {
	globalConfigManager.mu.RLock()
	defer globalConfigManager.mu.RUnlock()
	return globalConfigManager.config
}

// SetConfig sets the configuration in a thread-safe manner
func SetConfig(config types.Config) {
	globalConfigManager.mu.Lock()
	defer globalConfigManager.mu.Unlock()
	globalConfigManager.config = config
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

	// Save the merged config back to the global config manager
	SetConfig(libConfig)
	
	return libConfig
}
