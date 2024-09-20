package memory

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	config "github.com/igortoigildin/go-metrics-altering/config/agent"
	"github.com/igortoigildin/go-metrics-altering/internal/models"
	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
	"go.uber.org/zap"
)

func (m *MemoryStats) SendBatchMetrics(cfg *config.ConfigAgent) {
	for {
		metrics := []models.Metrics{}
		time.Sleep(cfg.PauseDuration)
		for i, v := range m.GaugeMetrics {
			metricGauge := models.GaugeConstructor(v, i)
			metrics = append(metrics, metricGauge)
		}
		metricCounter := models.CounterConstructor(int64(m.CounterMetric))
		metrics = append(metrics, metricCounter)
		err := sendAllMetrics(cfg, metrics)
		if err != nil {
			logger.Log.Info("unexpected sending batch metrics error:", zap.Error(err))
		}
	}
}

func sendAllMetrics(cfg *config.ConfigAgent, metrics []models.Metrics) error {
	agent := resty.New()
	req := agent.R().SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").SetHeader("Accept-Encoding", "gzip")
	metricsJSON, err := json.Marshal(metrics)
	if err != nil {
		logger.Log.Info("marshalling json error:", zap.Error(err))
		return err
	}
	// signing metric value with sha256 and setting header accordingly
	if cfg.FlagHashKey != "" {
		key := []byte(cfg.FlagHashKey)
		h := hmac.New(sha256.New, key)
		h.Write(metricsJSON)
		dst := h.Sum(nil)
		req.SetHeader("HashSHA256", fmt.Sprintf("%x", dst))
	}
	var compressedRequest bytes.Buffer
	writer := gzip.NewWriter(&compressedRequest)
	_, err = writer.Write(metricsJSON)
	if err != nil {
		logger.Log.Info("error while compressing request:", zap.Error(err))
		return err
	}
	err = writer.Close()
	if err != nil {
		logger.Log.Info("error while closing gzip writer:", zap.Error(err))
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
				logger.Log.Info("timeout error, server not reachable:", zap.Error(err))
			}
			return ErrConnectionFailed
		default:
			logger.Log.Info("unexpected sending metric error via URL:", zap.Error(err))
			return err
		}
	}
	return err
}
