package main

import (
	"context"

	"github.com/ex0rcist/metflix/internal/agent"
	"github.com/ex0rcist/metflix/internal/logging"
)

func main() {
	ctx := logging.Setup(context.Background())

	logging.LogInfo(ctx, "starting agent...")

	agnt, err := agent.New()
	if err != nil {
		logging.LogFatal(ctx, err)
	}

	err = agnt.ParseFlags()
	if err != nil {
		logging.LogFatal(ctx, err)
	}

	logging.LogInfo(ctx, agnt.Config.String())

	agnt.Run()

	logging.LogInfo(ctx, "agent ready")
}
