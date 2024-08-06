package logging

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func LogInfoF(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	logInfo(&log.Logger, message)
}

func LogInfo(messages ...string) {
	logInfo(&log.Logger, messages...)
}

func LogInfoCtx(ctx context.Context, messages ...string) {
	logger := loggerFromContext(ctx)
	logInfo(logger, messages...)
}

func logInfo(logger *zerolog.Logger, messages ...string) {
	msg := optMessagesToString(messages)
	logger.Info().Msg(msg)
}
