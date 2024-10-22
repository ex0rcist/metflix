package logging

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Log with level=info and formatting
func LogInfoF(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	logInfo(&log.Logger, message)
}

// Log with level=info
func LogInfo(messages ...string) {
	logInfo(&log.Logger, messages...)
}

// Log with context (request_id) and level=info
func LogInfoCtx(ctx context.Context, messages ...string) {
	logger := loggerFromContext(ctx)
	logInfo(logger, messages...)
}

func logInfo(logger *zerolog.Logger, messages ...string) {
	msg := optMessagesToString(messages)
	logger.Info().Msg(msg)
}
