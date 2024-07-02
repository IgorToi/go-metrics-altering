package server

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	config "github.com/IgorToi/go-metrics-altering/internal/config/server_config"
	"github.com/IgorToi/go-metrics-altering/internal/logger"
	"github.com/IgorToi/go-metrics-altering/internal/models"
	"go.uber.org/zap"
)

// iterate through memStorage
func (m MemStorage) ConvertToSlice() []models.Metrics {
    var metricSlice []models.Metrics
    var model models.Metrics
    for i, v := range  m.Gauge {
        model.ID = i
        model.MType = "gauge"
        model.Value = &v
        metricSlice = append(metricSlice, model)
    }
    for j, u := range m.Counter {
        model.ID = j
        model.MType = "counter"
        model.Delta = &u
        metricSlice = append(metricSlice, model)
    }
    return metricSlice
}
// save slice with metrics to the file
func Save(fname string, metricSlice []models.Metrics)  error {
    data, err := json.MarshalIndent(metricSlice, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(fname, data, 0606)
}
// load metrics from local file
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
// save metrics from memStorage to the file every StoreInterval
func (m *MemStorage) saveMetrics(cfg *config.ConfigServer) {
    interval, err  := strconv.Atoi(cfg.FlagStoreInterval)
    if err != nil {
        logger.Log.Debug("cannot decode time interval", zap.Error(err))
    }
    pauseDuration := time.Duration(interval) * time.Second
    for {
        time.Sleep(pauseDuration)
        metricSlice := m.ConvertToSlice()
        Save(cfg.FlagStorePath, metricSlice)
    }

}

