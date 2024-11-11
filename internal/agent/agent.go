package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/ex0rcist/metflix/internal/agent/exporter"
	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/security"
	"github.com/ex0rcist/metflix/internal/utils"
	"github.com/spf13/pflag"
)

const shutdownTimeout = 60 * time.Second

// Metric collecting agent (mr. Bond?).
type Agent struct {
	Config   *Config
	Stats    *Stats
	Exporter exporter.Exporter

	interrupt chan os.Signal

	wg sync.WaitGroup
}

// Agent config.
type Config struct {
	Address        entities.Address  `env:"ADDRESS" json:"address"`
	PollInterval   int               `env:"POLL_INTERVAL" json:"poll_interval"`
	ReportInterval int               `env:"REPORT_INTERVAL" json:"report_interval"`
	RateLimit      int               `env:"RATE_LIMIT" json:"-"`
	Secret         entities.Secret   `env:"KEY" json:"key"`
	PublicKeyPath  entities.FilePath `env:"CRYPTO_KEY" json:"crypto_key"`
}

// Constructor.
func New() (*Agent, error) {
	config := &Config{
		Address:        "0.0.0.0:8080",
		PollInterval:   2,
		ReportInterval: 10,
		RateLimit:      -1,
	}

	err := parseConfig(config)
	if err != nil {
		return nil, err
	}

	return &Agent{
		Config:    config,
		Stats:     NewStats(),
		interrupt: make(chan os.Signal, 1),
	}, nil
}

// Run agent.
func (a *Agent) Run() error {
	logging.LogInfo(a.Config.String())
	logging.LogInfo("agent ready")

	ctx, cancelBackgroundTasks := context.WithCancel(context.Background())
	defer cancelBackgroundTasks()

	exporter, err := newMetricsExporter(ctx, a.Config)
	if err != nil {
		return err
	}
	a.Exporter = exporter

	signal.Notify(a.interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	a.wg.Add(2)

	go func() {
		defer a.wg.Done()
		a.startPolling(ctx)
	}()

	go func() {
		defer a.wg.Done()
		a.startReporting(ctx)
	}()

	<-a.interrupt

	logging.LogInfo("shutting down agent...")
	cancelBackgroundTasks()

	stopped := make(chan struct{})

	stopCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	go func() {
		defer close(stopped)
		a.wg.Wait()
	}()

	select {
	case <-stopped:
		logging.LogInfo("agent shutdown successful")

	case <-stopCtx.Done():
		logging.LogWarn("exceeded shutdown timeout, force exit")
	}

	return nil
}

func (a *Agent) Shutdown() {
	a.interrupt <- os.Interrupt
}

func (a *Agent) startPolling(ctx context.Context) {
	ticker := time.NewTicker(utils.IntToDuration(a.Config.PollInterval))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			func() {
				err := a.Stats.Poll(ctx)
				if err != nil {
					logging.LogError(err)
				}
			}()

		case <-ctx.Done():
			logging.LogInfo("shutting down metrics polling")
			return
		}
	}
}

func (a *Agent) startReporting(ctx context.Context) {
	ticker := time.NewTicker(utils.IntToDuration(a.Config.ReportInterval))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			a.reportStats()

		case <-ctx.Done():
			logging.LogInfo("shutting down metrics reporting")
			return
		}
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

	a.Exporter.
		Add("FreeMemory", snapshot.System.FreeMemory).
		Add("TotalMemory", snapshot.System.TotalMemory)

	for i, u := range snapshot.System.CPUutilization {
		a.Exporter.Add(fmt.Sprintf("CPUutilization%d", i+1), u)
	}

	err := a.Exporter.Send()
	if err != nil {
		logging.LogErrorF("error sending metrics: %w", err)
	}

	a.Exporter.Reset()

	// because metrics.Counter adds value to itself
	a.Stats.PollCount -= snapshot.PollCount
}

