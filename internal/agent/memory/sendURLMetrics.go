package memory

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
func (m *MemoryStats) SendURLmetrics(cfg *config.ConfigAgent) {
    
    for {
        time.Sleep(cfg.PauseDuration)
        for i := range m.GaugeMetrics {
            err := SendURLGauge(cfg.URL, m.GaugeMetrics, i)
            if err != nil {
                logger.Log.Info("unexpected error:", zap.Error(err))
            }
        }
        err := sendURLCounter(cfg.URL,  m.CounterMetric)
        if err != nil {
            logger.Log.Info("unexpected error:", zap.Error(err))
        }
        logger.Log.Info("All URL metrics sent successfully")
    }
}

func sendURLCounter(url string, counter int) error {
    agent := resty.New()
    req := agent.R().SetPathParams(map[string]string{
        "metricType":  config.CountType,
        "metricName":  config.PollCount,
        "metricValue": strconv.Itoa(counter),
    }).SetHeader("Content-Type", "text/plain")
    _, err := req.Post(url + "/update/{metricType}/{metricName}/{metricValue}")
    if err != nil {
        switch  {
        case os.IsTimeout(err):
            for _, delay := range []time.Duration{time.Second, 2 * time.Second, 3 * time.Second} {
                time.Sleep(delay)
                if _, err = req.Post(url + urlParams); err == nil {
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
    return nil
}

func SendURLGauge(url string, gaugeMetrics map[string]float64, metricName string) error {
    agent := resty.New()
    req := agent.R().SetPathParams(map[string]string{
        "metricType":  config.GaugeType,
        "metricName":  metricName,
        "metricValue": strconv.FormatFloat(gaugeMetrics[metricName], 'f', 6, 64),
    }).SetHeader("Content-Type", "text/plain")
    _, err := req.Post(url + "/update/{metricType}/{metricName}/{metricValue}")
    if err != nil {
        switch {
        case os.IsTimeout(err):
            for _, delay := range []time.Duration{time.Second, 2 * time.Second, 3 * time.Second} {
                time.Sleep(delay)
                if _, err = req.Post(url + urlParams); err == nil {
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
    return nil
}

