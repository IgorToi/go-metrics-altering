package main

import (
	config "github.com/igortoigildin/go-metrics-altering/config/agent"
	"github.com/igortoigildin/go-metrics-altering/internal/models"
)

func prepareMetricBody(cfg *config.ConfigAgent, metricName string) []models.Metrics {
	var metrics []models.Metrics
	valueGauge := cfg.Memory[metricName]
	metric := models.Metrics{
		ID:    metricName,
		MType: config.GaugeType,
		Value: &valueGauge,
	}
	metrics = append(metrics, metric)

	valueDelta := int64(cfg.Count)
	metric = models.Metrics{
		ID:    config.PollCount,
		MType: config.CountType,
		Delta: &valueDelta,
	}
	metrics = append(metrics, metric)
	return metrics
}
