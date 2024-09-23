package cash

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	config "github.com/igortoigildin/go-metrics-altering/config/agent"
	"github.com/igortoigildin/go-metrics-altering/internal/models"
	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
	"go.uber.org/zap"
)

const pollCount = "PollCount"

type MemStorage struct {
	rm      sync.RWMutex
	Gauge   map[string]float64
	Counter map[string]int64
}

func InitLocalStorage() *MemStorage {
	var m MemStorage
	m.Counter = make(map[string]int64)
	m.Counter[pollCount] = 0
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

// iterate through memStorage
func (m *MemStorage) GetAllMetrics() []models.Metrics {
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

// load metrics from local file
func (m *MemStorage) LoadAllFromFile(fname string) error {
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

// SaveMetrics periodically saves metrics from memStorage.
func (m *MemStorage) SaveAllMetrics(FlagStoreInterval int, FlagStorePath string, fname string) error {
	pauseDuration := time.Duration(FlagStoreInterval) * time.Second
	for {
		time.Sleep(pauseDuration)
		metricSlice := m.GetAllMetrics()
		data, err := json.MarshalIndent(metricSlice, "", "  ")
		if err != nil {
			logger.Log.Info("marshalling error", zap.Error(err))
			return err
		}

		err = os.WriteFile(fname, data, 0606)
		if err != nil {
			logger.Log.Info("error saving metrics to the file", zap.Error(err))
			return err
		}
	}
}
