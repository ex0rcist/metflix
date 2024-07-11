package agent_test

import (
	"testing"

	"github.com/ex0rcist/metflix/internal/agent"
	"github.com/stretchr/testify/require"
)

func TestRuntimePoll(t *testing.T) {
	require := require.New(t)

	rs := agent.RuntimeStats{}
	require.Zero(rs.TotalAlloc)

	rs.Poll()
	require.NotZero(rs.TotalAlloc)

	newRs := agent.RuntimeStats{}
	newRs.Poll()

	require.NotEqual(newRs.TotalAlloc, rs.TotalAlloc)
}
