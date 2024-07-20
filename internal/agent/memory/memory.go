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
	GaugeMetrics 		map[string]float64
	CounterMetric		int
	Rtm                runtime.MemStats
}

func NewMemoryStats() *MemoryStats{
	return &MemoryStats{
		GaugeMetrics: make(map[string]float64),
	}
}

func (c *MemoryStats) UpdateMetrics(cfg *config.ConfigAgent) {
	PauseDuration := time.Duration(cfg.FlagPollInterval) * time.Second
	for {
		time.Sleep(PauseDuration)
		runtime.ReadMemStats(&c.Rtm)
		c.GaugeMetrics["Alloc"] = float64(c.Rtm.Alloc)
		c.GaugeMetrics["BuckHashSys"] = float64(c.Rtm.BuckHashSys)
		c.GaugeMetrics["Frees"] = float64(c.Rtm.Frees)
		c.GaugeMetrics["GCCPUFraction"] = float64(c.Rtm.GCCPUFraction)
		c.GaugeMetrics["GCSys"] = float64(c.Rtm.GCSys)
		c.GaugeMetrics["HeapAlloc"] = float64(c.Rtm.HeapAlloc)
		c.GaugeMetrics["HeapIdle"] = float64(c.Rtm.HeapIdle)
		c.GaugeMetrics["HeapInuse"] = float64(c.Rtm.HeapInuse)
		c.GaugeMetrics["HeapObjects"] = float64(c.Rtm.HeapObjects)
		c.GaugeMetrics["HeapReleased"] = float64(c.Rtm.HeapReleased)
		c.GaugeMetrics["HeapSys"] = float64(c.Rtm.HeapSys)
		c.GaugeMetrics["LastGC"] = float64(c.Rtm.LastGC)
		c.GaugeMetrics["Lookups"] = float64(c.Rtm.Lookups)
		c.GaugeMetrics["MCacheInuse"] = float64(c.Rtm.MCacheInuse)
		c.GaugeMetrics["MCacheSys"] = float64(c.Rtm.MCacheSys)
		c.GaugeMetrics["MSpanInuse"] = float64(c.Rtm.MSpanInuse)
		c.GaugeMetrics["MSpanSys"] = float64(c.Rtm.MSpanSys)
		c.GaugeMetrics["Mallocs"] = float64(c.Rtm.Mallocs)
		c.GaugeMetrics["NextGC"] = float64(c.Rtm.NextGC)
		c.GaugeMetrics["NumForcedGC"] = float64(c.Rtm.NumForcedGC)
		c.GaugeMetrics["NumGC"] = float64(c.Rtm.NumGC)
		c.GaugeMetrics["OtherSys"] = float64(c.Rtm.OtherSys)
		c.GaugeMetrics["NextGC"] = float64(c.Rtm.NextGC)
		c.GaugeMetrics["NumForcedGC"] = float64(c.Rtm.NumForcedGC)
		c.GaugeMetrics["NumGC"] = float64(c.Rtm.NumGC)
		c.GaugeMetrics["OtherSys"] = float64(c.Rtm.OtherSys)
		c.GaugeMetrics["PauseTotalNs"] = float64(c.Rtm.PauseTotalNs)
		c.GaugeMetrics["StackInuse"] = float64(c.Rtm.StackInuse)
		c.GaugeMetrics["StackSys"] = float64(c.Rtm.StackSys)
		c.GaugeMetrics["Sys"] = float64(c.Rtm.StackSys)
		c.GaugeMetrics["TotalAlloc"] = float64(c.Rtm.TotalAlloc)
		c.GaugeMetrics["RandomValue"] = rand.Float64()
		c.CounterMetric++
	}
}
