package logging

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func LogWarn(messages ...string) {
	logWarn(&log.Logger, messages...)
}

func LogWarnCtx(ctx context.Context, messages ...string) {
	logger := loggerFromContext(ctx)
	logWarn(logger, messages...)
}

func logWarn(logger *zerolog.Logger, messages ...string) {
	msg := optMessagesToString(messages)
	logger.Warn().Msg(msg)
}
