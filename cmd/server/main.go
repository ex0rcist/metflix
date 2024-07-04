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
		panic(err)
	}

	log.Info().Msg("server ready")

	err2 := srv.Run()
	if err2 != nil {
		panic(err2)
	}
}
