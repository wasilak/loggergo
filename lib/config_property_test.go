package lib

import (
	"log/slog"
	"sync"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/wasilak/loggergo/lib/types"
)

// **Feature: logger-library-audit-improvements, Property 7: Thread-safe configuration access**
// For any concurrent reads and writes to configuration, no data races should occur
// **Validates: Requirements 4.3**
func TestProperty_ThreadSafeConfigurationAccess(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("concurrent config access is race-free", prop.ForAll(
		func(configs []types.Config) bool {
			// Reset to a known state
			InitConfig()
			
			var wg sync.WaitGroup
			numGoroutines := 10
			
			// Launch multiple goroutines that read and write config concurrently
			for i := 0; i < numGoroutines; i++ {
				wg.Add(1)
				go func(idx int) {
					defer wg.Done()
					
					// Perform multiple operations
					for j := 0; j < 10; j++ {
						// Write operation
						if idx%2 == 0 && len(configs) > 0 {
							configIdx := (idx + j) % len(configs)
							SetConfig(configs[configIdx])
						}
						
						// Read operation
						_ = GetConfig()
					}
				}(i)
			}
			
			wg.Wait()
			
			// If we get here without a race, the test passes
			return true
		},
		genConfigSlice(),
	))

	properties.TestingRun(t)
}

// genConfigSlice generates a slice of random Config instances
func genConfigSlice() gopter.Gen {
	return gen.SliceOfN(5, genValidConfig())
}

// genValidConfig generates a valid Config instance
func genValidConfig() gopter.Gen {
	return gopter.CombineGens(
		genLevel(),
		genFormat(),
		gen.Bool(),
		genDevFlavor(),
		gen.Bool(),
		gen.AlphaString(),
		genOutput(),
		gen.AlphaString(),
		gen.Bool(),
	).Map(func(values []interface{}) types.Config {
		output := values[6].(types.OutputType)
		config := types.Config{
			Level:              values[0].(slog.Leveler),
			Format:             values[1].(types.LogFormat),
			DevMode:            values[2].(bool),
			DevFlavor:          values[3].(types.DevFlavor),
			OtelTracingEnabled: values[4].(bool),
			OtelLoggerName:     values[5].(string),
			Output:             output,
			OtelServiceName:    values[7].(string),
			SetAsDefault:       values[8].(bool),
			ContextKeys:        []interface{}{},
			ContextKeysDefault: nil,
		}
		
		// Ensure OTEL fields are set for OTEL/Fanout modes
		if output.String() == types.OutputOtel.String() || output.String() == types.OutputFanout.String() {
			if config.OtelLoggerName == "" {
				config.OtelLoggerName = "test-logger"
			}
			if config.OtelServiceName == "" {
				config.OtelServiceName = "test-service"
			}
		}
		
		return config
	})
}

// genLevel generates random slog.Leveler values (excluding nil for valid configs)
func genLevel() gopter.Gen {
	return gen.OneGenOf(
		gen.Const(slog.LevelDebug),
		gen.Const(slog.LevelInfo),
		gen.Const(slog.LevelWarn),
		gen.Const(slog.LevelError),
	)
}

// genFormat generates random LogFormat values
func genFormat() gopter.Gen {
	return gen.OneGenOf(
		gen.Const(types.LogFormatJSON),
		gen.Const(types.LogFormatText),
		gen.Const(types.LogFormatOtel),
	)
}

// genDevFlavor generates random DevFlavor values
func genDevFlavor() gopter.Gen {
	return gen.OneGenOf(
		gen.Const(types.DevFlavorTint),
		gen.Const(types.DevFlavorSlogor),
		gen.Const(types.DevFlavorDevslog),
	)
}

// genOutput generates random OutputType values (excluding zero value for valid configs)
func genOutput() gopter.Gen {
	return gen.OneGenOf(
		gen.Const(types.OutputConsole),
		gen.Const(types.OutputOtel),
		gen.Const(types.OutputFanout),
	)
}

// **Feature: logger-library-audit-improvements, Property 3: Configuration round-trip**
// For any valid Config, setting it via SetConfig and retrieving it via GetConfig should return an equivalent configuration
// **Validates: Requirements 2.4**
func TestProperty_ConfigurationRoundTrip(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("config round-trip preserves values", prop.ForAll(
		func(config types.Config) bool {
			// Set the config
			SetConfig(config)
			
			// Get the config back
			retrieved := GetConfig()
			
			// Compare all fields
			return configsEqual(config, retrieved)
		},
		genValidConfig(),
	))

	properties.TestingRun(t)
}

