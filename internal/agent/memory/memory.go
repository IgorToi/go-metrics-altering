package memory

import (
	"errors"
	"math/rand/v2"
	"runtime"
	"time"

	config "github.com/igortoigildin/go-metrics-altering/config/agent"
)

var (
    ErrConnectionFailed = errors.New("connection failed")
)

type MemoryStats struct {
    GaugeMetrics        map[string]float64
    CounterMetric       int
    Rtm                runtime.MemStats
}

func NewMemoryStats() *MemoryStats{
    return &MemoryStats{
        GaugeMetrics: make(map[string]float64),
    }
}

func (m *MemoryStats) UpdateMetrics(cfg *config.ConfigAgent) {
    PauseDuration := time.Duration(cfg.FlagPollInterval) * time.Second
    for {
        time.Sleep(PauseDuration)
        runtime.ReadMemStats(&m.Rtm)
        m.GaugeMetrics["Alloc"] = float64(m.Rtm.Alloc)
        m.GaugeMetrics["BuckHashSys"] = float64(m.Rtm.BuckHashSys)
        m.GaugeMetrics["Frees"] = float64(m.Rtm.Frees)
        m.GaugeMetrics["GCCPUFraction"] = float64(m.Rtm.GCCPUFraction)
        m.GaugeMetrics["GCSys"] = float64(m.Rtm.GCSys)
        m.GaugeMetrics["HeapAlloc"] = float64(m.Rtm.HeapAlloc)
        m.GaugeMetrics["HeapIdle"] = float64(m.Rtm.HeapIdle)
        m.GaugeMetrics["HeapInuse"] = float64(m.Rtm.HeapInuse)
        m.GaugeMetrics["HeapObjects"] = float64(m.Rtm.HeapObjects)
        m.GaugeMetrics["HeapReleased"] = float64(m.Rtm.HeapReleased)
        m.GaugeMetrics["HeapSys"] = float64(m.Rtm.HeapSys)
        m.GaugeMetrics["LastGC"] = float64(m.Rtm.LastGC)
        m.GaugeMetrics["Lookups"] = float64(m.Rtm.Lookups)
        m.GaugeMetrics["MCacheInuse"] = float64(m.Rtm.MCacheInuse)
        m.GaugeMetrics["MCacheSys"] = float64(m.Rtm.MCacheSys)
        m.GaugeMetrics["MSpanInuse"] = float64(m.Rtm.MSpanInuse)
        m.GaugeMetrics["MSpanSys"] = float64(m.Rtm.MSpanSys)
        m.GaugeMetrics["Mallocs"] = float64(m.Rtm.Mallocs)
        m.GaugeMetrics["NextGC"] = float64(m.Rtm.NextGC)
        m.GaugeMetrics["NumForcedGC"] = float64(m.Rtm.NumForcedGC)
        m.GaugeMetrics["NumGC"] = float64(m.Rtm.NumGC)
        m.GaugeMetrics["OtherSys"] = float64(m.Rtm.OtherSys)
        m.GaugeMetrics["NextGC"] = float64(m.Rtm.NextGC)
        m.GaugeMetrics["NumForcedGC"] = float64(m.Rtm.NumForcedGC)
        m.GaugeMetrics["NumGC"] = float64(m.Rtm.NumGC)
        m.GaugeMetrics["OtherSys"] = float64(m.Rtm.OtherSys)
        m.GaugeMetrics["PauseTotalNs"] = float64(m.Rtm.PauseTotalNs)
        m.GaugeMetrics["StackInuse"] = float64(m.Rtm.StackInuse)
        m.GaugeMetrics["StackSys"] = float64(m.Rtm.StackSys)
        m.GaugeMetrics["Sys"] = float64(m.Rtm.StackSys)
        m.GaugeMetrics["TotalAlloc"] = float64(m.Rtm.TotalAlloc)
        m.GaugeMetrics["RandomValue"] = rand.Float64()
        m.CounterMetric++
    }
}

