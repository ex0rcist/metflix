package agent

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	agnt := New()
	require.NotPanics(t, agnt.Run)
}

func TestReportStats(t *testing.T) {
	// HELP: как тестировать приватные?
}
