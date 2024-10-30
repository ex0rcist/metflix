package main

import (
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/server"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	logging.Setup()
	logging.LogInfoF("starting server v%s [%s, #%s]...", buildVersion, buildDate, buildCommit)

	srv, err := server.New()
	if err != nil {
		logging.LogFatal(err)
	}

	srv.Start()
}
