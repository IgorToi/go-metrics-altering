// Package memory accumulates updated metrics,
// which is being used and consumed by agent.

package memory

import (
	"runtime"
	"sync"
	"testing"

	config "github.com/igortoigildin/go-metrics-altering/config/agent"
)

func TestMemoryStats_UpdateRunTimeStat(t *testing.T) {
	type fields struct {
		GaugeMetrics  map[string]float64
		CounterMetric int
		RunTimeMem    *runtime.MemStats
		rwm           sync.RWMutex
	}
	type args struct {
		cfg *config.ConfigAgent
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryStats{
				GaugeMetrics:  tt.fields.GaugeMetrics,
				CounterMetric: tt.fields.CounterMetric,
				RunTimeMem:    tt.fields.RunTimeMem,
				rwm:           tt.fields.rwm,
			}
			m.UpdateRunTimeStat(tt.args.cfg)
		})
	}
}
