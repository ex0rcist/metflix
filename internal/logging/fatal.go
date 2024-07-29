package logging

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func LogFatal(err error, messages ...string) {
	logFatal(&log.Logger, err, messages...)
}

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
