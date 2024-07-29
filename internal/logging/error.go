package logging

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func LogError(err error, messages ...string) {
	logError(&log.Logger, err, messages...)
}

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
