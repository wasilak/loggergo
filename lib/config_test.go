package lib

import (
	"fmt"
	"log/slog"
	"sync"
	"testing"

	"github.com/wasilak/loggergo/lib/types"
)

// TestConcurrentGetConfig tests concurrent GetConfig calls
func TestConcurrentGetConfig(t *testing.T) {
	// Initialize with a known config
	InitConfig()
	
	var wg sync.WaitGroup
	numGoroutines := 100
	
	// Launch multiple goroutines that read config concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			// Perform multiple read operations
			for j := 0; j < 100; j++ {
				config := GetConfig()
				
				// Verify we got a valid config
				if config.Level == nil {
					t.Error("GetConfig returned config with nil Level")
				}
			}
		}()
	}
	
	wg.Wait()
}

// TestConcurrentSetConfig tests concurrent SetConfig calls
func TestConcurrentSetConfig(t *testing.T) {
	// Initialize with a known config
	InitConfig()
	
	var wg sync.WaitGroup
	numGoroutines := 50
	
	// Launch multiple goroutines that write config concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			
			// Create a unique config for this goroutine
			config := types.Config{
				Level:           slog.LevelInfo,
				Format:          types.LogFormatJSON,
				Output:          types.OutputConsole,
				OtelLoggerName:  "test-logger",
				OtelServiceName: "test-service",
			}
			
			// Perform multiple write operations
			for j := 0; j < 10; j++ {
				SetConfig(config)
			}
		}(i)
	}
	
	wg.Wait()
	
	// Verify we can still read the config after concurrent writes
	config := GetConfig()
	if config.Level == nil {
		t.Error("Config has nil Level after concurrent writes")
	}
}

// TestMixedReadWriteOperations tests concurrent reads and writes
func TestMixedReadWriteOperations(t *testing.T) {
	// Initialize with a known config
	InitConfig()
	
	var wg sync.WaitGroup
	numReaders := 50
	numWriters := 25
	
	// Launch reader goroutines
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			for j := 0; j < 100; j++ {
				config := GetConfig()
				
				// Verify we got a valid config
				if config.Level == nil {
					t.Error("GetConfig returned config with nil Level during mixed operations")
				}
			}
		}()
	}
	
	// Launch writer goroutines
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			
			config := types.Config{
				Level:           slog.LevelInfo,
				Format:          types.LogFormatJSON,
				Output:          types.OutputConsole,
				OtelLoggerName:  "test-logger",
				OtelServiceName: "test-service",
			}
			
			for j := 0; j < 50; j++ {
				SetConfig(config)
			}
		}(i)
	}
	
	wg.Wait()
	
	// Final verification
	config := GetConfig()
	if config.Level == nil {
		t.Error("Config has nil Level after mixed read/write operations")
	}
}

// TestConcurrentMergeConfig tests concurrent MergeConfig calls
func TestConcurrentMergeConfig(t *testing.T) {
	// Initialize with a known config
	InitConfig()
	
	var wg sync.WaitGroup
	numGoroutines := 50
	
	// Launch multiple goroutines that merge config concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			
			override := types.Config{
				Level:           slog.LevelDebug,
				Output:          types.OutputConsole,
				OtelLoggerName:  "merged-logger",
				OtelServiceName: "merged-service",
			}
			
			for j := 0; j < 10; j++ {
				_ = MergeConfig(override)
			}
		}(i)
	}
	
	wg.Wait()
	
	// Verify we can still read the config after concurrent merges
	config := GetConfig()
	if config.Level == nil {
		t.Error("Config has nil Level after concurrent merges")
	}
}

// TestConcurrentInitConfig tests concurrent InitConfig calls
func TestConcurrentInitConfig(t *testing.T) {
	var wg sync.WaitGroup
	numGoroutines := 50
	
	// Launch multiple goroutines that initialize config concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			for j := 0; j < 10; j++ {
				InitConfig()
			}
		}()
	}
	
	wg.Wait()
	
	// Verify the config is in a valid state
	config := GetConfig()
	if config.Level == nil {
		t.Error("Config has nil Level after concurrent initializations")
	}
	if config.Output == (types.OutputType{}) {
		t.Error("Config has zero Output after concurrent initializations")
	}
}

