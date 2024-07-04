package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Setup() {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339, // time.RFC822
	}

	l := zerolog.New(output).With().Timestamp()
	log.Logger = l.Logger()
}
