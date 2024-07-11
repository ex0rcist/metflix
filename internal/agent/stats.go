package agent

import (
	"math/rand"
	"time"

	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/metrics"
)

type Stats struct {
	Runtime     RuntimeStats
	PollCount   metrics.Counter
	RandomValue metrics.Gauge

	generator *rand.Rand
}

func NewStats() *Stats {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &Stats{generator: r}
}

func (m *Stats) Poll() error {
	logging.LogInfo("polling stats ... ")

	m.PollCount++
	m.RandomValue = metrics.Gauge(m.generator.Float64())
	m.Runtime.Poll()

	return nil
}
