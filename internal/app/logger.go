package app

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func createLogger() zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	logger := zerolog.New(output).With().Timestamp().Logger()
	logger.Level(zerolog.DebugLevel)

	return logger
}
