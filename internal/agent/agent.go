package agent

import (
	"time"

	"github.com/ex0rcist/metflix/internal/stats"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Address        string
	PollInterval   time.Duration
	ReportInterval time.Duration
	PollTimeout    time.Duration
	ExportTimeout  time.Duration
}

type Agent struct {
	config Config
	stats  *stats.Stats
	api    *API
}

func New(cfg Config) *Agent {
	stts := stats.NewStats()
	api := NewAPI(cfg.Address, nil)

	agnt := &Agent{config: cfg, stats: stts, api: api}

	return agnt
}

func (app *Agent) Run() {
	go app.startPolling()
	go app.startReporting()
}

func (app *Agent) startPolling() {
	for {
		err := app.stats.Poll()
		if err != nil {
			return // todo: handle errors
		}

		time.Sleep(app.config.PollInterval)
	}
}

func (app *Agent) startReporting() {
	for {
		time.Sleep(app.config.ReportInterval)

		app.reportStats() // todo: handle errors
	}
}

func (app *Agent) reportStats() {
	log.Info().Msg("reporting stats ... ")

	// agent continues polling while repor is in progress, take snapshot
	snapshot := *app.stats

	app.api.
		Report("Alloc", snapshot.Runtime.Alloc).
		Report("BuckHashSys", snapshot.Runtime.BuckHashSys).
		Report("Frees", snapshot.Runtime.Frees).
		Report("GCCPUFraction", snapshot.Runtime.GCCPUFraction).
		Report("GCSys", snapshot.Runtime.GCSys).
		Report("HeapAlloc", snapshot.Runtime.HeapAlloc).
		Report("HeapIdle", snapshot.Runtime.HeapIdle).
		Report("HeapInuse", snapshot.Runtime.HeapInuse).
		Report("HeapObjects", snapshot.Runtime.HeapObjects).
		Report("HeapReleased", snapshot.Runtime.HeapReleased).
		Report("HeapSys", snapshot.Runtime.HeapSys).
		Report("LastGC", snapshot.Runtime.LastGC).
		Report("Lookups", snapshot.Runtime.Lookups).
		Report("MCacheInuse", snapshot.Runtime.MCacheInuse).
		Report("MCacheSys", snapshot.Runtime.MCacheSys).
		Report("MSpanInuse", snapshot.Runtime.MSpanInuse).
		Report("MSpanSys", snapshot.Runtime.MSpanSys).
		Report("Mallocs", snapshot.Runtime.Mallocs).
		Report("NextGC", snapshot.Runtime.NextGC).
		Report("NumForcedGC", snapshot.Runtime.NumForcedGC).
		Report("NumGC", snapshot.Runtime.NumGC).
		Report("OtherSys", snapshot.Runtime.OtherSys).
		Report("PauseTotalNs", snapshot.Runtime.PauseTotalNs).
		Report("StackInuse", snapshot.Runtime.StackInuse).
		Report("StackSys", snapshot.Runtime.StackSys).
		Report("Sys", snapshot.Runtime.Sys).
		Report("TotalAlloc", snapshot.Runtime.TotalAlloc)

	app.api.
		Report("RandomValue", snapshot.RandomValue)

	app.api.
		Report("PollCount", snapshot.PollCount)
}