// TestMergeConfig_AllZeroValues tests merging with all zero values in override
func TestMergeConfig_AllZeroValues(t *testing.T) {
	// Set a base config with non-zero values
	base := types.Config{
		Level:              slog.LevelInfo,
		Format:             types.LogFormatJSON,
		DevMode:            true,
		DevFlavor:          types.DevFlavorTint,
		OtelTracingEnabled: true,
		OtelLoggerName:     "base-logger",
		Output:             types.OutputConsole,
		OtelServiceName:    "base-service",
		SetAsDefault:       true,
		ContextKeys:        []interface{}{"key1", "key2"},
		ContextKeysDefault: "default-value",
	}
	SetConfig(base)
	
	// Merge with all zero values
	override := types.Config{}
	result := MergeConfig(override)
	
	// Non-boolean base values should be preserved (zero values ignored)
	if result.Level.Level() != slog.LevelInfo {
		t.Errorf("Expected Level to be Info, got %v", result.Level.Level())
	}
	if result.Format.String() != types.LogFormatJSON.String() {
		t.Errorf("Expected Format to be JSON, got %v", result.Format.String())
	}
	if result.DevFlavor.String() != types.DevFlavorTint.String() {
		t.Errorf("Expected DevFlavor to be Tint, got %v", result.DevFlavor.String())
	}
	if result.OtelLoggerName != "base-logger" {
		t.Errorf("Expected OtelLoggerName to be 'base-logger', got %v", result.OtelLoggerName)
	}
	if result.Output.String() != types.OutputConsole.String() {
		t.Errorf("Expected Output to be Console, got %v", result.Output.String())
	}
	if result.OtelServiceName != "base-service" {
		t.Errorf("Expected OtelServiceName to be 'base-service', got %v", result.OtelServiceName)
	}
	if len(result.ContextKeys) != 2 {
		t.Errorf("Expected 2 ContextKeys, got %d", len(result.ContextKeys))
	}
	if result.ContextKeysDefault != "default-value" {
		t.Errorf("Expected ContextKeysDefault to be 'default-value', got %v", result.ContextKeysDefault)
	}
	
	// Boolean fields: false (zero value) overrides true when different
	// This is expected behavior - we can't distinguish "not set" from "set to false"
	if result.DevMode {
		t.Error("Expected DevMode to be false (zero value overrides true)")
	}
	if result.OtelTracingEnabled {
		t.Error("Expected OtelTracingEnabled to be false (zero value overrides true)")
	}
	if result.SetAsDefault {
		t.Error("Expected SetAsDefault to be false (zero value overrides true)")
	}
}

// TestMergeConfig_AllNonZeroValues tests merging with all non-zero values in override
func TestMergeConfig_AllNonZeroValues(t *testing.T) {
	// Set a base config
	base := types.Config{
		Level:              slog.LevelInfo,
		Format:             types.LogFormatJSON,
		DevMode:            true,
		DevFlavor:          types.DevFlavorTint,
		OtelTracingEnabled: true,
		OtelLoggerName:     "base-logger",
		Output:             types.OutputConsole,
		OtelServiceName:    "base-service",
		SetAsDefault:       true,
		ContextKeys:        []interface{}{"key1"},
		ContextKeysDefault: "base-default",
	}
	SetConfig(base)
	
	// Merge with all non-zero values
	override := types.Config{
		Level:              slog.LevelDebug,
		Format:             types.LogFormatText,
		DevMode:            false,
		DevFlavor:          types.DevFlavorSlogor,
		OtelTracingEnabled: false,
		OtelLoggerName:     "override-logger",
		Output:             types.OutputOtel,
		OtelServiceName:    "override-service",
		SetAsDefault:       false,
		ContextKeys:        []interface{}{"key2", "key3"},
		ContextKeysDefault: "override-default",
	}
	result := MergeConfig(override)
	
	// All override values should be used
	if result.Level.Level() != slog.LevelDebug {
		t.Errorf("Expected Level to be Debug, got %v", result.Level.Level())
	}
	if result.Format.String() != types.LogFormatText.String() {
		t.Errorf("Expected Format to be Text, got %v", result.Format.String())
	}
	if result.DevMode {
		t.Error("Expected DevMode to be false")
	}
	if result.DevFlavor.String() != types.DevFlavorSlogor.String() {
		t.Errorf("Expected DevFlavor to be Slogor, got %v", result.DevFlavor.String())
	}
	if result.OtelTracingEnabled {
		t.Error("Expected OtelTracingEnabled to be false")
	}
	if result.OtelLoggerName != "override-logger" {
		t.Errorf("Expected OtelLoggerName to be 'override-logger', got %v", result.OtelLoggerName)
	}
	if result.Output.String() != types.OutputOtel.String() {
		t.Errorf("Expected Output to be Otel, got %v", result.Output.String())
	}
	if result.OtelServiceName != "override-service" {
		t.Errorf("Expected OtelServiceName to be 'override-service', got %v", result.OtelServiceName)
	}
	if result.SetAsDefault {
		t.Error("Expected SetAsDefault to be false")
	}
	if len(result.ContextKeys) != 2 {
		t.Errorf("Expected 2 ContextKeys, got %d", len(result.ContextKeys))
	}
	if result.ContextKeysDefault != "override-default" {
		t.Errorf("Expected ContextKeysDefault to be 'override-default', got %v", result.ContextKeysDefault)
	}
}

