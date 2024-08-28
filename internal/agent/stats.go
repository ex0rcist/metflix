package agent

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/metrics"
	"golang.org/x/sync/errgroup"
)

type Stats struct {
	System      SystemStats
	Runtime     RuntimeStats
	PollCount   metrics.Counter
	RandomValue metrics.Gauge

	generator *rand.Rand
}

func NewStats() *Stats {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &Stats{generator: r}
}

func (m *Stats) Poll(ctx context.Context) error {
	logging.LogDebug("polling stats ... ")

	m.PollCount++
	m.RandomValue = metrics.Gauge(m.generator.Float64())

	g, _ := errgroup.WithContext(ctx)

	g.Go(func() error {
		m.Runtime.Poll()
		return nil
	})

	g.Go(func() error {
		return m.System.Poll(ctx)
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("failed to gather metrics: %w", err)
	}

	return nil
}
