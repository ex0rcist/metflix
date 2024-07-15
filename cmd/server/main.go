package main

import (
	"context"

	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/server"
)

func main() {
	ctx := logging.Setup(context.Background())

	logging.LogInfo(ctx, "starting server...")

	srv, err := server.New()
	if err != nil {
		logging.LogFatal(ctx, err)
	}

	err = srv.ParseFlags()
	if err != nil {
		logging.LogFatal(ctx, err)
	}

	logging.LogInfo(ctx, srv.Config.String())
	logging.LogInfo(ctx, "server ready") // TODO: must be after run?

	err = srv.Run()
	if err != nil {
		logging.LogFatal(ctx, err)
	}
}
