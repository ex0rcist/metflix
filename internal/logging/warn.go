package logging

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Log with level=warn and formatting
func LogWarnF(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	logWarn(&log.Logger, message)
}

// Log with level=warn
func LogWarn(messages ...string) {
	logWarn(&log.Logger, messages...)
}

// Log with context (request_id) and level=warn
func LogWarnCtx(ctx context.Context, messages ...string) {
	logger := loggerFromContext(ctx)
	logWarn(logger, messages...)
}

func logWarn(logger *zerolog.Logger, messages ...string) {
	msg := optMessagesToString(messages)
	logger.Warn().Msg(msg)
}
