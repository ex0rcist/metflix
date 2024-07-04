package agent_test

import (
	"testing"
	"time"

	"github.com/ex0rcist/metflix/internal/agent"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	cfg := agent.Config{
		Address:        "http://0.0.0.0:8080",
		ReportInterval: 10 * time.Second,
		PollInterval:   2 * time.Second,
		PollTimeout:    2 * time.Second,
		ExportTimeout:  4 * time.Second,
	}

	agnt := agent.New(cfg)
	require.NotPanics(t, agnt.Run)
}

func TestReportStats(t *testing.T) {
	// todo: как тестировать приватные?
}
