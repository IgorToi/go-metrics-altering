package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	agentHandlers "github.com/IgorToi/go-metrics-altering/internal/agent"
	agentConfig "github.com/IgorToi/go-metrics-altering/internal/config/agent_config"
	"github.com/go-resty/resty/v2"
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
			req := agent.R()
			req.SetPathParams(map[string]string{
				"metricType": agentConfig.GaugeType,
				"metricName": i,
				"metricValue": strconv.FormatFloat(v, 'f', 6, 64 ),
			}).SetHeader("Content-Type", "text/plain")
			req.URL = agentConfig.ProtocolScheme + cfg.FlagRunAddr
			_, err := agentHandlers.SendMetric(req.URL, agentConfig.GaugeType, i,  strconv.FormatFloat(v, 'f', 6, 64 ), req)
			if err != nil {
				fmt.Printf("unexpected sending metric error: %s", err)
			}
			fmt.Println("Metric has been sent successfully")
		}
		req := agent.R()
		req.SetPathParams(map[string]string{
			"metricType": agentConfig.CountType,
			"metricName": agentConfig.PollCount,
			"metricValue": strconv.Itoa(cfg.Count),
		}).SetHeader("Content-Type", "text/plain")
		req.URL = agentConfig.ProtocolScheme + cfg.FlagRunAddr
		_, err := agentHandlers.SendMetric(req.URL, agentConfig.CountType, agentConfig.PollCount,  strconv.Itoa(cfg.Count), req)
		if err != nil {
			fmt.Printf("unexpected sending metric error: %s", err)
		}
		fmt.Println("Metric has been sent successfully")
	}   
}

