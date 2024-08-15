package agent

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/services"
	"github.com/ex0rcist/metflix/internal/utils"
	"github.com/spf13/pflag"
)

type Agent struct {
	Config   *Config
	Stats    *Stats
	Exporter Exporter

	wg sync.WaitGroup
}

type Config struct {
	Address        entities.Address `env:"ADDRESS"`
	PollInterval   int              `env:"POLL_INTERVAL"`
	ReportInterval int              `env:"REPORT_INTERVAL"`
	Secret         entities.Secret  `env:"KEY"`
}

func New() (*Agent, error) {
	config := &Config{
		Address:        "0.0.0.0:8080",
		PollInterval:   2,
		ReportInterval: 10,
	}

	err := parseConfig(config)
	if err != nil {
		return nil, err
	}

	var signer services.Signer
	if len(config.Secret) > 0 {
		signer = services.NewSignerService(config.Secret)
	}

	exporter := NewMetricsExporter(&config.Address, http.DefaultTransport, signer)

	return &Agent{
		Config:   config,
		Stats:    NewStats(),
		Exporter: exporter,
	}, nil
}

func (a *Agent) Run() {
	logging.LogInfo(a.Config.String())
	logging.LogInfo("agent ready")

	a.wg.Add(2)

	go a.startPolling()
	go a.startReporting()

	a.wg.Wait()
}

func (a *Agent) startPolling() {
	defer a.wg.Done()

	for {
		err := a.Stats.Poll()
		if err != nil {
			logging.LogError(err)
		}

		time.Sleep(utils.IntToDuration(a.Config.PollInterval))
	}
}

func (a *Agent) startReporting() {
	defer a.wg.Done()

	for {
		time.Sleep(utils.IntToDuration(a.Config.ReportInterval))

		a.reportStats()
	}
}

func (a *Agent) reportStats() {
	logging.LogInfo("reporting stats ... ")

	// agent continues polling while report is in progress, take snapshot?
	snapshot := *a.Stats

	a.Exporter.
		Add("Alloc", snapshot.Runtime.Alloc).
		Add("BuckHashSys", snapshot.Runtime.BuckHashSys).
		Add("Frees", snapshot.Runtime.Frees).
		Add("GCCPUFraction", snapshot.Runtime.GCCPUFraction).
		Add("GCSys", snapshot.Runtime.GCSys).
		Add("HeapAlloc", snapshot.Runtime.HeapAlloc).
		Add("HeapIdle", snapshot.Runtime.HeapIdle).
		Add("HeapInuse", snapshot.Runtime.HeapInuse).
		Add("HeapObjects", snapshot.Runtime.HeapObjects).
		Add("HeapReleased", snapshot.Runtime.HeapReleased).
		Add("HeapSys", snapshot.Runtime.HeapSys).
		Add("LastGC", snapshot.Runtime.LastGC).
		Add("Lookups", snapshot.Runtime.Lookups).
		Add("MCacheInuse", snapshot.Runtime.MCacheInuse).
		Add("MCacheSys", snapshot.Runtime.MCacheSys).
		Add("MSpanInuse", snapshot.Runtime.MSpanInuse).
		Add("MSpanSys", snapshot.Runtime.MSpanSys).
		Add("Mallocs", snapshot.Runtime.Mallocs).
		Add("NextGC", snapshot.Runtime.NextGC).
		Add("NumForcedGC", snapshot.Runtime.NumForcedGC).
		Add("NumGC", snapshot.Runtime.NumGC).
		Add("OtherSys", snapshot.Runtime.OtherSys).
		Add("PauseTotalNs", snapshot.Runtime.PauseTotalNs).
		Add("StackInuse", snapshot.Runtime.StackInuse).
		Add("StackSys", snapshot.Runtime.StackSys).
		Add("Sys", snapshot.Runtime.Sys).
		Add("TotalAlloc", snapshot.Runtime.TotalAlloc)

	a.Exporter.
		Add("RandomValue", snapshot.RandomValue)

	a.Exporter.
		Add("PollCount", snapshot.PollCount)

	err := a.Exporter.Send()
	if err != nil {
		logging.LogErrorF("error sending metrics: %w", err)
	}

	a.Exporter.Reset()

	// because metrics.Counter adds value to itself
	a.Stats.PollCount -= snapshot.PollCount
}

func (c Config) String() string {

	str := []string{
		fmt.Sprintf("address=%s", c.Address),
		fmt.Sprintf("poll-interval=%v", c.PollInterval),
		fmt.Sprintf("report-interval=%v", c.ReportInterval),
	}

	if len(c.Secret) > 0 {
		str = append(str, fmt.Sprintf("secret=%v", c.Secret))
	}

	return "agent config: " + strings.Join(str, "; ")
}

func parseConfig(config *Config) error {
	address := config.Address
	pflag.VarP(&address, "address", "a", "address:port for HTTP API requests")

	secret := config.Secret
	pflag.VarP(&secret, "secret", "k", "a key to sign outgoing data")

	pflag.IntVarP(&config.PollInterval, "poll-interval", "p", config.PollInterval, "interval (s) for polling stats")
	pflag.IntVarP(&config.ReportInterval, "report-interval", "r", config.ReportInterval, "interval (s) for polling stats")

	pflag.Parse()

	// because VarP gets non-pointer value, set it manually
	pflag.Visit(func(f *pflag.Flag) {
		switch f.Name {
		case "address":
			config.Address = address
		case "secret":
			config.Secret = secret
		}
	})

	if err := env.Parse(config); err != nil {
		return err
	}

	return nil
}
