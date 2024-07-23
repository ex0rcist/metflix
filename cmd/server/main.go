package main

import (
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/server"
)

func main() {
	logging.Setup()

	logging.LogInfo("starting server...")

	srv, err := server.New()
	if err != nil {
		logging.LogFatal(err)
	}

	err = srv.ParseFlags()
	if err != nil {
		logging.LogFatal(err)
	}

	err = srv.Run()
	if err != nil {
		logging.LogFatal(err)
	}
}
