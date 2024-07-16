package agent

import (
	"encoding/json"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	config "github.com/igortoigildin/go-metrics-altering/config/agent"
	"github.com/igortoigildin/go-metrics-altering/internal/logger"
	"github.com/igortoigildin/go-metrics-altering/internal/models"
	"go.uber.org/zap"
)

func SendMetric(requestURL, metricType, metricName, metricValue string, req *resty.Request) (*resty.Response, error) {
    return req.Post(req.URL + "/update/{metricType}/{metricName}/{metricValue}")
}

func PrepareMetricBodyNew(cfg *config.ConfigAgent, metricName string) models.Metrics {
    var metric models.Metrics
    switch metricName {
    case config.PollCount:
        valueDelta := int64(cfg.Count)
        metric = models.Metrics{
            ID:    config.PollCount,
            MType: config.CountType,
            Delta: &valueDelta,
        }
    default:
        valueGauge := cfg.Memory[metricName]
        metric = models.Metrics{
            ID:    metricName,
            MType: config.GaugeType,
            Value: &valueGauge,
        }
    }
    agent := resty.New()
    req := agent.R().SetHeader("Content-Type", "application/json")
    req.URL = config.ProtocolScheme + cfg.FlagRunAddr
    metricsJSON, err := json.Marshal(metric)
    if err != nil {
        logger.Log.Debug("unexpected sending metric error:", zap.Error(err))
    }
    req.URL = req.URL + "/update/"
    _, err = req.SetBody(metricsJSON).Post(req.URL)
    if err != nil {
        if os.IsTimeout(err) {
            for n, t := 1, 1; n <= 3; n++ {
                time.Sleep(time.Duration(t) * time.Second)
                if _, err = req.Post(req.URL); err == nil {
                    break
                }
                t += 2
            }
        }
        logger.Log.Debug("unexpected sending metric error:", zap.Error(err))
    }
    return metric
}

func SendAllMetrics(cfg *config.ConfigAgent) {
    var metrics []models.Metrics
    agent := resty.New()
    req := agent.R().SetHeader("Content-Type", "application/json")
    req.URL = config.ProtocolScheme + cfg.FlagRunAddr
    durationPause := time.Duration(cfg.FlagReportInterval) * time.Second
    for {
        metrics = metrics[:0]
        time.Sleep(durationPause)
        for i := range cfg.Memory {
            _ = PrepareMetricBodyNew(cfg, i)
            _ = PrepareMetricBodyNew(cfg, config.PollCount)
        }
    }
}

