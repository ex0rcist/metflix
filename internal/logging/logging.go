package logging

import (
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func Setup() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	l := zerolog.New(output).With().Timestamp()
	log.Logger = l.Logger()
}

func NewError(err error) error { // wrap err to pkg/errors
	return errors.New(err.Error()) // TODO: can we remove this func from stack?
}

func LogError(err error, messages ...string) {
	msg := optMessagesToString(messages)
	log.Error().Stack().Err(err).Msg(msg)
}

func LogFatal(err error, messages ...string) {
	msg := optMessagesToString(messages)
	log.Fatal().Stack().Err(err).Msg(msg)
}

func LogInfo(messages ...string) {
	msg := optMessagesToString(messages)
	log.Info().Msg(msg)
}

func optMessagesToString(messages []string) string {
	if len(messages) == 0 {
		return ""
	}

	return strings.Join(messages, "; ")
}
