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

	agnt := agent.New()

	err := agnt.ParseFlags()
	if err != nil {
		panic(err)
	}

	log.Info().Msgf("agent args: a: %v, p: %v, r: %v", agnt.Config.Address, agnt.Config.PollInterval, agnt.Config.ReportInterval)

	agnt.Run()

	log.Info().Msgf("agent ready")

	for {
		// fixme: tmp hack for goroutine
		time.Sleep(time.Second * 1)
	}
}
