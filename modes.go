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
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

func consoleMode(defaultConfig LoggerGoConfig, opts slog.HandlerOptions) (slog.Handler, error) {
	var handler slog.Handler
	var err error

	if defaultConfig.Format == "otel" {
		return setupOtelFormat()
	}

	if defaultConfig.Format == "json" {
		handler = slog.NewJSONHandler(defaultConfig.OutputStream, &opts)
	}

	if defaultConfig.Format == "plain" {
		handler, err = setupPlainFormat(opts)
		if err != nil {
			return nil, err
		}
	}

	if defaultConfig.OtelTracingEnabled {
		handler = otelgoslog.NewTracingHandler(handler)
	}

	return handler, nil
}

func otelMode(ctx context.Context, defaultConfig LoggerGoConfig) (slog.Handler, error) {
	otelGoLogsConfig := otellogs.OtelGoLogsConfig{}

	_, provider, err := otellogs.Init(ctx, otelGoLogsConfig)
	if err != nil {
		return nil, err
	}

	return otelslog.NewHandler(defaultConfig.OtelLoggerName, otelslog.WithLoggerProvider(provider)), nil
}

func setupOtelFormat() (slog.Handler, error) {
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

func setupPlainFormat(opts slog.HandlerOptions) (slog.Handler, error) {
	if defaultConfig.DevMode {

		if defaultConfig.DevFlavor == "slogor" {
			return slogor.NewHandler(defaultConfig.OutputStream, slogor.Options{
				TimeFormat: time.Stamp,
				Level:      opts.Level.Level(),
				ShowSource: opts.AddSource,
			}), nil
		} else if defaultConfig.DevFlavor == "devslog" {
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
