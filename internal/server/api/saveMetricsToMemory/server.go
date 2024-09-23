package api

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	config "github.com/igortoigildin/go-metrics-altering/config/server"
	"github.com/igortoigildin/go-metrics-altering/internal/models"
	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
	"go.uber.org/zap"
)

type MemStorage struct {
	rm      sync.RWMutex
	Gauge   map[string]float64
	Counter map[string]int64
}

func InitStorage() *MemStorage {
	var m MemStorage
	m.Counter = make(map[string]int64)
	m.Counter["PollCount"] = 0
	m.Gauge = make(map[string]float64)
	return &m
}

func (m *MemStorage) UpdateGaugeMetric(metricName string, metricValue float64) {
	if m.Gauge == nil {
		m.Gauge = make(map[string]float64)
	}
	m.rm.Lock()
	m.Gauge[metricName] = metricValue
	m.rm.Unlock()
}

func (m *MemStorage) UpdateCounterMetric(metricName string, metricValue int64) {
	if m.Counter == nil {
		m.Counter = make(map[string]int64)
	}
	m.rm.Lock()
	m.Counter[metricName] += metricValue
	m.rm.Unlock()
}

func (m *MemStorage) GetGaugeMetricFromMemory(metricName string) float64 {
	m.rm.RLock()
	metric := m.Gauge[metricName]
	m.rm.RUnlock()
	return metric
}

func (m *MemStorage) GetCountMetricFromMemory(metricName string) int64 {
	m.rm.RLock()
	metric := m.Counter[metricName]
	m.rm.RUnlock()
	return metric
}

func (m *MemStorage) CheckIfGaugeMetricPresent(metricName string) bool {
	m.rm.RLock()
	_, ok := m.Gauge[metricName]
	m.rm.RUnlock()
	return ok
}

func (m *MemStorage) CheckIfCountMetricPresent(metricName string) bool {
	m.rm.RLock()
	_, ok := m.Counter[metricName]
	m.rm.RUnlock()
	return ok
}

func ConvertToSingleMap(a map[string]float64, b map[string]int64) map[string]interface{} {
	c := make(map[string]interface{}, 33)
	for i, v := range a {
		c[i] = v
	}
	for j, l := range b {
		c[j] = l
	}
	return c
}

// iterate through memStorage
func (m *MemStorage) ConvertToSlice() []models.Metrics {
	metricSlice := make([]models.Metrics, 33)
	var model models.Metrics
	for i, v := range m.Gauge {
		model.ID = i
		model.MType = config.GaugeType

		model.Value = &v
		metricSlice = append(metricSlice, model)
	}
	for j, u := range m.Counter {
		model.ID = j
		model.MType = config.CountType
		model.Delta = &u
		metricSlice = append(metricSlice, model)
	}
	return metricSlice
}

// save slice with metrics to the file
func Save(fname string, metricSlice []models.Metrics) error {
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
	pauseDuration := time.Duration(cfg.FlagStoreInterval) * time.Second
	for {
		time.Sleep(pauseDuration)
		metricSlice := m.ConvertToSlice()
		if err := Save(cfg.FlagStorePath, metricSlice); err != nil {
			logger.Log.Info("error saving metrics to the file", zap.Error(err))
		}
	}
}
