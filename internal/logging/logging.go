package logging

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func Setup(ctx context.Context) context.Context {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339Nano}
	// output := os.Stdout
	logger := zerolog.New(output).With().Timestamp().Logger()

	// set global logger
	log.Logger = logger

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.DefaultContextLogger = &logger // HELP: is it bad?   ->   logging.LogInfo(context.Background(), "some message")

	// added logger to ctx
	return logger.WithContext(ctx)
}

func NewError(err error) error { // wrap err to pkg/errors
	return errors.New(err.Error()) // TODO: can we remove this func from stack?
}

func LogError(ctx context.Context, err error, messages ...string) {
	msg := optMessagesToString(messages)
	log.Ctx(ctx).Error().Stack().Err(err).Msg(msg)
}

func LogFatal(ctx context.Context, err error, messages ...string) {
	msg := optMessagesToString(messages)
	log.Ctx(ctx).Fatal().Stack().Err(err).Msg(msg)
}

func LogInfo(ctx context.Context, messages ...string) {
	msg := optMessagesToString(messages)
	log.Ctx(ctx).Info().Msg(msg)
}

func optMessagesToString(messages []string) string {
	if len(messages) == 0 {
		return ""
	}

	return strings.Join(messages, "; ")
}
