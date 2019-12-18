package common

import (
	"github.com/mattn/go-isatty"
	"go.uber.org/zap"
	"os"
)

const Version = "0.0.1"

func NewLogger() *zap.SugaredLogger {
	var logger *zap.Logger
	if isatty.IsTerminal(os.Stderr.Fd()) {
		// New*() are shortcuts that provide pre-defined sets of config.
		logger, _ = zap.NewDevelopment() // text-formatting, debug level
	} else {
		logger, _ = zap.NewProduction() // JSON-formatting, info level
	}
	defer logger.Sync() // flushes buffer, if any

	return logger.Sugar()
}