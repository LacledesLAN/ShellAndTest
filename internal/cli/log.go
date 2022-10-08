package cli

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitZeroLog(logLevel string) {
	// set message field name
	zerolog.MessageFieldName = "msg"

	// init zerolog
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	switch logLevel {
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}
}
