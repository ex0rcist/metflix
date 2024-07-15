package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/spf13/pflag"
)

type Agent struct {
	Config *Config
	Stats  *Stats
	API    *API

	wg sync.WaitGroup
}

type Config struct {
	Address        entities.Address `env:"ADDRESS"`
	PollInterval   int              `env:"POLL_INTERVAL"`
	ReportInterval int              `env:"REPORT_INTERVAL"`
}

func New() (*Agent, error) {
	config := &Config{
		Address:        "0.0.0.0:8080",
		PollInterval:   2,
		ReportInterval: 10,
	}

	stats := NewStats()
	api := NewAPI(&config.Address, nil)

	return &Agent{
		Config: config,
		Stats:  stats,
		API:    api,
	}, nil
}

func (a *Agent) ParseFlags() error {
	address := a.Config.Address

	pflag.VarP(&address, "address", "a", "address:port for HTTP API requests") // HELP: "&"" because Set() has pointer receiver?

	pflag.IntVarP(&a.Config.PollInterval, "poll-interval", "p", a.Config.PollInterval, "interval (s) for polling stats")
	pflag.IntVarP(&a.Config.ReportInterval, "report-interval", "r", a.Config.ReportInterval, "interval (s) for polling stats")

	pflag.Parse()

	// because VarP gets non-pointer value, set it manually
	pflag.Visit(func(f *pflag.Flag) {
		switch f.Name {
		case "address":
			a.Config.Address = address
		}
	})

	if err := env.Parse(a.Config); err != nil {
		return err
	}

	return nil
}

func (a *Agent) Run() {
	a.wg.Add(2)

	go a.startPolling()
	go a.startReporting()

	a.wg.Wait()
}

func (a *Agent) startPolling() {
	defer a.wg.Done()

	ctx := context.Background()

	for {
		err := a.Stats.Poll()
		if err != nil {
			logging.LogError(ctx, err)
		}

		time.Sleep(intToDuration(a.Config.PollInterval))
	}
}

func (a *Agent) startReporting() {
	defer a.wg.Done()

	for {
		time.Sleep(intToDuration(a.Config.ReportInterval))

		a.reportStats()
	}
}

func (a *Agent) reportStats() {
	ctx := context.Background()

	logging.LogInfo(ctx, "reporting stats ... ")

	// agent continues polling while report is in progress, take snapshot?
	snapshot := *a.Stats

	a.API.
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

	a.API.
		Report("RandomValue", snapshot.RandomValue)

	a.API.
		Report("PollCount", snapshot.PollCount)

	// because metrics.Counter adds value to itself
	a.Stats.PollCount -= snapshot.PollCount
}

func intToDuration(s int) time.Duration {
	return time.Duration(s) * time.Second
}

func (c Config) String() string {
	out := "agent config: "

	out += fmt.Sprintf("address=%v \t", c.Address)
	out += fmt.Sprintf("poll-interval=%v \t", c.PollInterval)
	out += fmt.Sprintf("report-interval=%v \t", c.ReportInterval)
	return out
}
