package Mlogger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func NewZeroLogger() *zerolog.Logger {
	logger := zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).Level(zerolog.TraceLevel).With().Timestamp().Caller().Logger()
	return &logger
}
