package loggergo

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

// CustomContextAttributeHandler wraps an existing slog.Handler and automatically extracts
// values from context.Context to add as attributes to all log records.
//
// This handler enables context-aware logging by extracting specified keys from the context
// and including them in every log entry. It handles nil contexts, missing keys, and type
// conversion gracefully.
//
// Thread Safety:
//
// CustomContextAttributeHandler is safe for concurrent use. Multiple goroutines can call
// Handle, Enabled, WithAttrs, and WithGroup simultaneously without additional synchronization.
type CustomContextAttributeHandler struct {
	innerHandler       slog.Handler
	keys               []interface{}
	ContextKeysDefault interface{}
}

// NewCustomContextAttributeHandler creates a new handler that wraps the given handler
// and automatically extracts context values.
//
// Parameters:
//   - handler: The underlying slog.Handler to wrap
//   - keys: Slice of context keys to extract from context.Context
//   - contextKeysDefault: Default value to use when a key is not found in context (can be nil)
//
// Returns:
//   - *CustomContextAttributeHandler: A new handler that extracts context values
//
// Example:
//
//	type contextKey string
//	const requestIDKey contextKey = "request_id"
//
//	handler := slog.NewJSONHandler(os.Stdout, nil)
//	contextHandler := NewCustomContextAttributeHandler(
//	    handler,
//	    []interface{}{requestIDKey},
//	    "unknown",
//	)
//	logger := slog.New(contextHandler)
//
//	ctx := context.WithValue(context.Background(), requestIDKey, "req-123")
//	logger.InfoContext(ctx, "Processing request") // Will include request_id: "req-123"
func NewCustomContextAttributeHandler(handler slog.Handler, keys []interface{}, contextKeysDefault interface{}) *CustomContextAttributeHandler {
	return &CustomContextAttributeHandler{
		innerHandler:       handler,
		keys:               keys,
		ContextKeysDefault: contextKeysDefault,
	}
}

// Enabled reports whether the handler handles records at the given level.
// It delegates the check to the inner handler.
func (h *CustomContextAttributeHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.innerHandler.Enabled(ctx, level)
}

// Handle processes a log record by extracting context values and delegating to the inner handler.
//
// It extracts values for all configured keys from the context and adds them as attributes
// to the log record. If a key is not found, it uses the default value (if configured) or
// omits the field.
//
// Error Handling:
//
// Handle never panics. All panics are recovered and returned as errors.
// If the context is nil, context.Background() is used as a fallback.
//
// Thread Safety:
//
// Handle is safe to call concurrently from multiple goroutines.
func (h *CustomContextAttributeHandler) Handle(ctx context.Context, record slog.Record) (err error) {
	// Panic recovery to ensure Handle never panics
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in Handle: %v", r)
			// Log the panic to stderr for debugging
			fmt.Fprintf(os.Stderr, "PANIC recovered in Handle(): %v\n", r)
		}
	}()

	// Handle nil context by using context.Background() as fallback
	if ctx == nil {
		ctx = context.Background()
	}

	for _, key := range h.keys {
		// Safe context value extraction with error handling
		val := ctx.Value(key)

		if val == nil {
			// Use default value for missing keys
			if h.ContextKeysDefault != nil {
				record.Add(slog.Any(fmt.Sprintf("%v", key), h.ContextKeysDefault))
			}
			// If no default is set, omit the field (graceful handling)
		} else {
			// Add the extracted value to the log record
			// slog.Any handles type formatting appropriately
			record.Add(slog.Any(fmt.Sprintf("%v", key), val))
		}
	}

	// Delegate to the inner handler
	return h.innerHandler.Handle(ctx, record)
}

// WithAttrs returns a new handler with the given attributes added.
// The new handler preserves the context extraction behavior.
func (h *CustomContextAttributeHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewCustomContextAttributeHandler(h.innerHandler.WithAttrs(attrs), h.keys, h.ContextKeysDefault)
}

// WithGroup returns a new handler with the given group name.
// The new handler preserves the context extraction behavior.
func (h *CustomContextAttributeHandler) WithGroup(name string) slog.Handler {
	return NewCustomContextAttributeHandler(h.innerHandler.WithGroup(name), h.keys, h.ContextKeysDefault)
}
