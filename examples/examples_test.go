package examples_test

import (
	"bytes"
	"os/exec"
	"testing"
	"time"
)

// TestAllExamplesCompile verifies that all examples compile successfully
func TestAllExamplesCompile(t *testing.T) {
	examples := []string{
		"github.com/wasilak/loggergo/examples/simple",
		"github.com/wasilak/loggergo/examples/otel",
		"github.com/wasilak/loggergo/examples/context",
		"github.com/wasilak/loggergo/examples/advanced",
	}

	for _, example := range examples {
		t.Run(example, func(t *testing.T) {
			cmd := exec.Command("go", "build", "-o", "/dev/null", example)
			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			
			if err := cmd.Run(); err != nil {
				t.Fatalf("example failed to compile: %v\nstderr: %s", err, stderr.String())
			}
		})
	}
}

// TestSimpleExampleRuns verifies that the simple example runs without errors
func TestSimpleExampleRuns(t *testing.T) {
	cmd := exec.Command("go", "run", "github.com/wasilak/loggergo/examples/simple",
		"-log-level=info", "-log-format=json")
	
	// Set a timeout to prevent hanging
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()
	
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("simple example failed to run: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("simple example timed out")
	}
}

// TestContextExampleRuns verifies that the context example runs without errors
func TestContextExampleRuns(t *testing.T) {
	cmd := exec.Command("go", "run", "github.com/wasilak/loggergo/examples/context")
	
	// Set a timeout to prevent hanging
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()
	
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("context example failed to run: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("context example timed out")
	}
}

// TestAdvancedExampleRuns verifies that the advanced example runs without errors
func TestAdvancedExampleRuns(t *testing.T) {
	cmd := exec.Command("go", "run", "github.com/wasilak/loggergo/examples/advanced")
	
	// Set a timeout to prevent hanging
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()
	
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("advanced example failed to run: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("advanced example timed out")
	}
}

// TestOtelExampleRuns verifies that the OTEL example runs without errors
// Note: This may fail gracefully if OTEL endpoint is not configured
func TestOtelExampleRuns(t *testing.T) {
	cmd := exec.Command("go", "run", "github.com/wasilak/loggergo/examples/otel")
	
	// Set a timeout to prevent hanging
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()
	
	select {
	case err := <-done:
		if err != nil {
			// OTEL example may fail if endpoint is not configured, which is expected
			t.Logf("otel example exited with error (expected if OTEL not configured): %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("otel example timed out")
	}
}
