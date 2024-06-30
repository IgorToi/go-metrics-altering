package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"log"
	"time"

	agentConfig "github.com/IgorToi/go-metrics-altering/internal/config/agent_config"
	"github.com/IgorToi/go-metrics-altering/internal/logger"
	"github.com/IgorToi/go-metrics-altering/internal/models"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func main() {
	cfg, err := agentConfig.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
    // start goroutine to update metrics every pollInterval
	go cfg.UpdateMetrics()  
	agent := resty.New()

	
	durationPause := time.Duration(cfg.FlagReportInterval) * time.Second
	for {
		time.Sleep(durationPause)
		for i, v := range cfg.Memory {

			requestBody := models.Metrics{
				ID: i,
				MType: agentConfig.GaugeType,
				Value: &v,
			}

			var b bytes.Buffer
			gw := gzip.NewWriter(&b)
			err := json.NewEncoder(gw).Encode(requestBody)
			if err != nil {
				logger.Log.Info("unexpected encoding body error",zap.Error(err))
			}
			gw.Close()

			req := agent.R().SetBody(&b).SetHeader("Content-Type", "application/json")
			req.SetHeader("Content-Encoding", "gzip")
			req.URL = agentConfig.ProtocolScheme + cfg.FlagRunAddr
			_, err = req.Post(req.URL + "/update/")

			if err != nil {
				logger.Log.Info("unexpected sending metric error", zap.Error(err))
			}
			cfg.Count++
			logger.Log.Info("metric sent")	
		}

		requestBody := models.Metrics{
			ID: agentConfig.PollCount,
			MType: agentConfig.CountType,
			Delta: &cfg.Count,
		}
		var b bytes.Buffer
		gw := gzip.NewWriter(&b)
		err := json.NewEncoder(gw).Encode(requestBody)
		if err != nil {
			logger.Log.Info("unexpected encoding body error",zap.Error(err))
		}
		gw.Close()
		req := agent.R().SetBody(&b).SetHeader("Content-Type", "application/json")
		req.SetHeader("Content-Encoding", "gzip")


		req.URL = agentConfig.ProtocolScheme + cfg.FlagRunAddr
		_, err = req.Post(req.URL + "/update/")

		if err != nil {
			logger.Log.Info("unexpected sending metric error", zap.Error(err))
		}
		logger.Log.Info("metric sent")	
	}   
}

