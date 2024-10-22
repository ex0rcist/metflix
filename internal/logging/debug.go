package logging

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Log with level=debug and formatting
func LogDebugF(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	logDebug(&log.Logger, message)
}

// Log with level=debug
func LogDebug(messages ...string) {
	logDebug(&log.Logger, messages...)
}

// Log with context (request_id) and level=debug
func LogDebugCtx(ctx context.Context, messages ...string) {
	logger := loggerFromContext(ctx)
	logDebug(logger, messages...)
}

func logDebug(logger *zerolog.Logger, messages ...string) {
	msg := optMessagesToString(messages)
	logger.Debug().Msg(msg)
}
