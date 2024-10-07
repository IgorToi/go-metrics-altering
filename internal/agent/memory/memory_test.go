// Package memory accumulates updated metrics,
// which is being used and consumed by agent.

package memory

import (
	"testing"
	"time"

	config "github.com/igortoigildin/go-metrics-altering/config/agent"
	"github.com/stretchr/testify/assert"
)

func TestMemoryStats_UpdateRunTimeStat(t *testing.T) {
	m := NewMemoryStats()
	cfg, _ := config.LoadConfig()
	cfg.PauseDuration = 0
	go m.UpdateRunTimeStat(cfg)

	time.Sleep(1 * time.Second)

	assert.Equal(t, 28, len(m.GaugeMetrics))
}

func TestMemoryStats_UpdateCPURAMStat(t *testing.T) {
	m := NewMemoryStats()

	cfg := config.ConfigAgent{}
	cfg.PauseDuration = 0

	go m.UpdateCPURAMStat(&cfg)

	time.Sleep(1 * time.Second)

	assert.Equal(t, 3, len(m.GaugeMetrics))
}
