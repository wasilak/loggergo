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