// **Feature: logger-library-audit-improvements, Property 2: Configuration merge consistency**
// For any two Config instances, merging them should follow consistent precedence rules where non-zero values in the override config replace defaults
// **Validates: Requirements 2.2**
func TestProperty_ConfigurationMergeConsistency(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("merge follows consistent precedence rules", prop.ForAll(
		func(base, override types.Config) bool {
			// Set the base config
			SetConfig(base)
			
			// Merge with override
			result := MergeConfig(override)
			
			// Verify merge precedence rules:
			// Non-zero values in override should replace base values
			
			// Enum fields: non-zero override values should be used
			if override.Format != (types.LogFormat{}) {
				if result.Format.String() != override.Format.String() {
					return false
				}
			} else {
				if result.Format.String() != base.Format.String() {
					return false
				}
			}
			
			if override.DevFlavor != (types.DevFlavor{}) {
				if result.DevFlavor.String() != override.DevFlavor.String() {
					return false
				}
			} else {
				if result.DevFlavor.String() != base.DevFlavor.String() {
					return false
				}
			}
			
			if override.Output != (types.OutputType{}) {
				if result.Output.String() != override.Output.String() {
					return false
				}
			} else {
				if result.Output.String() != base.Output.String() {
					return false
				}
			}
			
			// Pointer fields: non-nil override values should be used
			if override.Level != nil {
				if result.Level == nil || result.Level.Level() != override.Level.Level() {
					return false
				}
			} else {
				if base.Level != nil {
					if result.Level == nil || result.Level.Level() != base.Level.Level() {
						return false
					}
				}
			}
			
			// String fields: non-empty override values should be used
			if override.OtelLoggerName != "" {
				if result.OtelLoggerName != override.OtelLoggerName {
					return false
				}
			} else {
				if result.OtelLoggerName != base.OtelLoggerName {
					return false
				}
			}
			
			if override.OtelServiceName != "" {
				if result.OtelServiceName != override.OtelServiceName {
					return false
				}
			} else {
				if result.OtelServiceName != base.OtelServiceName {
					return false
				}
			}
			
			// Boolean fields: override if different from base
			// This allows false to override true
			if override.DevMode != base.DevMode {
				if result.DevMode != override.DevMode {
					return false
				}
			} else {
				if result.DevMode != base.DevMode {
					return false
				}
			}
			
			if override.OtelTracingEnabled != base.OtelTracingEnabled {
				if result.OtelTracingEnabled != override.OtelTracingEnabled {
					return false
				}
			} else {
				if result.OtelTracingEnabled != base.OtelTracingEnabled {
					return false
				}
			}
			
			if override.SetAsDefault != base.SetAsDefault {
				if result.SetAsDefault != override.SetAsDefault {
					return false
				}
			} else {
				if result.SetAsDefault != base.SetAsDefault {
					return false
				}
			}
			
			// Slice fields: non-empty override values should be used
			if len(override.ContextKeys) > 0 {
				if len(result.ContextKeys) != len(override.ContextKeys) {
					return false
				}
				for i := range override.ContextKeys {
					if result.ContextKeys[i] != override.ContextKeys[i] {
						return false
					}
				}
			} else {
				if len(result.ContextKeys) != len(base.ContextKeys) {
					return false
				}
			}
			
			// Interface fields: non-nil override values should be used
			if override.ContextKeysDefault != nil {
				if result.ContextKeysDefault != override.ContextKeysDefault {
					return false
				}
			} else {
				if result.ContextKeysDefault != base.ContextKeysDefault {
					return false
				}
			}
			
			return true
		},
		genValidConfig(),
		genValidConfig(),
	))

	properties.TestingRun(t)
}

// configsEqual compares two Config instances for equality
func configsEqual(a, b types.Config) bool {
	// Compare Level
	if (a.Level == nil) != (b.Level == nil) {
		return false
	}
	if a.Level != nil && b.Level != nil {
		if a.Level.Level() != b.Level.Level() {
			return false
		}
	}
	
	// Compare Format
	if a.Format.String() != b.Format.String() {
		return false
	}
	
	// Compare DevMode
	if a.DevMode != b.DevMode {
		return false
	}
	
	// Compare DevFlavor
	if a.DevFlavor.String() != b.DevFlavor.String() {
		return false
	}
	
	// Compare OtelTracingEnabled
	if a.OtelTracingEnabled != b.OtelTracingEnabled {
		return false
	}
	
	// Compare OtelLoggerName
	if a.OtelLoggerName != b.OtelLoggerName {
		return false
	}
	
	// Compare Output
	if a.Output.String() != b.Output.String() {
		return false
	}
	
	// Compare OtelServiceName
	if a.OtelServiceName != b.OtelServiceName {
		return false
	}
	
	// Compare SetAsDefault
	if a.SetAsDefault != b.SetAsDefault {
		return false
	}
	
	// Compare ContextKeys length
	if len(a.ContextKeys) != len(b.ContextKeys) {
		return false
	}
	
	// Compare ContextKeys elements
	for i := range a.ContextKeys {
		if a.ContextKeys[i] != b.ContextKeys[i] {
			return false
		}
	}
	
	// Compare ContextKeysDefault
	if a.ContextKeysDefault != b.ContextKeysDefault {
		return false
	}
	
	// Note: We don't compare OutputStream because io.Writer is an interface
	// and may not be directly comparable
	
	return true
}
