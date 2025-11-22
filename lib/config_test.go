package lib

import (
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
