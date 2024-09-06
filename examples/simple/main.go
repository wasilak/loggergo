package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/wasilak/loggergo"
)

func main() {
	ctx := context.Background()
	logLevel := flag.String("log-level", os.Getenv("LOG_LEVEL"), "log level (debug, info, warn, error, fatal)")
	logFormat := flag.String("log-format", os.Getenv("LOG_FORMAT"), "log format (json, plain, otel)")
	devMode := flag.Bool("dev-mode", false, "Development mode")
	otelEnabled := flag.Bool("otel-enabled", false, "OpenTelemetry traces enabled")
	flag.Parse()

	loggerConfig := loggergo.Config{
		Level:        loggergo.LogLevelFromString(*logLevel),
		Format:       loggergo.LogFormatFromString(*logFormat),
		OutputStream: os.Stdout,
		DevMode:      *devMode,
		Output:       loggergo.OutputConsole,
	}

	if *otelEnabled {

		loggerConfig.OtelServiceName = os.Getenv("OTEL_SERVICE_NAME")
		loggerConfig.Output = loggergo.OutputFanout
		loggerConfig.OtelLoggerName = "github.com/wasilak/loggergo" // your package name goes here
		loggerConfig.OtelTracingEnabled = false
	}

	// this registers loggergo as slog.Default()
	_, err := loggergo.LoggerInit(ctx, loggerConfig)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		os.Exit(1)
	}

	slog.InfoContext(ctx, "Hello, World!")
}