// TestMergeConfig_PartialOverrides tests merging with partial overrides
func TestMergeConfig_PartialOverrides(t *testing.T) {
	// Set a base config
	base := types.Config{
		Level:              slog.LevelInfo,
		Format:             types.LogFormatJSON,
		DevMode:            true,
		DevFlavor:          types.DevFlavorTint,
		OtelTracingEnabled: true,
		OtelLoggerName:     "base-logger",
		Output:             types.OutputConsole,
		OtelServiceName:    "base-service",
		SetAsDefault:       true,
	}
	SetConfig(base)
	
	// Merge with partial overrides (only some fields set)
	override := types.Config{
		Level:           slog.LevelDebug,
		OtelLoggerName:  "override-logger",
		DevMode:         false,
	}
	result := MergeConfig(override)
	
	// Override values should be used where provided
	if result.Level.Level() != slog.LevelDebug {
		t.Errorf("Expected Level to be Debug, got %v", result.Level.Level())
	}
	if result.OtelLoggerName != "override-logger" {
		t.Errorf("Expected OtelLoggerName to be 'override-logger', got %v", result.OtelLoggerName)
	}
	if result.DevMode {
		t.Error("Expected DevMode to be false")
	}
	
	// Base values should be preserved where not overridden (non-boolean fields)
	if result.Format.String() != types.LogFormatJSON.String() {
		t.Errorf("Expected Format to be JSON, got %v", result.Format.String())
	}
	if result.DevFlavor.String() != types.DevFlavorTint.String() {
		t.Errorf("Expected DevFlavor to be Tint, got %v", result.DevFlavor.String())
	}
	if result.Output.String() != types.OutputConsole.String() {
		t.Errorf("Expected Output to be Console, got %v", result.Output.String())
	}
	if result.OtelServiceName != "base-service" {
		t.Errorf("Expected OtelServiceName to be 'base-service', got %v", result.OtelServiceName)
	}
	
	// Boolean fields: false (zero value) overrides true when different
	// OtelTracingEnabled and SetAsDefault are not in override, so they get zero value (false)
	// which overrides the base true value
	if result.OtelTracingEnabled {
		t.Error("Expected OtelTracingEnabled to be false (zero value overrides true)")
	}
	if result.SetAsDefault {
		t.Error("Expected SetAsDefault to be false (zero value overrides true)")
	}
}

// TestMergeConfig_BooleanFieldMerging tests boolean field merge logic
func TestMergeConfig_BooleanFieldMerging(t *testing.T) {
	tests := []struct {
		name                string
		baseDevMode         bool
		baseOtelEnabled     bool
		baseSetAsDefault    bool
		overrideDevMode     bool
		overrideOtelEnabled bool
		overrideSetAsDefault bool
		expectDevMode       bool
		expectOtelEnabled   bool
		expectSetAsDefault  bool
	}{
		{
			name:                "false overrides true",
			baseDevMode:         true,
			baseOtelEnabled:     true,
			baseSetAsDefault:    true,
			overrideDevMode:     false,
			overrideOtelEnabled: false,
			overrideSetAsDefault: false,
			expectDevMode:       false,
			expectOtelEnabled:   false,
			expectSetAsDefault:  false,
		},
		{
			name:                "true overrides false",
			baseDevMode:         false,
			baseOtelEnabled:     false,
			baseSetAsDefault:    false,
			overrideDevMode:     true,
			overrideOtelEnabled: true,
			overrideSetAsDefault: true,
			expectDevMode:       true,
			expectOtelEnabled:   true,
			expectSetAsDefault:  true,
		},
		{
			name:                "same values preserved",
			baseDevMode:         true,
			baseOtelEnabled:     false,
			baseSetAsDefault:    true,
			overrideDevMode:     true,
			overrideOtelEnabled: false,
			overrideSetAsDefault: true,
			expectDevMode:       true,
			expectOtelEnabled:   false,
			expectSetAsDefault:  true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set base config
			base := types.Config{
				Level:              slog.LevelInfo,
				Format:             types.LogFormatJSON,
				Output:             types.OutputConsole,
				DevMode:            tt.baseDevMode,
				OtelTracingEnabled: tt.baseOtelEnabled,
				SetAsDefault:       tt.baseSetAsDefault,
				OtelLoggerName:     "test-logger",
				OtelServiceName:    "test-service",
			}
			SetConfig(base)
			
			// Merge with override
			override := types.Config{
				DevMode:            tt.overrideDevMode,
				OtelTracingEnabled: tt.overrideOtelEnabled,
				SetAsDefault:       tt.overrideSetAsDefault,
			}
			result := MergeConfig(override)
			
			// Verify boolean fields
			if result.DevMode != tt.expectDevMode {
				t.Errorf("Expected DevMode to be %v, got %v", tt.expectDevMode, result.DevMode)
			}
			if result.OtelTracingEnabled != tt.expectOtelEnabled {
				t.Errorf("Expected OtelTracingEnabled to be %v, got %v", tt.expectOtelEnabled, result.OtelTracingEnabled)
			}
			if result.SetAsDefault != tt.expectSetAsDefault {
				t.Errorf("Expected SetAsDefault to be %v, got %v", tt.expectSetAsDefault, result.SetAsDefault)
			}
		})
	}
}

