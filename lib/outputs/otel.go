package outputs

import (
	"context"
	"log/slog"

	"github.com/wasilak/loggergo/lib"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	otellog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

const sevOffset = slog.Level(otellog.SeverityDebug) - slog.LevelDebug

// Custom processor to filter logs by level
type levelFilterProcessor struct {
	minLevel  slog.Level
	processor log.Processor
}

// OnEmit filters log records by level and delegates to the wrapped processor
func (p *levelFilterProcessor) OnEmit(ctx context.Context, record *log.Record) error {
	sev := slog.Level(record.Severity()) - sevOffset

	if sev >= p.minLevel {
		return p.processor.OnEmit(ctx, record)
	}
	return nil // Ignore logs below the minimum level
}

// Shutdown cleans up resources used by the processor
func (p *levelFilterProcessor) Shutdown(ctx context.Context) error {
	return p.processor.Shutdown(ctx)
}

// ForceFlush ensures all pending records are flushed
func (p *levelFilterProcessor) ForceFlush(ctx context.Context) error {
	return p.processor.ForceFlush(ctx)
}

// setupOtelFormat sets up a slog.Handler for OpenTelemetry format.
// It merges the default resource with the service name attribute, creates a stdoutlog exporter,
// and sets up a log processor and logger provider with the merged resource and exporter.
// Returns the handler and any error encountered.
func SetupOtelFormat() (slog.Handler, error) {
	resource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(lib.GetConfig().OtelServiceName),
		),
	)
	if err != nil {
		return nil, err
	}

	exporter, err := stdoutlog.New()
	if err != nil {
		return nil, err
	}

	// Wrap the exporter with a simple processor
	baseProcessor := log.NewSimpleProcessor(exporter)

	// Wrap the processor with a level filter
	filteredProcessor := &levelFilterProcessor{
		minLevel:  lib.GetConfig().Level.Level(),
		processor: baseProcessor,
	}

	// processor := log.NewSimpleProcessor(exporter)
	stdoutProvider := log.NewLoggerProvider(
		log.WithResource(resource),
		// log.WithProcessor(processor),
		log.WithProcessor(filteredProcessor),
	)

	return otelslog.NewHandler(lib.GetConfig().OtelLoggerName, otelslog.WithLoggerProvider(stdoutProvider)), nil
}
