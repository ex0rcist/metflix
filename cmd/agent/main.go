package main

import (
	"github.com/ex0rcist/metflix/internal/agent"
	"github.com/ex0rcist/metflix/internal/logging"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	logging.Setup()

	logging.LogInfoF("starting agent v%s [%s, #%s]...", buildVersion, buildDate, buildCommit)

	agnt, err := agent.New()
	if err != nil {
		logging.LogFatal(err)
	}

	agnt.Run()
}
