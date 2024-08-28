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

	agnt.Run()
}
