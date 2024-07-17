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

func sendJSONMetrics(cfg *config.ConfigAgent) {
	agent := resty.New()
	for {
		time.Sleep(cfg.PauseDuration)
		for i := range cfg.Memory {
			err := SendJSONGauge(cfg, i, agent)
			if err != nil {
				logger.Log.Fatal("unexpected sending json metric error:", zap.Error(err))
				return
			}
		}
		err := SendJSONCounter(cfg, agent)
		if err != nil {
			logger.Log.Fatal("unexpected sending json metric error:", zap.Error(err))
			return
		}
	}
}

func SendJSONGauge(cfg *config.ConfigAgent, metricName string, agent *resty.Client) error {
	valueGauge := cfg.Memory[metricName]
	metric := models.Metrics{
		ID:    metricName,
		MType: config.GaugeType,
		Value: &valueGauge,
	}
	req := agent.R().SetHeader("Content-Type", "application/json")
	metricsJSON, err := json.Marshal(metric)
	if err != nil {
		logger.Log.Debug("marshalling json error:", zap.Error(err))
		return err
	}
	req.URL = cfg.URL + "/update/"
	_, err = req.SetBody(metricsJSON).Post(req.URL)
	if err != nil {
		//send again n times if timeout error
		switch {
		case os.IsTimeout(err):
			for _, delay := range []time.Duration{time.Second, 2 * time.Second, 3 * time.Second} {
				time.Sleep(delay)
				if _, err = req.Post(req.URL); err == nil {
					break
				}
				logger.Log.Debug("timeout error, server not reachable:", zap.Error(err))
			}
			return ErrConnectionFailed
		default:
			logger.Log.Debug("unexpected sending metric error via URL:", zap.Error(err))
			return err
		}
	}
	return err
}

func SendJSONCounter(cfg *config.ConfigAgent, agent *resty.Client) error {
	valueDelta := int64(cfg.Count)
	metric := models.Metrics{
		ID:    config.PollCount,
		MType: config.CountType,
		Delta: &valueDelta,
	}
	req := agent.R().SetHeader("Content-Type", "application/json")
	metricsJSON, err := json.Marshal(metric)
	if err != nil {
		logger.Log.Debug("marshalling json error:", zap.Error(err))
		return err
	}
	req.URL = cfg.URL + "/update/"
	_, err = req.SetBody(metricsJSON).Post(req.URL)
	if err != nil {
		//send again n times if timeout error
		switch {
		case os.IsTimeout(err):
			for _, delay := range []time.Duration{time.Second, 2 * time.Second, 3 * time.Second} {
				time.Sleep(delay)
				if _, err = req.Post(req.URL); err == nil {
					break
				}
				logger.Log.Debug("timeout error, server not reachable:", zap.Error(err))
			}
			return ErrConnectionFailed
		default:
			logger.Log.Debug("unexpected sending metric error via URL:", zap.Error(err))
			return err
		}
	}
	return err
}
