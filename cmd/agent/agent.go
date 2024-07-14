package main

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	config "github.com/igortoigildin/go-metrics-altering/config/agent"
	httpAgent "github.com/igortoigildin/go-metrics-altering/internal/agent"
	"github.com/igortoigildin/go-metrics-altering/internal/logger"
	"github.com/igortoigildin/go-metrics-altering/internal/models"
	"go.uber.org/zap"
)

// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Log.Fatal("error while logading config", zap.Error(err))
	}
	// start goroutine to update metrics every pollInterval
	go cfg.UpdateMetrics()
	agent := resty.New()
	durationPause := time.Duration(cfg.FlagReportInterval) * time.Second
	var metrics []models.Metrics
	for {
		time.Sleep(durationPause)
		for i, v := range cfg.Memory {
			req := agent.R()
			req.SetPathParams(map[string]string{
				"metricType":  config.GaugeType,
				"metricName":  i,
				"metricValue": strconv.FormatFloat(v, 'f', 6, 64),
			}).SetHeader("Content-Type", "text/plain")
			req.URL = config.ProtocolScheme + cfg.FlagRunAddr
			_, err := httpAgent.SendMetric(req.URL, config.GaugeType, i, strconv.FormatFloat(v, 'f', 6, 64), req)
			if err != nil {
				logger.Log.Debug("unexpected sending metric error:", zap.Error(err))
			}
			logger.Log.Info("Metric has been sent successfully")

			var metric models.Metrics
			metric.ID = i
			metric.MType = config.GaugeType
			metric.Value = &v
			metrics = append(metrics, metric)
		}
		req := agent.R()
		req.SetPathParams(map[string]string{
			"metricType":  config.CountType,
			"metricName":  config.PollCount,
			"metricValue": strconv.Itoa(cfg.Count),
		}).SetHeader("Content-Type", "text/plain")

		req.URL = config.ProtocolScheme + cfg.FlagRunAddr
		_, err := httpAgent.SendMetric(req.URL, config.CountType, config.PollCount, strconv.Itoa(cfg.Count), req)
		if err != nil {
			logger.Log.Debug("unexpected sending metric error:", zap.Error(err))
		}
		logger.Log.Info("Metric has been sent successfully")

		var metric models.Metrics
		var delta int64 = int64(cfg.Count)
		metric.ID = config.PollCount
		metric.MType = config.CountType
		metric.Delta = &delta
		metrics = append(metrics, metric)
		metricsJSON, err  := json.Marshal(metrics)
		if err != nil {
			logger.Log.Debug("unexpected sending metric error:", zap.Error(err))
		}
		_, err = req.SetBody(metricsJSON).SetHeader("Content-Type", "text/plain").Post(req.URL)
		if err != nil {
			logger.Log.Debug("unexpected sending metric error:", zap.Error(err))
		}

	}
}
