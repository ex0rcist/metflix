package stats_test

import (
	"testing"

	"github.com/ex0rcist/metflix/internal/stats"
	"github.com/stretchr/testify/require"
)

func TestRuntimePoll(t *testing.T) {
	require := require.New(t)

	rs := stats.RuntimeStats{}
	require.Zero(rs.TotalAlloc)

	rs.Poll()
	require.NotZero(rs.TotalAlloc)

	newRs := stats.RuntimeStats{}
	newRs.Poll()

	require.NotEqual(newRs.TotalAlloc, rs.TotalAlloc)
}
