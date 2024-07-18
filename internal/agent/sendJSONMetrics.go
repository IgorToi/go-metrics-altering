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
    for {
        time.Sleep(cfg.PauseDuration)
        for i := range cfg.Memory {
            err := SendJSONGauge(cfg, i)
            if err != nil {
                logger.Log.Debug("unexpected sending json metric error:", zap.Error(err))
                return
            }
        }
        err := SendJSONCounter(cfg)
        if err != nil {
            logger.Log.Debug("unexpected sending json metric error:", zap.Error(err))
            return
        }
    }
}

func SendJSONGauge(cfg *config.ConfigAgent, metricName string) error {
    agent := resty.New()
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
    req.Method = resty.MethodPost
    _, err = req.SetBody(metricsJSON).Post(cfg.URL + "/update/")
    if err != nil {
        // send again n times if timeout error
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
            logger.Log.Debug("unexpected sending metric error:", zap.Error(err))
            return err
        }
    }
    return nil
}

func SendJSONCounter(cfg *config.ConfigAgent) error {
    agent := resty.New()
    valueDelta := int64(cfg.Count)
    metric := models.Metrics{
        ID:    config.PollCount,
        MType: config.CountType,
        Delta: &valueDelta,
    }
    req := agent.R().SetHeader("Content-Type", "application/json")
    metricJSON, err := json.Marshal(metric)
    if err != nil {
        logger.Log.Debug("marshalling json error:", zap.Error(err))
        return err
    }
    _, err = req.SetBody(metricJSON).Post(cfg.URL + "/update/")
    if err != nil {
        // send again n times if timeout error
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
    return nil
}

