package logging

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func LogDebugF(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	logDebug(&log.Logger, message)
}

func LogDebug(messages ...string) {
	logDebug(&log.Logger, messages...)
}

func LogDebugCtx(ctx context.Context, messages ...string) {
	logger := loggerFromContext(ctx)
	logDebug(logger, messages...)
}

func logDebug(logger *zerolog.Logger, messages ...string) {
	msg := optMessagesToString(messages)
	logger.Debug().Msg(msg)
}
