package common

import (
	"os"

	"github.com/mattn/go-isatty"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapgrpc"
)

const Version = "0.0.1"

func NewCommonModule() fx.Option {
	var logger *zap.Logger
	if isatty.IsTerminal(os.Stderr.Fd()) {
		// New*() are shortcuts that provide pre-defined sets of config.
		logger, _ = zap.NewDevelopment() // text-formatting, debug level
	} else {
		logger, _ = zap.NewProduction() // JSON-formatting, info level
	}

	return fx.Options(
		fx.Provide(func() *zap.SugaredLogger {
			logger.Debug("NewLogger")
			return logger.Sugar()

		}),
		fx.Logger(zapgrpc.NewLogger(logger)),
	)

}
