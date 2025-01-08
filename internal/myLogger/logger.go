package myLogger

import (
	"io"
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

var Log *zerolog.Logger

func NewZeroLogger() *zerolog.Logger {
	logger := zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).Level(zerolog.TraceLevel).With().Timestamp().Caller().Logger()
	return &logger
}

func NewZeroLoggerV2() *zerolog.Logger {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// Fall back to only logging to stderr if the file cannot be opened
		logWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
		logger := zerolog.New(logWriter).
			Level(zerolog.TraceLevel).
			With().
			Timestamp().
			Caller().
			Logger()

		logger.Warn().Msg("Logger initialized but only console logging availible")
		return &logger
	}

	multiWriter := io.MultiWriter(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}, // Console output
		file, // File output
	)

	logger := zerolog.New(multiWriter).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	return &logger
}
