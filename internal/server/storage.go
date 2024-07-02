package server

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	config "github.com/IgorToi/go-metrics-altering/internal/config/server_config"
	"github.com/IgorToi/go-metrics-altering/internal/logger"
	"github.com/IgorToi/go-metrics-altering/internal/models"
	"go.uber.org/zap"
)

// convert memStorage to slice of models.Metrics
func (memory MemStorage) ConvertToSlice() []models.Metrics {
	var metricSlice []models.Metrics
	var model models.Metrics


	for i, v := range  memory.Gauge {
		model.ID = i
		model.MType = "gauge"
		model.Value = &v
		metricSlice = append(metricSlice, model)
	}


	for j, u := range memory.Counter {
		model.ID = j
		model.MType = "counter"
		model.Delta = &u
		metricSlice = append(metricSlice, model)
	}
	return metricSlice
}

func Save(fname string, metricSlice []models.Metrics)  error {
	data, err := json.MarshalIndent(metricSlice, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fname, data, 0606)
}

func (m *MemStorage) Load(fname string) error {
	data, err := os.ReadFile(fname)
	if err != nil {
		return err
	}

	var metricSlice []models.Metrics
	if err := json.Unmarshal(data, &metricSlice); err != nil {
		return err
	}

	for _, v := range metricSlice {
		if v.MType == "gauge" {
			m.Gauge[v.ID] = *v.Value
		} else if v.MType == "counter" {
			m.Counter[v.ID] = *v.Delta
		}
	}
	return nil
}

func (memory *MemStorage) saveMetrics(cfg *config.ConfigServer) {

	interval, err  := strconv.Atoi(cfg.FlagStoreInterval)
	if err != nil {
		fmt.Println(err)
		logger.Log.Debug("cannot decode time interval", zap.Error(err))
	}
	pauseDuration := time.Duration(interval) * time.Second
	for {
		time.Sleep(pauseDuration)
		metricSlice := memory.ConvertToSlice()
		Save(cfg.FlagStorePath, metricSlice)
	}

}
