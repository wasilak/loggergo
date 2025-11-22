// Package lib provides internal configuration management for the logger.
//
// Configuration Merge Behavior Warning:
//
// The MergeConfig function has important behavior regarding boolean fields that users
// should be aware of. Because Go cannot distinguish between "not set" and "explicitly
// set to false" for boolean fields, any boolean field not explicitly set in an override
// Config will default to false and override a true value in the base config.
//
// This means partial configuration overrides require careful handling of boolean fields.
// See MergeConfig documentation for details and workarounds.
package lib

import (
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/wasilak/loggergo/lib/types"
)

// configManager provides thread-safe access to the global configuration
type configManager struct {
	mu            sync.RWMutex
	config        types.Config
	cleanupFuncs  []func() error
	cleanupMu     sync.Mutex
}

// globalConfigManager is the singleton instance for configuration management
var globalConfigManager = &configManager{}

// InitConfig initializes the global configuration with default values.
//
// This function is called automatically by Init() and should not typically be called directly.
// It sets up sensible defaults for all configuration fields.
//
// Thread Safety:
//
// InitConfig is safe to call concurrently from multiple goroutines.
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

// GetConfig returns a copy of the current configuration in a thread-safe manner.
//
// Thread Safety:
//
// GetConfig is safe to call concurrently from multiple goroutines. It uses a read lock
// to ensure consistent reads even when SetConfig or MergeConfig are being called.
//
// Example:
//
//	config := lib.GetConfig()
//	fmt.Printf("Current log level: %v\n", config.Level)
func GetConfig() types.Config {
	globalConfigManager.mu.RLock()
	defer globalConfigManager.mu.RUnlock()
	return globalConfigManager.config
}

// SetConfig sets the configuration in a thread-safe manner.
//
// This function completely replaces the current configuration with the provided one.
// For partial updates, use MergeConfig instead.
//
// Thread Safety:
//
// SetConfig is safe to call concurrently from multiple goroutines. It uses a write lock
// to ensure exclusive access during the update.
//
// Example:
//
//	newConfig := types.Config{
//	    Level:  slog.LevelDebug,
//	    Format: types.LogFormatJSON,
//	    Output: types.OutputConsole,
//	}
//	lib.SetConfig(newConfig)
func SetConfig(config types.Config) {
	globalConfigManager.mu.Lock()
	defer globalConfigManager.mu.Unlock()
	globalConfigManager.config = config
}

