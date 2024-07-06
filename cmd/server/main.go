package main

import (
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/server"

	"github.com/rs/zerolog/log"
)

func main() {
	logging.Setup()

	log.Info().Msg("starting server...")

	srv, err := server.New()
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	err = srv.ParseFlags()
	if err != nil {

		log.Error().Err(err).Msg("")
		return
	}

	log.Info().Msgf(srv.Config.String())
	log.Info().Msg("server ready") // TODO: must be after run?

	err = srv.Run()
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}
}
