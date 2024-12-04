package loggergo

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/golang-cz/devslog"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	otellogs "github.com/wasilak/otelgo/logs"
	otelgoslog "github.com/wasilak/otelgo/slog"
	"gitlab.com/greyxor/slogor"
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

// consoleMode returns a slog.Handler based on the provided defaultConfig and opts.
// It checks the defaultConfig.Format and sets up the appropriate handler based on the format.
// If defaultConfig.OtelTracingEnabled is true, it wraps the handler with otelgoslog.NewTracingHandler.
// Returns the handler and any error encountered.
func consoleMode(defaultConfig Config, opts slog.HandlerOptions) (slog.Handler, error) {
	var handler slog.Handler
	var err error

	if defaultConfig.Format == LogFormatOtel {
		return setupOtelFormat(defaultConfig)
	}

	if defaultConfig.Format == LogFormatJSON {
		handler = slog.NewJSONHandler(defaultConfig.OutputStream, &opts)
	}

	if defaultConfig.Format == LogFormatText {
		handler, err = setupPlainFormat(opts, defaultConfig)
		if err != nil {
			return nil, err
		}
	}

	if defaultConfig.OtelTracingEnabled {
		handler = otelgoslog.NewTracingHandler(handler)
	}

	return handler, nil
}

// otelMode returns a slog.Handler for OpenTelemetry mode based on the provided defaultConfig.
// It initializes the otellogs package and returns a handler with the otelslog.WithLoggerProvider option.
// Returns the handler and any error encountered.
func otelMode(ctx context.Context, defaultConfig Config) (slog.Handler, context.Context, error) {
	otelGoLogsConfig := otellogs.OtelGoLogsConfig{}

	ctx, provider, err := otellogs.Init(ctx, otelGoLogsConfig)
	if err != nil {
		return nil, ctx, err
	}

	return otelslog.NewHandler(defaultConfig.OtelLoggerName, otelslog.WithLoggerProvider(provider)), ctx, nil
}

// setupOtelFormat sets up a slog.Handler for OpenTelemetry format.
// It merges the default resource with the service name attribute, creates a stdoutlog exporter,
// and sets up a log processor and logger provider with the merged resource and exporter.
// Returns the handler and any error encountered.
func setupOtelFormat(defaultConfig Config) (slog.Handler, error) {
	resource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(defaultConfig.OtelServiceName),
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
		minLevel:  defaultConfig.Level.Level(),
		processor: baseProcessor,
	}

	// processor := log.NewSimpleProcessor(exporter)
	stdoutProvider := log.NewLoggerProvider(
		log.WithResource(resource),
		// log.WithProcessor(processor),
		log.WithProcessor(filteredProcessor),
	)

	return otelslog.NewHandler(defaultConfig.OtelLoggerName, otelslog.WithLoggerProvider(stdoutProvider)), nil
}

// setupPlainFormat sets up a slog.Handler for plain format.
// If defaultConfig.DevMode is true, it checks the defaultConfig.DevFlavor and sets up the appropriate handler based on the flavor.
// Returns the handler and any error encountered.
func setupPlainFormat(opts slog.HandlerOptions, defaultConfig Config) (slog.Handler, error) {
	if defaultConfig.DevMode {

		if defaultConfig.DevFlavor == DevFlavorSlogor {
			return slogor.NewHandler(defaultConfig.OutputStream, slogor.ShowSource(), slogor.SetTimeFormat(time.Stamp), slogor.SetLevel(opts.Level.Level())), nil
		} else if defaultConfig.DevFlavor == DevFlavorDevslog {
			return devslog.NewHandler(defaultConfig.OutputStream, &devslog.Options{
				HandlerOptions:    &opts,
				MaxSlicePrintSize: 10,
				SortKeys:          true,
			}), nil
		} else {
			return tint.NewHandler(defaultConfig.OutputStream, &tint.Options{
				Level:     opts.Level,
				NoColor:   !isatty.IsTerminal(os.Stderr.Fd()),
				AddSource: opts.AddSource,
			}), nil
		}
	}

	return slog.NewTextHandler(defaultConfig.OutputStream, &opts), nil
}