// Stringer.
func (c Config) String() string {
	str := []string{
		fmt.Sprintf("address=%s", c.Address),
		fmt.Sprintf("poll-interval=%v", c.PollInterval),
		fmt.Sprintf("report-interval=%v", c.ReportInterval),
		fmt.Sprintf("rate-limit=%v", c.RateLimit),
	}

	if len(c.Secret) > 0 {
		str = append(str, fmt.Sprintf("secret=%v", c.Secret))
	}

	if len(c.PublicKeyPath) > 0 {
		str = append(str, fmt.Sprintf("public-key=%v", c.PublicKeyPath))
	}

	return "agent config: " + strings.Join(str, "; ")
}

func detectExporterKind(c *Config) string {
	var ek string

	switch {
	case c.RateLimit > 0:
		ek = exporter.KindLimited
	default:
		ek = exporter.KindBatch
	}

	return ek
}

func newMetricsExporter(ctx context.Context, config *Config) (exporter.Exporter, error) {
	var exp exporter.Exporter
	var signer security.Signer
	var err error

	if len(config.Secret) > 0 {
		signer = security.NewSignerService(config.Secret)
	}

	var publicKey security.PublicKey
	if len(config.PublicKeyPath) != 0 {
		publicKey, err = security.NewPublicKey(config.PublicKeyPath)
		if err != nil {
			return nil, err
		}
	}

	exporterKind := detectExporterKind(config)

	switch exporterKind {
	case exporter.KindLimited:
		exp = exporter.NewLimitedExporter(ctx, &config.Address, signer, config.RateLimit, publicKey)
	case exporter.KindBatch:
		exp = exporter.NewBatchExporter(ctx, &config.Address, signer, publicKey)
	default:
		exp, err = nil, fmt.Errorf("unknown exporter type")
	}

	return exp, err
}

func parseConfig(config *Config) error {
	err := tryLoadJSONConfig(config)
	if err != nil {
		return err
	}

	address := config.Address
	pflag.VarP(&address, "address", "a", "address:port for HTTP API requests")

	secret := config.Secret
	pflag.VarP(&secret, "secret", "k", "a key to sign outgoing data")

	publicKeyPath := config.PublicKeyPath
	pflag.VarP(&publicKeyPath, "crypto-key", "", "path to public key to encrypt agent -> server communications")

	configPath := entities.FilePath("") // register var for compatibility
	pflag.VarP(&configPath, "config", "c", "path to configuration file in JSON format")

	pflag.IntVarP(&config.PollInterval, "poll-interval", "p", config.PollInterval, "interval (s) for polling stats")
	pflag.IntVarP(&config.ReportInterval, "report-interval", "r", config.ReportInterval, "interval (s) for polling stats")
	pflag.IntVarP(&config.RateLimit, "rate-limit", "l", config.RateLimit, "number of max simultaneous requests to server")

	pflag.Parse()

	pflag.Visit(func(f *pflag.Flag) {
		switch f.Name {
		case "address":
			config.Address = address
		case "secret":
			config.Secret = secret
		case "crypto-key":
			config.PublicKeyPath = publicKeyPath
		}
	})

	if err := env.Parse(config); err != nil {
		return err
	}

	return nil
}

func tryLoadJSONConfig(dst *Config) error {
	configArg := os.Getenv("CONFIG")

	// args is higher prior
	for i, arg := range os.Args {
		if (arg == "-c" || arg == "--config") && i+1 < len(os.Args) {
			configArg = os.Args[i+1]
			break
		}
	}

	if len(configArg) > 0 {
		err := loadConfigFromFile(entities.FilePath(configArg), dst)
		if err != nil {
			return err
		}
	}

	return nil
}

func loadConfigFromFile(src entities.FilePath, dst *Config) error {
	data, err := os.ReadFile(src.String())
	if err != nil {
		return fmt.Errorf("agent.loadConfigFromFile - os.ReadFile: %w", err)
	}

	if err := json.Unmarshal(data, dst); err != nil {
		return fmt.Errorf("agent.loadConfigFromFile - json.Unmarshal: %w", err)
	}

	return nil
}
