package stats

import (
	"math/rand"
	"time"

	"github.com/ex0rcist/metflix/internal/metrics"
	"github.com/rs/zerolog/log"
)

type Stats struct {
	Runtime     RuntimeStats
	PollCount   metrics.Gauge
	RandomValue metrics.Gauge

	generator *rand.Rand
}

func NewStats() *Stats {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &Stats{generator: r}
}

func (m *Stats) Poll() error {
	log.Info().Msg("polling stats ... ")

	m.PollCount++
	m.RandomValue = metrics.Gauge(m.generator.Float64())
	m.Runtime.Poll()

	// todo: errors?

	return nil
}
