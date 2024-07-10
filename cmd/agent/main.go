package main

import (
	"github.com/ex0rcist/metflix/internal/agent"
	"github.com/ex0rcist/metflix/internal/logging"
)

func main() {
	logging.Setup()

	logging.LogInfo("starting agent...")

	agnt, err := agent.New()
	if err != nil {
		logging.LogFatal(err)
	}

	err = agnt.ParseFlags()
	if err != nil {
		logging.LogFatal(err)
	}

	logging.LogInfo(agnt.Config.String())

	agnt.Run()

	logging.LogInfo("agent ready")
}
