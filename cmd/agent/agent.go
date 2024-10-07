package main

import (
	"sync"

	config "github.com/igortoigildin/go-metrics-altering/config/agent"
	"github.com/igortoigildin/go-metrics-altering/internal/agent"
	"github.com/igortoigildin/go-metrics-altering/internal/agent/memory"
	"github.com/igortoigildin/go-metrics-altering/internal/models"
	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Log.Fatal("error while logading config", zap.Error(err))
	}
	if err := logger.Initialize(cfg.FlagLogLevel); err != nil {
		logger.Log.Fatal("error while initializing logger", zap.Error(err))
	}

	logger.Log.Info("loading metrics...")

	memoryStats := memory.NewMemoryStats()
	var wg sync.WaitGroup
	metricsChan := make(chan models.Metrics, 33)

	wg.Add(1)
	go func() {
		defer wg.Done()
		memoryStats.UpdateRunTimeStat(cfg)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		memoryStats.UpdateCPURAMStat(cfg)
	}()

	for w := 1; w <= cfg.FlagRateLimit; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			agent.SendMetrics(metricsChan, cfg)
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		memoryStats.ReadMetrics(cfg, metricsChan)
	}()

	wg.Wait()
}
