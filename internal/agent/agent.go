package agent

import (
	"time"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
)

type Agent struct {
	Config *Config
	Stats  *Stats
	API    *API
}

type Config struct {
	Address        entities.Address
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func New() *Agent {
	config := &Config{
		Address:        "0.0.0.0:8080",
		PollInterval:   2 * time.Second,
		ReportInterval: 10 * time.Second,
	}

	stats := NewStats()
	api := NewAPI(&config.Address, nil)

	return &Agent{
		Config: config,
		Stats:  stats,
		API:    api,
	}
}

func (app *Agent) ParseFlags() error {
	address := app.Config.Address

	pflag.VarP(&address, "address", "a", "address:port for HTTP API requests") // HELP: "&"" because Set() has pointer receiver?

	// Task requires us to receive intervals in seconds, not duration, so we have to do it dirty
	pollFlag := pflag.IntP("poll-interval", "p", durationToInt(app.Config.PollInterval), "interval (s) for polling stats")
	reportFlag := pflag.IntP("report-interval", "r", durationToInt(app.Config.ReportInterval), "interval (s) for polling stats")

	pflag.Parse()

	app.Config.PollInterval = intToDuration(*pollFlag)
	app.Config.ReportInterval = intToDuration(*reportFlag)

	// because VarP gets non-pointer value, set it manually
	pflag.Visit(func(f *pflag.Flag) {
		switch f.Name {
		case "address":
			app.Config.Address = address
		}
	})

	return nil
}

func (app *Agent) Run() {
	go app.startPolling()
	go app.startReporting()
}

func (app *Agent) startPolling() {
	for {
		err := app.Stats.Poll()
		if err != nil {
			return // todo: handle errors
		}

		time.Sleep(app.Config.PollInterval)
	}
}

func (app *Agent) startReporting() {
	for {
		time.Sleep(app.Config.ReportInterval)

		app.reportStats() // todo: handle errors
	}
}

func (app *Agent) reportStats() {
	log.Info().Msg("reporting stats ... ")

	// agent continues polling while report is in progress, take snapshot?
	snapshot := *app.Stats

	app.API.
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

	app.API.
		Report("RandomValue", snapshot.RandomValue)

	app.API.
		Report("PollCount", snapshot.PollCount)
}

func intToDuration(s int) time.Duration {
	return time.Duration(s) * time.Second
}

func durationToInt(d time.Duration) int {
	return int(d.Seconds())
}
