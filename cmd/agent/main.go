package main

import (
	"time"

	"github.com/ex0rcist/metflix/internal/agent"
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/rs/zerolog/log"
)

func main() {
	logging.Setup()

	log.Info().Msg("starting agent...")

	agnt, err := agent.New()
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	err = agnt.ParseFlags()
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	log.Info().Msg(agnt.Config.String())

	err = agnt.Run()
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	log.Info().Msg("agent ready")

	for { // fixme: tmp hack for goroutine
		time.Sleep(time.Second * 1)
	}
}
