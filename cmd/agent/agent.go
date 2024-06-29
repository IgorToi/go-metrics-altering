package main

import (
	"fmt"
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
			req := agent.R().SetBody(models.Metrics{
				ID:	i,
				MType: agentConfig.GaugeType,
				Value: &v,
			}).SetHeader("Content-Type", "application/json")
			req.URL = agentConfig.ProtocolScheme + cfg.FlagRunAddr
			_, err := req.Post(req.URL + "/update/")
			if err != nil {
				logger.Log.Debug("unexpected sending metric error", zap.Error(err))
			}
			cfg.Count++
			logger.Log.Info("metric sent")	
		}
		req := agent.R().SetBody(models.Metrics{
			ID: agentConfig.PollCount,
			MType: agentConfig.CountType,
			Delta: &cfg.Count,
		}).SetHeader("Content-Type", "application/json")
		req.URL = agentConfig.ProtocolScheme + cfg.FlagRunAddr
		_, err := req.Post(req.URL + "/update/")
		fmt.Println(req.Body)
		if err != nil {
			logger.Log.Debug("unexpected sending metric error", zap.Error(err))
		}
		logger.Log.Info("metric sent")	
	}   
}

