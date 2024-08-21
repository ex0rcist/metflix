package agent_test

import (
	"context"
	"testing"

	"github.com/ex0rcist/metflix/internal/agent"
	"github.com/ex0rcist/metflix/internal/metrics"
	"github.com/stretchr/testify/require"
)

func TestStatsPoll(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	s := agent.NewStats()
	require.Zero(s.Runtime.Alloc)
	require.Zero(s.PollCount)
	require.Zero(s.RandomValue)

	err := s.Poll(ctx)
	require.NoError(err)
	require.Equal(s.PollCount, metrics.Counter(1))
	require.NotZero(s.Runtime.Alloc)
	require.NotZero(s.RandomValue)
	require.NotEmpty(s.System.CPUutilization)

	prev := *s
	err = s.Poll(ctx)
	require.NoError(err)
	require.Equal(s.PollCount, metrics.Counter(2))
	require.NotEqual(prev.RandomValue, s.RandomValue)
	require.NotEqual(prev.Runtime.Alloc, s.Runtime.Alloc)
	require.NotEqual(prev.System.CPUutilization, s.System.CPUutilization)

	err = s.Poll(ctx)
	require.NoError(err)
	require.Equal(s.PollCount, metrics.Counter(3))
}
