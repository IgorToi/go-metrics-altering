package main

import (
	"encoding/json"
	"fmt"
	"os"
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
		logger.Log.Fatal("error while logading config", zap.Error(err))
	}


	if err := logger.Initialize(cfg.FlagLogLevel); err != nil {
		logger.Log.Fatal("error while initializing logger", zap.Error(err))
	}



	// start goroutine to update metrics every pollInterval
	go cfg.UpdateMetrics()

	// start goroutine to send bunch metrics 
	go SendAllMetrics(cfg)



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
				// if error due to timeout - try send again
				fmt.Println("ERROR sending metric", err)	
				if os.IsTimeout(err) {
						for n, t := 1, 1; n <= 3; n++ {
						time.Sleep(time.Duration(t) * time.Second)
						if _, err := httpAgent.SendMetric(req.URL, config.GaugeType, i, strconv.FormatFloat(cfg.Memory[i], 'f', 6, 64), req); err == nil {
							break
						}
						t += 2
					}
				}
				logger.Log.Debug("unexpected sending metric error:", zap.Error(err))
			}
			// logger.Log.Info("Metric has been sent successfully")

			// preparing and sending slice of metrics to /updates/
			metrics := PrepareMetricBody(cfg, i)
			metricsJSON, err := json.Marshal(metrics)
			if err != nil {
				logger.Log.Debug("unexpected sending metric error:", zap.Error(err))
			}
			_, err = req.SetBody(metricsJSON).SetHeader("Content-Type", "application/json").Post(req.URL + "/updates/")
			if err != nil {
				// if error due to timeout - try send again
				fmt.Println("ERROR sending metric", err)
				if os.IsTimeout(err) {
					for n, t := 1, 1; n <= 3; n++ {
						time.Sleep(time.Duration(t) * time.Second)
						metrics := PrepareMetricBody(cfg, i)
						metricsJSON, err := json.Marshal(metrics)
						if err != nil {
							logger.Log.Debug("unexpected sending metric error:", zap.Error(err))
						}
						if _, err = req.SetBody(metricsJSON).SetHeader("Content-Type", "application/json").Post(req.URL + "/updates/"); err == nil {
							break
						}
						t += 2
					}
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
			// if error due to timeout - try send again
			fmt.Println("ERROR sending metric", err)
			if os.IsTimeout(err) {
				for n, t := 1, 1; n <= 3; n++ {
					time.Sleep(time.Duration(t) * time.Second)
					if _, err := httpAgent.SendMetric(req.URL, config.CountType, config.PollCount, strconv.Itoa(cfg.Count), req); err == nil {
						break
					}
					t += 2
				}
			}
			logger.Log.Debug("unexpected sending metric error:", zap.Error(err))
		}
		logger.Log.Info("Metric has been sent successfully")
	}
}

func PrepareMetricBody(cfg *config.ConfigAgent, metricName string) []models.Metrics {
	var metrics []models.Metrics
	valueGauge := cfg.Memory[metricName]
	metric := models.Metrics{
		ID:    metricName,
		MType: config.GaugeType,
		Value: &valueGauge,
	}
	metrics = append(metrics, metric)

	valueDelta := int64(cfg.Count)
	metric = models.Metrics{
		ID:    config.PollCount,
		MType: config.CountType,
		Delta: &valueDelta,
	}
	metrics = append(metrics, metric)
	return metrics
}



func SendAllMetrics(cfg *config.ConfigAgent) () {
	var metrics []models.Metrics
	agent := resty.New()
	req := agent.R().SetHeader("Content-Type", "application/json")
	req.URL = config.ProtocolScheme + cfg.FlagRunAddr
	durationPause := time.Duration(cfg.FlagReportInterval) * time.Second
	for {
		metrics = metrics[:0]
		time.Sleep(durationPause)
		for i := range cfg.Memory {
			metric := PrepareMetricBodyNew(cfg, i)
			metrics = append(metrics, metric)
		}
		fmt.Println(metrics)
		metricsJSON, err := json.Marshal(metrics)
		if err != nil {
			logger.Log.Debug("unexpected sending metric error:", zap.Error(err))
		}
		req.URL = req.URL + "/updates/"
		_, err = req.SetBody(metricsJSON).Post(req.URL)
		if err != nil {
			fmt.Println("NEW ERROR")

				// if error due to timeout - try send again

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
}


 

 func PrepareMetricBodyNew(cfg *config.ConfigAgent, metricName string) models.Metrics {
	var metric models.Metrics
	switch metricName {
	case config.PollCount:
		valueDelta := int64(cfg.Count)
		metric = models.Metrics{
			ID:    	config.PollCount,
			MType: 	config.CountType,
			Delta: 	&valueDelta,
		}
	default:
		valueGauge := cfg.Memory[metricName]
		metric = models.Metrics{
			ID:    metricName,
			MType: config.GaugeType,
			Value: &valueGauge,
		}
	}
	return metric
}