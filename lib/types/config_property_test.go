package types

import (
	"log/slog"
	"os"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// **Feature: logger-library-audit-improvements, Property 1: Invalid configuration detection**
// For any Config with invalid field values, calling Validate() should return an error identifying the invalid fields
// **Validates: Requirements 1.4, 2.1**
func TestProperty_InvalidConfigurationDetection(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("invalid configs are detected", prop.ForAll(
		func(config Config) bool {
			err := config.Validate()
			
			// Check if config is invalid and should produce an error
			if isInvalidConfig(config) {
				return err != nil
			}
			
			// If config is valid, Validate should not return an error
			return err == nil
		},
		genConfig(),
	))

	properties.TestingRun(t)
}

// isInvalidConfig checks if a config should be considered invalid
func isInvalidConfig(c Config) bool {
	// Level is nil
	if c.Level == nil {
		return true
	}
	
	// Output mode is zero value (not specified)
	if c.Output == (OutputType{}) {
		return true
	}
	
	// OTEL-specific validation
	if c.Output.String() == OutputOtel.String() || c.Output.String() == OutputFanout.String() {
		if c.OtelLoggerName == "" || c.OtelServiceName == "" {
			return true
		}
	}
	
	// ContextKeysDefault set but no ContextKeys
	if c.ContextKeysDefault != nil && len(c.ContextKeys) == 0 {
		return true
	}
	
	return false
}

// genConfig generates random Config instances for property testing
func genConfig() gopter.Gen {
	return gopter.CombineGens(
		genLevel(),
		genFormat(),
		gen.Bool(),
		genDevFlavor(),
		genOutputStream(),
		gen.Bool(),
		gen.AlphaString(),
		genOutput(),
		gen.AlphaString(),
		gen.Bool(),
		genContextKeys(),
		genContextKeysDefault(),
	).Map(func(values []interface{}) Config {
		return Config{
			Level:              values[0].(slog.Leveler),
			Format:             values[1].(LogFormat),
			DevMode:            values[2].(bool),
			DevFlavor:          values[3].(DevFlavor),
			OutputStream:       values[4].(*os.File),
			OtelTracingEnabled: values[5].(bool),
			OtelLoggerName:     values[6].(string),
			Output:             values[7].(OutputType),
			OtelServiceName:    values[8].(string),
			SetAsDefault:       values[9].(bool),
			ContextKeys:        values[10].([]interface{}),
			ContextKeysDefault: values[11],
		}
	})
}

// genLevel generates random slog.Leveler values (including nil)
func genLevel() gopter.Gen {
	return gen.OneGenOf(
		gen.Const(nil),
		gen.Const(slog.LevelDebug),
		gen.Const(slog.LevelInfo),
		gen.Const(slog.LevelWarn),
		gen.Const(slog.LevelError),
	)
}

// genFormat generates random LogFormat values
func genFormat() gopter.Gen {
	return gen.OneGenOf(
		gen.Const(LogFormatJSON),
		gen.Const(LogFormatText),
		gen.Const(LogFormatOtel),
		gen.Const(LogFormat{}), // zero value
	)
}

// genDevFlavor generates random DevFlavor values
func genDevFlavor() gopter.Gen {
	return gen.OneGenOf(
		gen.Const(DevFlavorTint),
		gen.Const(DevFlavorSlogor),
		gen.Const(DevFlavorDevslog),
		gen.Const(DevFlavor{}), // zero value
	)
}

// genOutputStream generates random io.Writer values
func genOutputStream() gopter.Gen {
	return gen.Const(os.Stdout)
}

// genOutput generates random OutputType values (including zero value)
func genOutput() gopter.Gen {
	return gen.OneGenOf(
		gen.Const(OutputConsole),
		gen.Const(OutputOtel),
		gen.Const(OutputFanout),
		gen.Const(OutputType{}), // zero value - invalid
	)
}

// genContextKeys generates random context keys slices
func genContextKeys() gopter.Gen {
	return gen.OneGenOf(
		gen.Const([]interface{}{}),
		gen.SliceOf(gen.AlphaString()).Map(func(strs []string) []interface{} {
			result := make([]interface{}, len(strs))
			for i, s := range strs {
				result[i] = s
			}
			return result
		}),
	)
}

// genContextKeysDefault generates random context keys default values
func genContextKeysDefault() gopter.Gen {
	return gen.OneGenOf(
		gen.Const(nil),
		gen.AlphaString(),
		gen.Int(),
	)
}

// **Feature: logger-library-audit-improvements, Property 12: Required field validation**
// For any Config missing required fields (based on output mode), Validate() should return an error
// **Validates: Requirements 12.3**
func TestProperty_RequiredFieldValidation(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	properties.Property("missing required fields are detected", prop.ForAll(
		func(output OutputType, hasOtelLoggerName bool, hasOtelServiceName bool, hasLevel bool) bool {
			// Build a config with potentially missing required fields
			config := Config{
				Output: output,
			}
			
			// Set Level if hasLevel is true
			if hasLevel {
				config.Level = slog.LevelInfo
			}
			
			// For Console mode, don't set OTEL fields to avoid conflict errors
			// For OTEL/Fanout modes, set fields based on flags
			if output.String() == OutputConsole.String() {
				// Don't set OTEL fields for Console mode
			} else {
				// Set OTEL fields based on flags for OTEL/Fanout modes
				if hasOtelLoggerName {
					config.OtelLoggerName = "test-logger"
				}
				if hasOtelServiceName {
					config.OtelServiceName = "test-service"
				}
			}
			
			err := config.Validate()
			
			// Level is always required
			if !hasLevel {
				if err == nil {
					t.Logf("Expected error for missing Level, but got nil")
					return false
				}
				return true
			}
			
			// Output mode is always required
			if output == (OutputType{}) {
				if err == nil {
					t.Logf("Expected error for missing Output, but got nil")
					return false
				}
				return true
			}
			
			// For OTEL and Fanout modes, OTEL fields are required
			if output.String() == OutputOtel.String() || output.String() == OutputFanout.String() {
				if !hasOtelLoggerName || !hasOtelServiceName {
					if err == nil {
						t.Logf("Expected error for missing OTEL fields with output=%s, but got nil", output.String())
						return false
					}
					return true
				}
			}
			
			// If all required fields are present, should not error
			return err == nil
		},
		genOutput(),
		gen.Bool(),
		gen.Bool(),
		gen.Bool(),
	))

	properties.TestingRun(t)
}
