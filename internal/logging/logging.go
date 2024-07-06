package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func Setup() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339, // time.RFC822
	}

	// HELP: didn't manage to put error stack :(
	l := zerolog.New(output).With().Timestamp()

	log.Logger = l.Logger()
}
