package agent

import (
	"context"
	"fmt"

	"github.com/ex0rcist/metflix/internal/metrics"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type SystemStats struct {
	TotalMemory metrics.Gauge
	FreeMemory  metrics.Gauge

	CPUutilization []metrics.Gauge
}

func (s *SystemStats) Poll(ctx context.Context) error {
	vMem, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("mem.VirtualMemory: %w", err)
	}

	s.TotalMemory = metrics.Gauge(vMem.Total)
	s.FreeMemory = metrics.Gauge(vMem.Free)

	utilisation, err := cpu.PercentWithContext(ctx, 0, true)
	if err != nil {
		return fmt.Errorf("cpu.PercentWithContext: %w", err)
	}

	s.CPUutilization = []metrics.Gauge{}
	for _, u := range utilisation {
		s.CPUutilization = append(s.CPUutilization, metrics.Gauge(u))
	}

	return nil
}