// MergeConfig merges an override configuration with the current base configuration.
//
// Merge Precedence Rules:
//   - Non-zero values in override config replace values in base config
//   - Zero values in override config are ignored (base config values retained)
//   - For pointer fields (Level, OutputStream), nil values are ignored
//   - For slice fields (ContextKeys), empty slices are ignored
//   - For string fields, empty strings are ignored
//   - For enum fields (Format, DevFlavor, Output), zero values are ignored
//
// IMPORTANT - Boolean Field Behavior:
//   Boolean fields (DevMode, OtelTracingEnabled, SetAsDefault) have special behavior
//   because Go cannot distinguish between "not set" and "explicitly set to false".
//   
//   When creating a partial override Config:
//     - If you don't set a boolean field, it defaults to false
//     - This false value WILL override a true value in the base config
//     - This is counterintuitive but unavoidable without using pointer booleans
//   
//   Example Problem:
//     base := Config{DevMode: true, Level: slog.LevelInfo}
//     override := Config{Level: slog.LevelDebug}  // DevMode not set, defaults to false
//     result := MergeConfig(override)
//     // result.DevMode is now false (not true as you might expect!)
//   
//   Workaround:
//     To preserve boolean values when doing partial overrides, you must explicitly
//     set them in the override config:
//     override := Config{Level: slog.LevelDebug, DevMode: true}  // Explicitly preserve DevMode
//   
//   Best Practice:
//     - For full config replacement: use SetConfig() directly
//     - For partial overrides: explicitly set all boolean fields you want to preserve
//     - Consider using InitConfig() + SetConfig() instead of MergeConfig() for clarity
func MergeConfig(override types.Config) types.Config {
	libConfig := GetConfig()

	// Enum fields: override if non-zero
	if override.Format != (types.LogFormat{}) {
		libConfig.Format = override.Format
	}
	if override.DevFlavor != (types.DevFlavor{}) {
		libConfig.DevFlavor = override.DevFlavor
	}
	if override.Output != (types.OutputType{}) {
		libConfig.Output = override.Output
	}

	// Pointer fields: override if non-nil
	if override.Level != nil {
		libConfig.Level = override.Level
	}
	if override.OutputStream != nil {
		libConfig.OutputStream = override.OutputStream
	}

	// String fields: override if non-empty
	if override.OtelLoggerName != "" {
		libConfig.OtelLoggerName = override.OtelLoggerName
	}
	if override.OtelServiceName != "" {
		libConfig.OtelServiceName = override.OtelServiceName
	}

	// Boolean fields: We need special handling to allow false to override true
	// We only skip the override if both values are the same (no change intended)
	if override.DevMode != libConfig.DevMode {
		libConfig.DevMode = override.DevMode
	}
	if override.OtelTracingEnabled != libConfig.OtelTracingEnabled {
		libConfig.OtelTracingEnabled = override.OtelTracingEnabled
	}
	if override.SetAsDefault != libConfig.SetAsDefault {
		libConfig.SetAsDefault = override.SetAsDefault
	}

	// Slice fields: override if non-empty
	if len(override.ContextKeys) > 0 {
		libConfig.ContextKeys = override.ContextKeys
	}

	// Interface fields: override if non-nil
	if override.ContextKeysDefault != nil {
		libConfig.ContextKeysDefault = override.ContextKeysDefault
	}

	// Save the merged config back to the global config manager
	SetConfig(libConfig)
	
	return libConfig
}

// RegisterCleanup registers a cleanup function to be called during Shutdown.
//
// Cleanup functions are called in reverse order of registration (LIFO).
// This ensures that resources are cleaned up in the opposite order they were created.
//
// Thread Safety:
//
// RegisterCleanup is safe to call concurrently from multiple goroutines.
//
// Example:
//
//	lib.RegisterCleanup(func() error {
//	    return provider.Shutdown(context.Background())
//	})
func RegisterCleanup(cleanup func() error) {
	globalConfigManager.cleanupMu.Lock()
	defer globalConfigManager.cleanupMu.Unlock()
	globalConfigManager.cleanupFuncs = append(globalConfigManager.cleanupFuncs, cleanup)
}

// Shutdown performs cleanup of all registered resources.
//
// This function should be called when the application is shutting down to ensure
// proper cleanup of resources like OTEL providers, file handles, etc.
//
// Cleanup functions are called in reverse order of registration (LIFO).
// If any cleanup function returns an error, Shutdown continues with remaining
// cleanup functions and returns a combined error at the end.
//
// Thread Safety:
//
// Shutdown is safe to call concurrently, but should typically only be called once
// during application shutdown.
//
// Example:
//
//	ctx, logger, err := loggergo.Init(ctx, config)
//	if err != nil {
//	    panic(err)
//	}
//	defer loggergo.Shutdown()
//	
//	// Use logger...
func Shutdown() error {
	globalConfigManager.cleanupMu.Lock()
	defer globalConfigManager.cleanupMu.Unlock()
	
	var errors []error
	
	// Call cleanup functions in reverse order (LIFO)
	for i := len(globalConfigManager.cleanupFuncs) - 1; i >= 0; i-- {
		if err := globalConfigManager.cleanupFuncs[i](); err != nil {
			errors = append(errors, err)
		}
	}
	
	// Clear cleanup functions after execution
	globalConfigManager.cleanupFuncs = nil
	
	if len(errors) > 0 {
		return fmt.Errorf("shutdown errors: %v", errors)
	}
	
	return nil
}
