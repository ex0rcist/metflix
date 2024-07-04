package agent

import (
	"time"

	"github.com/ex0rcist/metflix/internal/stats"
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
	app.api.
		Report("Alloc", app.stats.Runtime.Alloc).
		Report("BuckHashSys", app.stats.Runtime.BuckHashSys).
		Report("Frees", app.stats.Runtime.Frees).
		Report("GCCPUFraction", app.stats.Runtime.GCCPUFraction).
		Report("GCSys", app.stats.Runtime.GCSys).
		Report("HeapAlloc", app.stats.Runtime.HeapAlloc).
		Report("HeapIdle", app.stats.Runtime.HeapIdle).
		Report("HeapInuse", app.stats.Runtime.HeapInuse).
		Report("HeapObjects", app.stats.Runtime.HeapObjects).
		Report("HeapReleased", app.stats.Runtime.HeapReleased).
		Report("HeapSys", app.stats.Runtime.HeapSys).
		Report("LastGC", app.stats.Runtime.LastGC).
		Report("Lookups", app.stats.Runtime.Lookups).
		Report("MCacheInuse", app.stats.Runtime.MCacheInuse).
		Report("MCacheSys", app.stats.Runtime.MCacheSys).
		Report("MSpanInuse", app.stats.Runtime.MSpanInuse).
		Report("MSpanSys", app.stats.Runtime.MSpanSys).
		Report("Mallocs", app.stats.Runtime.Mallocs).
		Report("NextGC", app.stats.Runtime.NextGC).
		Report("NumForcedGC", app.stats.Runtime.NumForcedGC).
		Report("NumGC", app.stats.Runtime.NumGC).
		Report("OtherSys", app.stats.Runtime.OtherSys).
		Report("PauseTotalNs", app.stats.Runtime.PauseTotalNs).
		Report("StackInuse", app.stats.Runtime.StackInuse).
		Report("StackSys", app.stats.Runtime.StackSys).
		Report("Sys", app.stats.Runtime.Sys).
		Report("TotalAlloc", app.stats.Runtime.TotalAlloc)

	app.api.
		Report("RandomValue", app.stats.RandomValue)

	app.api.
		Report("PollCount", app.stats.PollCount)
}
