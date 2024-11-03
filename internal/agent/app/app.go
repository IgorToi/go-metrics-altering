package app

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"
	"time"

	config "github.com/igortoigildin/go-metrics-altering/config/agent"
	"github.com/igortoigildin/go-metrics-altering/internal/agent"
	"github.com/igortoigildin/go-metrics-altering/internal/agent/memory"
	"github.com/igortoigildin/go-metrics-altering/internal/models"
	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
)

var buildVersion string = "N/A"
var buildDate string = "N/A"
var buildCommit string = "N/A"

func Run(cfg *config.ConfigAgent) {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	logger.Log.Info("loading metrics...")

	memoryStats := memory.New()
	var wg sync.WaitGroup
	metricsChan := make(chan models.Metrics, 33)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	wg.Add(1)
	go func() {
		for {
			time.Sleep(cfg.PauseDuration)
			select {
			case <-ctx.Done():
				logger.Log.Info("Stop updating inmemory runtime metrics")
				wg.Done()
				return
			default:
				memoryStats.UpdateRunTimeStat(cfg)
			}
		}
	}()

	wg.Add(1)
	go func() {
		for {
			time.Sleep(cfg.PauseDuration)
			select {
			case <-ctx.Done():
				logger.Log.Info("Stop updating inmemory CPU RAM metrics")
				wg.Done()
				return
			default:
				memoryStats.UpdateCPURAMStat(cfg)
			}
		}
	}()

	for w := 1; w <= cfg.FlagRateLimit; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			agent.SendMetrics(metricsChan, cfg)
		}()
	}

	go func() {
		for {
			time.Sleep(cfg.PauseDuration)
			select {
			case <-ctx.Done():
				logger.Log.Info("Stop filling metrics chan with new metrics")
				close(metricsChan)
				return
			default:
				wg.Add(1)
				memoryStats.ReadMetrics(cfg, metricsChan)
				wg.Done()
			}
		}
	}()

	wg.Wait()

	logger.Log.Info("Graceful agent shutdown complete...")
}
