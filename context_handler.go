package loggergo

import (
	"context"
	"fmt"
	"log/slog"
)

// CustomContextAttributeHandler wraps an existing handler and adds a custom attribute to all logs.
type CustomContextAttributeHandler struct {
	innerHandler       slog.Handler
	keys               []interface{}
	ContextKeysDefault interface{}
}

// NewCustomContextAttributeHandler creates a new handler that wraps the given handler and adds a custom attribute.
func NewCustomContextAttributeHandler(handler slog.Handler, keys []interface{}, contextKeysDefault interface{}) *CustomContextAttributeHandler {
	return &CustomContextAttributeHandler{
		innerHandler:       handler,
		keys:               keys,
		ContextKeysDefault: contextKeysDefault,
	}
}

// Enabled delegates the Enabled check to the inner handler.
func (h *CustomContextAttributeHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.innerHandler.Enabled(ctx, level)
}

// Handle adds the custom attribute and delegates the log processing to the inner handler.
func (h *CustomContextAttributeHandler) Handle(ctx context.Context, record slog.Record) error {

	for _, key := range h.keys {

		val := ctx.Value(key)

		if val == nil {
			if h.ContextKeysDefault != nil {
				record.Add(slog.Any(fmt.Sprintf("%v", key), h.ContextKeysDefault))
			}
		} else {
			record.Add(slog.Any(fmt.Sprintf("%v", key), val))
		}
	}

	// Delegate to the inner handler
	return h.innerHandler.Handle(ctx, record)
}

// WithAttrs creates a new handler with additional attributes, preserving the custom attribute.
func (h *CustomContextAttributeHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewCustomContextAttributeHandler(h.innerHandler.WithAttrs(attrs), h.keys, h.ContextKeysDefault)
}

// WithGroup creates a new handler with a group, preserving the custom attribute.
func (h *CustomContextAttributeHandler) WithGroup(name string) slog.Handler {
	return NewCustomContextAttributeHandler(h.innerHandler.WithGroup(name), h.keys, h.ContextKeysDefault)
}
