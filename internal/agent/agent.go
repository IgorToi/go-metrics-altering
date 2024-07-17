package agent

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	config "github.com/igortoigildin/go-metrics-altering/config/agent"
	"github.com/igortoigildin/go-metrics-altering/internal/logger"
	"github.com/igortoigildin/go-metrics-altering/internal/models"
	"go.uber.org/zap"
)

func RunAgent() {
    cfg, err := config.LoadConfig()
    if err != nil {
        logger.Log.Fatal("error while logading config", zap.Error(err))
    }
    if err := logger.Initialize(cfg.FlagLogLevel); err != nil {
        logger.Log.Fatal("error while initializing logger", zap.Error(err))
    }
    // updating metrics in memory every pollInterval
    go cfg.UpdateMetrics()
    // v2 - metrics in body json
    go sendMetricsInJSON(cfg)
    // v1 - metrics in url path
    sendMetricsViaURL(cfg)
}

func sendMetricsViaURL(cfg *config.ConfigAgent) {
// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
    agent := resty.New()
    for {
        time.Sleep(cfg.PauseDuration)
        for i := range cfg.Memory {
            sendGaugeURLPath(cfg, agent, i)
        }
        sendCounterURLPath(cfg, agent)
        logger.Log.Info("Metric has been sent successfully")
    }
}

func sendCounterURLPath(cfg *config.ConfigAgent, agent *resty.Client) {
    _, err := agent.R().SetPathParams(map[string]string{
        "metricType":  config.CountType,
        "metricName":  config.PollCount,
        "metricValue": strconv.Itoa(cfg.Count),
    }).SetHeader("Content-Type", "text/plain").Post(cfg.URL + "/update/{metricType}/{metricName}/{metricValue}")
    if err != nil {
        logger.Log.Debug("unexpected sending metric error:", zap.Error(err))
    }
}

func sendGaugeURLPath(cfg *config.ConfigAgent, agent *resty.Client, metricName string) {
    _, err := agent.R().SetPathParams(map[string]string{
        "metricType":  config.GaugeType,
        "metricName": metricName,
        "metricValue": strconv.FormatFloat(cfg.Memory[metricName], 'f', 6, 64),
    }).SetHeader("Content-Type", "text/plain").Post(cfg.URL + "/update/{metricType}/{metricName}/{metricValue}")
    if err != nil {
        logger.Log.Debug("unexpected sending metric error:", zap.Error(err))
    }
}

func sendMetricsInJSON(cfg *config.ConfigAgent) {
    agent := resty.New()
    for {
        time.Sleep(cfg.PauseDuration)
        for i := range cfg.Memory {
            SendGaugeInJSON(cfg, i, agent)
        }
        SendCounterInJSON(cfg, agent)
    }
}

func SendGaugeInJSON(cfg *config.ConfigAgent, metricName string, agent *resty.Client) {
    valueGauge := cfg.Memory[metricName]
    metric := models.Metrics{
        ID:    metricName,
        MType: config.GaugeType,
        Value: &valueGauge,
    }
    req := agent.R().SetHeader("Content-Type", "application/json")
    metricsJSON, err := json.Marshal(metric)
    if err != nil {
        logger.Log.Debug("unexpected sending metric error:", zap.Error(err))
    }
    req.URL = cfg.URL + "/update/"
    _, err = req.SetBody(metricsJSON).Post(req.URL)
    if err != nil {
        //send again n times if timeout error
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
}

func SendCounterInJSON(cfg *config.ConfigAgent, agent *resty.Client) {
    valueDelta := int64(cfg.Count)
    metric := models.Metrics{
        ID:    config.PollCount,
        MType: config.CountType,
        Delta: &valueDelta,
    }
    req := agent.R().SetHeader("Content-Type", "application/json")
    metricsJSON, err := json.Marshal(metric)
    if err != nil {
        logger.Log.Debug("unexpected sending metric error:", zap.Error(err))
    }
    req.URL = cfg.URL + "/update/"
    _, err = req.SetBody(metricsJSON).Post(req.URL)
    if err != nil {
        //send again n times if timeout error
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
}