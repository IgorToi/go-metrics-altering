// Package memory accumulates updated metrics,
// which is being used and consumed by agent.

package memory

import (
	"testing"
	"time"

	config "github.com/igortoigildin/go-metrics-altering/config/agent"
	"github.com/igortoigildin/go-metrics-altering/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestMemoryStats_UpdateRunTimeStat(t *testing.T) {
	m := New()
	cfg, _ := config.LoadConfig()
	cfg.PauseDuration = 0
	go m.UpdateRunTimeStat(cfg)

	time.Sleep(1 * time.Second)

	assert.Equal(t, 28, len(m.GaugeMetrics))
}

func TestMemoryStats_UpdateCPURAMStat(t *testing.T) {
	m := New()

	cfg := config.ConfigAgent{}
	cfg.PauseDuration = 0

	go m.UpdateCPURAMStat(&cfg)

	time.Sleep(1 * time.Second)

	assert.Equal(t, 3, len(m.GaugeMetrics))
}

func TestMemoryStats_ReadMetrics(t *testing.T) {
	mem := New()
	mem.CounterMetric = 2
	mem.GaugeMetrics = map[string]float64{"first": float64(5)}

	ch := make(chan models.Metrics, 2)
	cfg := config.ConfigAgent{}
	mem.ReadMetrics(&cfg, ch)

	fm := <-ch
	assert.Equal(t, "first", fm.ID)

	sm := <-ch
	assert.NotNil(t, sm.Delta)
}
