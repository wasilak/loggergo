package loggergo

import (
	"context"
	"log/slog"
)

// CustomContextAttributeHandler wraps an existing handler and adds a custom attribute to all logs.
type CustomContextAttributeHandler struct {
	innerHandler       slog.Handler
	keys               []string
	ContextKeysDefault interface{}
	CTX                context.Context
}

// NewCustomContextAttributeHandler creates a new handler that wraps the given handler and adds a custom attribute.
func NewCustomContextAttributeHandler(ctx context.Context, handler slog.Handler, keys []string, contextKeysDefault interface{}) *CustomContextAttributeHandler {
	return &CustomContextAttributeHandler{
		innerHandler:       handler,
		keys:               keys,
		ContextKeysDefault: contextKeysDefault,
		CTX:                ctx,
	}
}

// Enabled delegates the Enabled check to the inner handler.
func (h *CustomContextAttributeHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.innerHandler.Enabled(ctx, level)
}

// Handle adds the custom attribute and delegates the log processing to the inner handler.
func (h *CustomContextAttributeHandler) Handle(ctx context.Context, record slog.Record) error {

	for _, key := range h.keys {

		val := h.CTX.Value(key)

		if val == nil {
			if h.ContextKeysDefault != nil {
				record.Add(slog.Any(key, h.ContextKeysDefault))
			}
		} else {
			record.Add(slog.Any(key, val))
		}
	}

	// Delegate to the inner handler
	return h.innerHandler.Handle(ctx, record)
}

// WithAttrs creates a new handler with additional attributes, preserving the custom attribute.
func (h *CustomContextAttributeHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewCustomContextAttributeHandler(h.CTX, h.innerHandler.WithAttrs(attrs), h.keys, h.ContextKeysDefault)
}

// WithGroup creates a new handler with a group, preserving the custom attribute.
func (h *CustomContextAttributeHandler) WithGroup(name string) slog.Handler {
	return NewCustomContextAttributeHandler(h.CTX, h.innerHandler.WithGroup(name), h.keys, h.ContextKeysDefault)
}