// Example_mergeConfigBooleanBehavior demonstrates the boolean field behavior
// that users need to be aware of when using MergeConfig.
func Example_mergeConfigBooleanBehavior() {
	// Set up a base config with DevMode enabled
	base := types.Config{
		Level:              slog.LevelInfo,
		Format:             types.LogFormatJSON,
		Output:             types.OutputConsole,
		DevMode:            true,  // Base has DevMode enabled
		OtelTracingEnabled: true,
		SetAsDefault:       true,
		OtelLoggerName:     "test-logger",
		OtelServiceName:    "test-service",
	}
	SetConfig(base)
	
	// PROBLEM: Partial override without explicitly setting boolean fields
	// This will unintentionally change DevMode from true to false!
	problemOverride := types.Config{
		Level: slog.LevelDebug,  // Only want to change the level
		// DevMode not set - defaults to false and will override base's true!
	}
	problemResult := MergeConfig(problemOverride)
	
	// DevMode is now false (unexpected!)
	fmt.Printf("Problem - DevMode after merge: %v\n", problemResult.DevMode)
	
	// SOLUTION: Explicitly set boolean fields you want to preserve
	correctOverride := types.Config{
		Level:              slog.LevelDebug,  // Change the level
		DevMode:            true,             // Explicitly preserve DevMode
		OtelTracingEnabled: true,             // Explicitly preserve OtelTracingEnabled
		SetAsDefault:       true,             // Explicitly preserve SetAsDefault
	}
	
	// Reset base config
	SetConfig(base)
	correctResult := MergeConfig(correctOverride)
	
	// DevMode is preserved (expected!)
	fmt.Printf("Solution - DevMode after merge: %v\n", correctResult.DevMode)
	
	// Output:
	// Problem - DevMode after merge: false
	// Solution - DevMode after merge: true
}

// TestRegisterCleanup tests that cleanup functions are registered correctly
func TestRegisterCleanup(t *testing.T) {
	// Reset cleanup functions
	globalConfigManager.cleanupMu.Lock()
	globalConfigManager.cleanupFuncs = nil
	globalConfigManager.cleanupMu.Unlock()
	
	called := false
	cleanup := func() error {
		called = true
		return nil
	}
	
	RegisterCleanup(cleanup)
	
	// Verify cleanup function was registered
	globalConfigManager.cleanupMu.Lock()
	count := len(globalConfigManager.cleanupFuncs)
	globalConfigManager.cleanupMu.Unlock()
	
	if count != 1 {
		t.Errorf("Expected 1 cleanup function, got %d", count)
	}
	
	// Call Shutdown to execute cleanup
	err := Shutdown()
	if err != nil {
		t.Errorf("Shutdown returned error: %v", err)
	}
	
	if !called {
		t.Error("Cleanup function was not called")
	}
}

// TestShutdown_MultipleCleanups tests that multiple cleanup functions are called in LIFO order
func TestShutdown_MultipleCleanups(t *testing.T) {
	// Reset cleanup functions
	globalConfigManager.cleanupMu.Lock()
	globalConfigManager.cleanupFuncs = nil
	globalConfigManager.cleanupMu.Unlock()
	
	var order []int
	
	// Register cleanup functions
	RegisterCleanup(func() error {
		order = append(order, 1)
		return nil
	})
	RegisterCleanup(func() error {
		order = append(order, 2)
		return nil
	})
	RegisterCleanup(func() error {
		order = append(order, 3)
		return nil
	})
	
	// Call Shutdown
	err := Shutdown()
	if err != nil {
		t.Errorf("Shutdown returned error: %v", err)
	}
	
	// Verify LIFO order (3, 2, 1)
	if len(order) != 3 {
		t.Fatalf("Expected 3 cleanup calls, got %d", len(order))
	}
	if order[0] != 3 || order[1] != 2 || order[2] != 1 {
		t.Errorf("Expected cleanup order [3, 2, 1], got %v", order)
	}
}

