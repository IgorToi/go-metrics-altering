package agent

import (
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	config "github.com/igortoigildin/go-metrics-altering/config/agent"
	"github.com/igortoigildin/go-metrics-altering/internal/logger"
	"go.uber.org/zap"
)

const urlParams = "/update/{metricType}/{metricName}/{metricValue}"

// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func sendURLmetrics(cfg *config.ConfigAgent) error {
	agent := resty.New()
	for {
		time.Sleep(cfg.PauseDuration)
		for i := range cfg.Memory {
			err := sendURLGauge(cfg, agent, i)
			if err != nil {
				return err
			}
		}
		err := sendURLCounter(cfg, agent)
		if err != nil {
			return err
		}
		logger.Log.Info("All URL metrics sent successfully")
	}
}

func sendURLCounter(cfg *config.ConfigAgent, agent *resty.Client) error {
	req := agent.R().SetPathParams(map[string]string{
		"metricType":  config.CountType,
		"metricName":  config.PollCount,
		"metricValue": strconv.Itoa(cfg.Count),
	}).SetHeader("Content-Type", "text/plain")
	_, err := req.Post(cfg.URL + urlParams)
	if err != nil {
		switch {
		case os.IsTimeout(err):
			for _, delay := range []time.Duration{time.Second, 2 * time.Second, 3 * time.Second} {
				time.Sleep(delay)
				if _, err = req.Post(cfg.URL + urlParams); err == nil {
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

func sendURLGauge(cfg *config.ConfigAgent, agent *resty.Client, metricName string) error {
	req := agent.R().SetPathParams(map[string]string{
		"metricType":  config.GaugeType,
		"metricName":  metricName,
		"metricValue": strconv.FormatFloat(cfg.Memory[metricName], 'f', 6, 64),
	}).SetHeader("Content-Type", "text/plain")
	_, err := req.Post(cfg.URL + urlParams)
	if err != nil {
		switch {
		case os.IsTimeout(err):
			for _, delay := range []time.Duration{time.Second, 2 * time.Second, 3 * time.Second} {
				time.Sleep(delay)
				if _, err = req.Post(cfg.URL + urlParams); err == nil {
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
