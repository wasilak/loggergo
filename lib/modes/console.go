package modes

import (
	"log/slog"

	"github.com/wasilak/loggergo/lib"
	"github.com/wasilak/loggergo/lib/outputs"
	"github.com/wasilak/loggergo/lib/types"
	otelgoslog "github.com/wasilak/otelgo/slog"
)

// consoleMode returns a slog.Handler based on the provided defaultConfig and opts.
// It checks the defaultConfig.Format and sets up the appropriate handler based on the format.
// If defaultConfig.OtelTracingEnabled is true, it wraps the handler with otelgoslog.NewTracingHandler.
// Returns the handler and any error encountered.
func ConsoleMode(opts slog.HandlerOptions) (slog.Handler, error) {
	var handler slog.Handler
	var err error

	if lib.GetConfig().Format == types.LogFormatOtel {
		return outputs.SetupOtelFormat()
	}

	if lib.GetConfig().Format == types.LogFormatJSON {
		handler = slog.NewJSONHandler(lib.GetConfig().OutputStream, &opts)
	}

	if lib.GetConfig().Format == types.LogFormatText {
		handler, err = outputs.SetupPlainFormat(opts)
		if err != nil {
			return nil, err
		}
	}

	if lib.GetConfig().OtelTracingEnabled {
		handler = otelgoslog.NewTracingHandler(handler)
	}

	return handler, nil
}
