package logging

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

type config struct {
	Env string `env:"APP_ENV" envDefault:"development"`
}

func Setup(ctx context.Context) context.Context {
	cfg := parseConfig()

	var output io.Writer
	switch cfg.Env {
	case "development":
		output = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339Nano}
	case "production":
		output = os.Stdout
	}

	logger := zerolog.New(output).With().Timestamp().Logger()

	// set global logger
	log.Logger = logger

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.DefaultContextLogger = &logger

	// added logger to ctx
	return logger.WithContext(ctx)
}

func NewError(err error) error { // wrap err to pkg/errors
	return errors.New(err.Error()) // TODO: can we remove this func from stack?
}

func LogError(ctx context.Context, err error, messages ...string) {
	msg := optMessagesToString(messages)
	zerolog.Ctx(ctx).Error().Stack().Err(err).Msg(msg)
}

func LogFatal(ctx context.Context, err error, messages ...string) {
	msg := optMessagesToString(messages)
	zerolog.Ctx(ctx).Fatal().Stack().Err(err).Msg(msg)
}

func LogInfo(ctx context.Context, messages ...string) {
	msg := optMessagesToString(messages)
	zerolog.Ctx(ctx).Info().Msg(msg)
}

func optMessagesToString(messages []string) string {
	if len(messages) == 0 {
		return ""
	}

	return strings.Join(messages, "; ")
}

func parseConfig() config {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	return cfg
}

func (c config) isDevelopment() bool {
	return c.Env == "development"
}

func (c config) isProduction() bool {
	return c.Env == "production"
}
