package logging

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Log with level=fatal and formatting
func LogFatalF(format string, err error) {
	fErr := fmt.Errorf(format, err)
	logFatal(&log.Logger, fErr)
}

// Log with level=fatal
func LogFatal(err error, messages ...string) {
	logFatal(&log.Logger, err, messages...)
}

// Log with context (request_id) and level=fatal
func LogFatalCtx(ctx context.Context, err error, messages ...string) {
	logger := loggerFromContext(ctx)
	logFatal(logger, err, messages...)
}

func logFatal(logger *zerolog.Logger, err error, messages ...string) {
	msg := optMessagesToString(messages)

	if isDebugLevel() {
		logger.Fatal().Stack().Err(err).Msg(msg) // Stack() must be called before Err()
	} else {
		logger.Fatal().Err(err).Msg(msg)
	}
}
