package logging

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Log with level=error and formatting
func LogErrorF(format string, err error) {
	fErr := fmt.Errorf(format, err)
	logError(&log.Logger, fErr)
}

// Log with level=error
func LogError(err error, messages ...string) {
	logError(&log.Logger, err, messages...)
}

// Log with context (request_id) and level=error
func LogErrorCtx(ctx context.Context, err error, messages ...string) {
	logger := loggerFromContext(ctx)
	logError(logger, err, messages...)
}

func logError(logger *zerolog.Logger, err error, messages ...string) {
	msg := optMessagesToString(messages)

	if isDebugLevel() {
		logger.Error().Stack().Err(err).Msg(msg) // Stack() must be called before Err()
	} else {
		logger.Error().Err(err).Msg(msg)
	}
}