// TestShutdown_WithErrors tests that Shutdown continues even when cleanup functions return errors
func TestShutdown_WithErrors(t *testing.T) {
	// Reset cleanup functions
	globalConfigManager.cleanupMu.Lock()
	globalConfigManager.cleanupFuncs = nil
	globalConfigManager.cleanupMu.Unlock()
	
	var called []int
	
	// Register cleanup functions, some with errors
	RegisterCleanup(func() error {
		called = append(called, 1)
		return nil
	})
	RegisterCleanup(func() error {
		called = append(called, 2)
		return fmt.Errorf("error from cleanup 2")
	})
	RegisterCleanup(func() error {
		called = append(called, 3)
		return fmt.Errorf("error from cleanup 3")
	})
	
	// Call Shutdown
	err := Shutdown()
	
	// Should return an error
	if err == nil {
		t.Error("Expected Shutdown to return error, got nil")
	}
	
	// All cleanup functions should have been called despite errors
	if len(called) != 3 {
		t.Errorf("Expected 3 cleanup calls, got %d", len(called))
	}
	
	// Verify LIFO order (3, 2, 1)
	if called[0] != 3 || called[1] != 2 || called[2] != 1 {
		t.Errorf("Expected cleanup order [3, 2, 1], got %v", called)
	}
}

// TestShutdown_ClearsCleanupFunctions tests that cleanup functions are cleared after Shutdown
func TestShutdown_ClearsCleanupFunctions(t *testing.T) {
	// Reset cleanup functions
	globalConfigManager.cleanupMu.Lock()
	globalConfigManager.cleanupFuncs = nil
	globalConfigManager.cleanupMu.Unlock()
	
	// Register a cleanup function
	RegisterCleanup(func() error {
		return nil
	})
	
	// Verify it was registered
	globalConfigManager.cleanupMu.Lock()
	count := len(globalConfigManager.cleanupFuncs)
	globalConfigManager.cleanupMu.Unlock()
	
	if count != 1 {
		t.Errorf("Expected 1 cleanup function before Shutdown, got %d", count)
	}
	
	// Call Shutdown
	err := Shutdown()
	if err != nil {
		t.Errorf("Shutdown returned error: %v", err)
	}
	
	// Verify cleanup functions were cleared
	globalConfigManager.cleanupMu.Lock()
	count = len(globalConfigManager.cleanupFuncs)
	globalConfigManager.cleanupMu.Unlock()
	
	if count != 0 {
		t.Errorf("Expected 0 cleanup functions after Shutdown, got %d", count)
	}
}

// TestConcurrentRegisterCleanup tests concurrent RegisterCleanup calls
func TestConcurrentRegisterCleanup(t *testing.T) {
	// Reset cleanup functions
	globalConfigManager.cleanupMu.Lock()
	globalConfigManager.cleanupFuncs = nil
	globalConfigManager.cleanupMu.Unlock()
	
	var wg sync.WaitGroup
	numGoroutines := 50
	
	// Launch multiple goroutines that register cleanup functions concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			
			RegisterCleanup(func() error {
				return nil
			})
		}(i)
	}
	
	wg.Wait()
	
	// Verify all cleanup functions were registered
	globalConfigManager.cleanupMu.Lock()
	count := len(globalConfigManager.cleanupFuncs)
	globalConfigManager.cleanupMu.Unlock()
	
	if count != numGoroutines {
		t.Errorf("Expected %d cleanup functions, got %d", numGoroutines, count)
	}
	
	// Clean up
	Shutdown()
}

// TestShutdown_EmptyCleanupList tests that Shutdown works with no cleanup functions
func TestShutdown_EmptyCleanupList(t *testing.T) {
	// Reset cleanup functions
	globalConfigManager.cleanupMu.Lock()
	globalConfigManager.cleanupFuncs = nil
	globalConfigManager.cleanupMu.Unlock()
	
	// Call Shutdown with no cleanup functions
	err := Shutdown()
	if err != nil {
		t.Errorf("Shutdown returned error with empty cleanup list: %v", err)
	}
}
