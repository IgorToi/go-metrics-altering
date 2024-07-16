package main

import (
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	config "github.com/igortoigildin/go-metrics-altering/config/agent"
	httpAgent "github.com/igortoigildin/go-metrics-altering/internal/agent"
	"github.com/igortoigildin/go-metrics-altering/internal/logger"
	"go.uber.org/zap"
)

// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Log.Fatal("error while logading config", zap.Error(err))
	}
	if err := logger.Initialize(cfg.FlagLogLevel); err != nil {
		logger.Log.Fatal("error while initializing logger", zap.Error(err))
	}
	// start goroutine to update metrics every pollInterval
	go cfg.UpdateMetrics()
	// start goroutine to send metrics via /update/ - variant 2
	go httpAgent.SendAllMetrics(cfg)
	agent := resty.New()
	durationPause := time.Duration(cfg.FlagReportInterval) * time.Second
	for {
		time.Sleep(durationPause)
		for i, v := range cfg.Memory {
			req := agent.R()
			// preparing and sending metric via url
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
