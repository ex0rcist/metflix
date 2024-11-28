package agent

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ex0rcist/metflix/internal/agent/exporter"
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/security"
	"github.com/ex0rcist/metflix/internal/utils"
)

const shutdownTimeout = 60 * time.Second

// Metric collecting agent (mr. Bond?).
type Agent struct {
	Config    *Config
	Stats     *Stats
	Exporter  exporter.Exporter
	interrupt chan os.Signal
	wg        sync.WaitGroup
}

// Constructor.
func New() (*Agent, error) {
	config, err := NewConfig()
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

// Shutdown agent
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

func newMetricsExporter(ctx context.Context, config *Config) (exporter.Exporter, error) {
	var (
		exp    exporter.Exporter
		signer security.Signer
		err    error
	)

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

	switch {
	case config.RateLimit > 0:
		exp = exporter.NewLimitedExporter(ctx, &config.Address, signer, config.RateLimit, publicKey)
	default:
		exp = exporter.NewBatchExporter(ctx, &config.Address, signer, publicKey)
	}

	return exp, err
}
