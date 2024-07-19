package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	config "github.com/igortoigildin/go-metrics-altering/config/agent"
	"github.com/igortoigildin/go-metrics-altering/internal/logger"
	"github.com/igortoigildin/go-metrics-altering/internal/models"
	"go.uber.org/zap"
)

func sendBatchMetrics(cfg *config.ConfigAgent) {
	agent := resty.New()
	for {
		metrics := []models.Metrics{}
		time.Sleep(cfg.PauseDuration)
		for i, v := range cfg.Memory {
			metricGauge := models.Metrics{
				ID:    i,
				MType: config.GaugeType,
				Value: &v,
			}
			metrics = append(metrics, metricGauge)
		}
		countDelta := int64(cfg.Count)
		metricCounter := models.Metrics{
			ID:    config.PollCount,
			MType: config.CountType,
			Delta: &countDelta,
		}
		metrics = append(metrics, metricCounter)
		err := sendAllMetrics(cfg, metrics, agent)
		if err != nil {
			logger.Log.Debug("unexpected sending batch metrics error:", zap.Error(err))
			return
		}
		logger.Log.Info("Metrics batch sent successfully")
	}
}

func sendAllMetrics(cfg *config.ConfigAgent, metrics []models.Metrics, agent *resty.Client) error {
	req := agent.R().SetHeader("Content-Type", "application/json").SetHeader("Content-Encoding", "gzip")
	metricsJSON, err := json.Marshal(metrics)
	if err != nil {
		logger.Log.Debug("marshalling json error:", zap.Error(err))
		return err
	}
	var compressedRequest bytes.Buffer
	writer := gzip.NewWriter(&compressedRequest)
	_, err = writer.Write(metricsJSON)
	if err != nil {
		logger.Log.Debug("error while compressing request:", zap.Error(err))
		return err
	}
	err = writer.Close()
	if err != nil {
		logger.Log.Debug("error while closing gzip writer:", zap.Error(err))
		return err
	}
	_, err = req.SetBody(compressedRequest.Bytes()).Post(cfg.URL + "/updates/")
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
