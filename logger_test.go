package loggergo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"testing"
)

func TestLoggerInit_SetAsDefault(t *testing.T) {
	ctx := context.Background()

	// Create a buffer to capture log output
	var buf bytes.Buffer

	// Configure the logger to use the buffer
	config := Config{
		OutputStream: &buf,
		Output:       OutputConsole,
		SetAsDefault: true,
		Format:       LogFormatJSON,
	}

	// Initialize logger
	_, err := LoggerInit(ctx, config)
	if err != nil {
		t.Fatalf("Logger initialization failed: %v", err)
	}

	expectedMsgs := []string{"Test message"}

	// Log using slog.Default() and verify the output
	slog.Default().InfoContext(ctx, "Test message")
	slog.Default().InfoContext(ctx, "Test message")

	// Split the buffer by newlines to handle multiple JSON log entries
	logLines := bytes.Split(buf.Bytes(), []byte("\n"))

	// Iterate over each log line and check the contents
	for _, line := range logLines {
		if len(line) == 0 {
			continue // Skip empty lines
		}

		var logEntry map[string]interface{}
		if err := json.Unmarshal(line, &logEntry); err != nil {
			t.Fatalf("Log output is not valid JSON: %v", err)
		}

		// Check if the 'msg' field is in the list of expected msgs
		if v, ok := logEntry["msg"]; ok {
			found := false
			for _, msg := range expectedMsgs {
				if v.(string) == msg {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected 'msg' field to be one of '%v', but got: '%s'", expectedMsgs, v)
			}
		}

	}

}

func TestLoggerInit_SetAsLogLevelInfo(t *testing.T) {
	ctx := context.Background()

	// Create a buffer to capture log output
	var buf bytes.Buffer

	// Configure the logger to use the buffer
	config := Config{
		OutputStream: &buf,
		Output:       OutputConsole,
		SetAsDefault: true,
		Format:       LogFormatJSON,
		Level:        slog.LevelInfo,
	}

	// Initialize logger
	_, err := LoggerInit(ctx, config)
	if err != nil {
		t.Fatalf("Logger initialization failed: %v", err)
	}

	expectedMsgs := []string{"Test message"}
	NotExpectedMsgs := []string{"Debug message"}

	expectedLevels := []string{"INFO"}

	// Log using slog.Default() and verify the output
	slog.Default().InfoContext(ctx, expectedMsgs[0])
	slog.Default().DebugContext(ctx, NotExpectedMsgs[0])

	// Split the buffer by newlines to handle multiple JSON log entries
	logLines := bytes.Split(buf.Bytes(), []byte("\n"))

	// Iterate over each log line and check the contents
	for _, line := range logLines {
		if len(line) == 0 {
			continue // Skip empty lines
		}

		var logEntry map[string]interface{}
		if err := json.Unmarshal(line, &logEntry); err != nil {
			t.Fatalf("Log output is not valid JSON: %v", err)
		}

		// Check if the 'level' field is in the list of expected levels
		if v, ok := logEntry["level"]; ok {
			found := false
			for _, level := range expectedLevels {
				if v.(string) == level {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected 'level' field to be one of '%v', but got: '%s'", expectedLevels, v)
			}
		}

		// Check if the 'msg' field is in the list of expected msgs
		if v, ok := logEntry["msg"]; ok {
			found := false
			for _, msg := range expectedMsgs {
				if v.(string) == msg {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected 'msg' field to be one of '%v', but got: '%s'", expectedMsgs, v)
			}
		}

		// Check if the 'msg' field is in the list of NOT expected msgs
		if v, ok := logEntry["msg"]; ok {
			found := true
			for _, msg := range NotExpectedMsgs {
				if v.(string) == msg {
					found = false
					break
				}
			}
			if !found {
				t.Errorf("Expected 'msg' field to be not one of '%v', but got: '%s'", NotExpectedMsgs, v)
			}
		}

	}
}

func TestLoggerInit_SetAsLogLevelDebug(t *testing.T) {
	ctx := context.Background()

	// Create a buffer to capture log output
	var buf bytes.Buffer

	// Configure the logger to use the buffer
	config := Config{
		OutputStream: &buf,
		Output:       OutputConsole,
		SetAsDefault: true,
		Format:       LogFormatJSON,
		Level:        slog.LevelDebug,
	}

	// Initialize logger
	_, err := LoggerInit(ctx, config)
	if err != nil {
		t.Fatalf("Logger initialization failed: %v", err)
	}

	expectedMsgs := []string{"Test message", "Debug message"}

	expectedLevels := []string{"DEBUG", "INFO"}

	// Log using slog.Default() and verify the output
	slog.Default().InfoContext(ctx, expectedMsgs[0])
	slog.Default().DebugContext(ctx, expectedMsgs[1])

	// Split the buffer by newlines to handle multiple JSON log entries
	logLines := bytes.Split(buf.Bytes(), []byte("\n"))

	// Iterate over each log line and check the contents
	for _, line := range logLines {
		if len(line) == 0 {
			continue // Skip empty lines
		}

		var logEntry map[string]interface{}
		if err := json.Unmarshal(line, &logEntry); err != nil {
			t.Fatalf("Log output is not valid JSON: %v", err)
		}

		// Check if the 'level' field is in the list of expected levels
		if v, ok := logEntry["level"]; ok {
			found := false
			for _, level := range expectedLevels {
				if v.(string) == level {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected 'level' field to be one of '%v', but got: '%s'", expectedLevels, v)
			}
		}

		// Check if the 'msg' field is in the list of expected msgs
		if v, ok := logEntry["msg"]; ok {
			found := false
			for _, msg := range expectedMsgs {
				if v.(string) == msg {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected 'msg' field to be one of '%v', but got: '%s'", expectedMsgs, v)
			}
		}

	}
}

func TestLoggerInit_SetAsDefault_PlainText(t *testing.T) {
	ctx := context.Background()

	// Create a buffer to capture log output
	var buf bytes.Buffer

	// Configure the logger to use the buffer with plain text format
	config := Config{
		OutputStream: &buf,
		Output:       OutputConsole,
		SetAsDefault: true,
		Format:       LogFormatText, // Use plain text format
	}

	// Initialize logger
	_, err := LoggerInit(ctx, config)
	if err != nil {
		t.Fatalf("Logger initialization failed: %v", err)
	}

	expectedMsgs := []string{"Test message", "Test message"}

	// Log using slog.Default() and verify the output
	for _, msg := range expectedMsgs {
		slog.Default().InfoContext(ctx, msg)
	}

	// Split the buffer by newlines to handle multiple log entries
	logLines := bytes.Split(buf.Bytes(), []byte("\n"))

	// Iterate over each log line and check the contents
	for i, line := range logLines {
		if len(line) == 0 {
			continue // Skip empty lines
		}

		logLine := string(line)

		// Check for each expected message
		if !bytes.Contains(line, []byte(fmt.Sprintf("msg=%q", expectedMsgs[i]))) {
			t.Errorf("Expected message '%s' in log output, but got: %s", expectedMsgs[i], logLine)
		}

		// Check that the level is correctly set
		if !bytes.Contains(line, []byte("level=INFO")) {
			t.Errorf("Expected 'level=INFO' in log output, but got: %s", logLine)
		}

		// Check if the time field is present
		if !bytes.Contains(line, []byte("time=")) {
			t.Errorf("Expected 'time' field in log output, but got: %s", logLine)
		}
	}
}

func TestLoggerInit_SetAsDefault_OTEL(t *testing.T) {
	ctx := context.Background()

	// Create a buffer to capture log output
	var buf bytes.Buffer

	// Configure the logger to use the buffer with plain text format
	config := Config{
		OutputStream: &buf,
		Output:       OutputConsole,
		SetAsDefault: true,
		Format:       LogFormatOtel,
	}

	// Initialize logger
	_, err := LoggerInit(ctx, config)
	if err != nil {
		t.Fatalf("Logger initialization failed: %v", err)
	}

	expectedMsgs := []string{"Test message", "Test message"}

	// Log using slog.Default() and verify the output
	for _, msg := range expectedMsgs {
		slog.Default().InfoContext(ctx, msg)
	}

	// Split the buffer by newlines to handle multiple log entries
	logLines := bytes.Split(buf.Bytes(), []byte("\n"))

	// Iterate over each log line and check the contents
	for i, line := range logLines {
		if len(line) == 0 {
			continue // Skip empty lines
		}

		logLine := string(line)

		// Check for each expected message
		if !bytes.Contains(line, []byte(fmt.Sprintf("msg=%q", expectedMsgs[i]))) {
			t.Errorf("Expected message '%s' in log output, but got: %s", expectedMsgs[i], logLine)
		}

		// Check that the level is correctly set
		if !bytes.Contains(line, []byte("level=INFO")) {
			t.Errorf("Expected 'level=INFO' in log output, but got: %s", logLine)
		}

		// Check if the time field is present
		if !bytes.Contains(line, []byte("time=")) {
			t.Errorf("Expected 'time' field in log output, but got: %s", logLine)
		}
	}
}
