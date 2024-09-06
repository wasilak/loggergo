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
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

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
func otelMode(ctx context.Context, defaultConfig Config) (slog.Handler, error) {
	otelGoLogsConfig := otellogs.OtelGoLogsConfig{}

	_, provider, err := otellogs.Init(ctx, otelGoLogsConfig)
	if err != nil {
		return nil, err
	}

	return otelslog.NewHandler(defaultConfig.OtelLoggerName, otelslog.WithLoggerProvider(provider)), nil
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

	processor := log.NewSimpleProcessor(exporter)
	stdoutProvider := log.NewLoggerProvider(
		log.WithResource(resource),
		log.WithProcessor(processor),
	)

	return otelslog.NewHandler(defaultConfig.OtelLoggerName, otelslog.WithLoggerProvider(stdoutProvider)), nil
}

// setupPlainFormat sets up a slog.Handler for plain format.
// If defaultConfig.DevMode is true, it checks the defaultConfig.DevFlavor and sets up the appropriate handler based on the flavor.
// Returns the handler and any error encountered.
func setupPlainFormat(opts slog.HandlerOptions, defaultConfig Config) (slog.Handler, error) {
	if defaultConfig.DevMode {

		if defaultConfig.DevFlavor == DevFlavorSlogor {
			return slogor.NewHandler(defaultConfig.OutputStream, slogor.Options{
				TimeFormat: time.Stamp,
				Level:      opts.Level.Level(),
				ShowSource: opts.AddSource,
			}), nil
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
