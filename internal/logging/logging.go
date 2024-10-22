package logging

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

type config struct {
	ENV string `env:"APP_ENV" envDefault:"development"`
}

// Initialize logging
func Setup() {
	cfg := parseConfig()

	var output io.Writer
	var logger zerolog.Logger

	switch cfg.ENV {
	case "tracing":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

		output = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339Nano}
	case "development":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)

		output = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339Nano}
	case "production":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)

		output = os.Stdout
	}

	loggerCtx := zerolog.New(output).With().Timestamp()
	switch {
	case isTraceLevel():
		logger = loggerCtx.Caller().Logger()
	default:
		logger = loggerCtx.Logger()
	}

	log.Logger = logger
	zerolog.DefaultContextLogger = &logger
}

func parseConfig() config {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	return cfg
}

func optMessagesToString(messages []string) string {
	if len(messages) == 0 {
		return ""
	}

	// remove empty
	var result []string
	for _, str := range messages {
		if str != "" {
			result = append(result, str)
		}
	}

	return strings.Join(result, "; ")
}

func loggerFromContext(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}

func isDebugLevel() bool {
	return zerolog.GlobalLevel() == zerolog.DebugLevel
}

func isTraceLevel() bool {
	return zerolog.GlobalLevel() == zerolog.TraceLevel
}
