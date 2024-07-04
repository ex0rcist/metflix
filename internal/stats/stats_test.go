package stats_test

import (
	"testing"

	"github.com/ex0rcist/metflix/internal/metrics"
	"github.com/ex0rcist/metflix/internal/stats"
	"github.com/stretchr/testify/require"
)

func TestStatsPoll(t *testing.T) {
	require := require.New(t)

	s := stats.NewStats()
	require.Zero(s.Runtime.Alloc)
	require.Zero(s.PollCount)
	require.Zero(s.RandomValue)

	err := s.Poll()
	require.NoError(err)
	require.Equal(s.PollCount, metrics.Counter(1))
	require.NotZero(s.Runtime.Alloc)
	require.NotZero(s.RandomValue)

	prev := *s
	err = s.Poll()
	require.NoError(err)
	require.Equal(s.PollCount, metrics.Counter(2))
	require.NotEqual(prev.RandomValue, s.RandomValue)
	require.NotEqual(prev.Runtime.Alloc, s.Runtime.Alloc)
}
