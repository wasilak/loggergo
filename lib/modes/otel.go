package modes

import (
	"context"
	"log/slog"

	"github.com/wasilak/loggergo/lib/types"
	otellogs "github.com/wasilak/otelgo/logs"
	"go.opentelemetry.io/contrib/bridges/otelslog"
)

// otelMode returns a slog.Handler for OpenTelemetry mode based on the provided defaultConfig.
// It initializes the otellogs package and returns a handler with the otelslog.WithLoggerProvider option.
// Returns the handler and any error encountered.
func OtelMode(ctx context.Context, defaultConfig types.Config) (slog.Handler, context.Context, error) {
	otelGoLogsConfig := otellogs.OtelGoLogsConfig{}

	ctx, provider, err := otellogs.Init(ctx, otelGoLogsConfig)
	if err != nil {
		return nil, ctx, err
	}

	return otelslog.NewHandler(defaultConfig.OtelLoggerName, otelslog.WithLoggerProvider(provider)), ctx, nil
}
