package loggergo

import (
	"os"
	"strings"

	"log/slog"

	otelgoslog "github.com/wasilak/otelgo/slog"

	"github.com/golang-cz/devslog"

	"dario.cat/mergo"
)

type LoggerGoConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
	Dev    bool   `json:"dev"`
}

var defaultConfig = LoggerGoConfig{
	Level:  "info",
	Format: "plain",
	Dev:    false,
}

func LoggerInit(config LoggerGoConfig, additionalAttrs ...any) *slog.Logger {

	mergo.Merge(&defaultConfig, config)

	var logLevel slog.Leveler

	switch strings.ToLower(config.Level) {
	case "info":
		logLevel = slog.LevelInfo
	case "error":
		logLevel = slog.LevelError
	case "warn":
		logLevel = slog.LevelWarn
	case "debug":
		logLevel = slog.LevelDebug
	default:
		logLevel = slog.LevelInfo
	}

	opts := slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	}

	if strings.ToLower(config.Format) == "json" {
		slog.SetDefault(slog.New(otelgoslog.NewTracingHandler(slog.NewJSONHandler(os.Stderr, &opts))))
	} else {
		if config.Dev {
			devOpts := &devslog.Options{
				HandlerOptions:    &opts,
				MaxSlicePrintSize: 10,
				SortKeys:          true,
			}
			slog.SetDefault(slog.New(otelgoslog.NewTracingHandler(devslog.NewHandler(os.Stderr, devOpts))))
		} else {
			slog.SetDefault(slog.New(otelgoslog.NewTracingHandler(slog.NewTextHandler(os.Stderr, &opts))))
		}
	}

	for _, v := range additionalAttrs {
		slog.SetDefault(slog.Default().With(v))
	}

	return slog.Default()
}
