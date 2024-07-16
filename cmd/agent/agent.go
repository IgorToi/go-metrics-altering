package main

import (
	"encoding/json"
	"errors"
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
		if errors.Is(err, config.ErrParsingFlag) {
			logger.Log.Fatal("error while parsing flag", zap.Error(err))
		} else {
			logger.Log.Fatal("unexpected flag error", zap.Error(err))
		}
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
			// Varint 1 with memory
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
			// Varint 2 with DB
			metricGauge := models.Metrics{
				ID:    i,
				MType: config.GaugeType,
				Value: &v,
			}
			metrics = append(metrics, metricGauge)
			delta := int64(cfg.Count)
			metricCounter := models.Metrics{
				ID:    config.PollCount,
				MType: config.CountType,
				Delta: &delta,
			}
			metrics = append(metrics, metricCounter)
			metricsJSON, err := json.Marshal(metrics)
			if err != nil {
				logger.Log.Debug("unexpected sending metric error:", zap.Error(err))
			}
			_, err = req.SetBody(metricsJSON).SetHeader("Content-Type", "application/json").Post(req.URL + "/updates/")
			if err != nil {
				// urlErr := err.(*url.Error)
				// if urlErr != nil {
				// attempt to send metric again
					logger.Log.Debug("unexpected sending metric error:", zap.Error(err))
					for n, t := 1, 1; n <= 3; n++ {
						time.Sleep(time.Duration(t) * time.Second)
						if _, err = req.Post(req.URL + "/updates/"); err == nil {
							logger.Log.Info("Metric has been sent successfully")
							break
						}
						t += 2
					}
			}

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
	}
}
