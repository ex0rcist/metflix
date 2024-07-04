package main

import (
	"time"

	"github.com/ex0rcist/metflix/internal/agent"
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/rs/zerolog/log"
)

func main() {
	logging.Setup()

	conf := agent.Config{
		Address:        "http://0.0.0.0:8080",
		ReportInterval: 10 * time.Second,
		PollInterval:   2 * time.Second,
		PollTimeout:    2 * time.Second,
		ExportTimeout:  4 * time.Second,
	} // todo: yml?

	log.Info().Msg("starting agent...")

	agnt := agent.New(conf)
	agnt.Run()

	log.Info().Msg("agent ready")

	for {
		// fixme: tmp hack for goroutine
		time.Sleep(time.Second * 1)
	}

}
