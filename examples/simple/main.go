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

	ctx = context.WithValue(ctx, "test", "aaaaaaa")

	loggerConfig := loggergo.Config{
		Level:        loggergo.LogLevelFromString(*logLevel),
		Format:       loggergo.LogFormatFromString(*logFormat),
		OutputStream: os.Stdout,
		DevMode:      *devMode,
		Output:       loggergo.OutputConsole,
		ContextKeys:  []string{"test", "test2", "test3"},
		// ContextKeysDefault: "default",
	}

	ctx = context.WithValue(ctx, "test3", "bbbbbb")

	if *otelEnabled {

		loggerConfig.OtelServiceName = os.Getenv("OTEL_SERVICE_NAME")
		loggerConfig.Output = loggergo.OutputFanout
		loggerConfig.OtelLoggerName = "github.com/wasilak/loggergo" // your package name goes here
		loggerConfig.OtelTracingEnabled = false
	}

	// this registers loggergo as slog.Default()
	ctx, _, err := loggergo.LoggerInit(ctx, loggerConfig)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		os.Exit(1)
	}

	logLevelConfig := loggergo.GetLogLevelAccessor()

	// every log below should be have fields test & test3 with values from abov, but not test2, as it is not in the context
	// in case ContextKeysDefault: "default", test2 would have value "default"

	slog.InfoContext(ctx, "Hello, World!")
	slog.DebugContext(ctx, "Hello, Debug #1!") // will not be printed as the default log level is info

	logLevelConfig.Set(slog.LevelDebug)
	slog.DebugContext(ctx, "Hello, Debug #2!") // now it will be printed as the default log level is debug
}
