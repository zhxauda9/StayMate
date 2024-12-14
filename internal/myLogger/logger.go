package myLogger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

/*
	Logger Levels:

	Trace(), Debug(), Info(), Warn(), and Error() are different logging levels provided by zerolog.

	Adding Context:
	Use methods like Str(), Int(), Bool() to add structured data to the logs.

	Final Message:
	The Msg("message text") method outputs the log message.
*/

func NewZeroLogger() *zerolog.Logger {
	logger := zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).Level(zerolog.TraceLevel).With().Timestamp().Caller().Logger()

	return &logger
}
