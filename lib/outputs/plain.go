package outputs

import (
	"log/slog"
	"os"
	"time"

	"github.com/golang-cz/devslog"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"github.com/wasilak/loggergo/lib/types"
	"gitlab.com/greyxor/slogor"
)

// setupPlainFormat sets up a slog.Handler for plain format.
// If defaultConfig.DevMode is true, it checks the defaultConfig.DevFlavor and sets up the appropriate handler based on the flavor.
// Returns the handler and any error encountered.
func SetupPlainFormat(opts slog.HandlerOptions, defaultConfig types.Config) (slog.Handler, error) {
	if defaultConfig.DevMode {

		if defaultConfig.DevFlavor == types.DevFlavorSlogor {
			return slogor.NewHandler(defaultConfig.OutputStream, slogor.ShowSource(), slogor.SetTimeFormat(time.Stamp), slogor.SetLevel(opts.Level.Level())), nil
		} else if defaultConfig.DevFlavor == types.DevFlavorDevslog {
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
