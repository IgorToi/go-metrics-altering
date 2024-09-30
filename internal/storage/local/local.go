package local

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	config "github.com/igortoigildin/go-metrics-altering/config/agent"
	"github.com/igortoigildin/go-metrics-altering/internal/models"
	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
	processmap "github.com/igortoigildin/go-metrics-altering/pkg/processMap"
	"go.uber.org/zap"
)

const pollCount = "PollCount"

type LocalStorage struct {
	rm      sync.RWMutex
	Gauge   map[string]float64
	Counter map[string]int64
}

func InitLocalStorage() *LocalStorage {
	var m LocalStorage
	m.Counter = make(map[string]int64)
	m.Counter[pollCount] = 0
	m.Gauge = make(map[string]float64)
	return &m
}

func (m *LocalStorage) Update(ctx context.Context, metricType string, metricName string, metricValue any) error {
	switch metricType {
	case config.GaugeType:
		if m.Gauge == nil {
			m.Gauge = make(map[string]float64)
		}
		m.rm.Lock()
		v, _ := metricValue.(*float64)
		var temp float64
		if v == nil {
			temp = float64(0)
		} else {
			temp = *v
		}
		m.Gauge[metricName] = temp
		m.rm.Unlock()
	case config.CountType:
		if m.Counter == nil {
			m.Counter = make(map[string]int64)
		}
		m.rm.Lock()
		v, _ := metricValue.(*int64)
		var temp int64
		if v == nil {
			temp = int64(0)
		} else {
			temp = *v
		}
		m.Counter[metricName] += temp
		m.rm.Unlock()
	default:
		return errors.New("undefined metric type")
	}
	return nil
}

func (m *LocalStorage) Get(ctx context.Context, metricType string, metricName string) (models.Metrics, error) {
	var metric models.Metrics

	switch metricType {
	case config.GaugeType:
		m.rm.RLock()
		v := m.Gauge[metricName]
		metric.Value = &v
		m.rm.RUnlock()
	case config.CountType:
		m.rm.RLock()
		d := m.Counter[metricName]
		metric.Delta = &d
		m.rm.RUnlock()
	default:
		return metric, errors.New("undefined metric type")
	}
	metric.MType = metricType

	return metric, nil
}

func (m *LocalStorage) GetAll(ctx context.Context) (map[string]any, error) {
	return processmap.ConvertToSingleMap(m.Gauge, m.Counter), nil
}

// LoadMetricsFromFile loads metrics from the stated file.
func (m *LocalStorage) LoadMetricsFromFile(fname string) error {
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

func (m *LocalStorage) Ping(ctx context.Context) error {
	if m.Gauge == nil {
		logger.Log.Info("gauge local storage not initialized")
		return errors.New("gauge local storage not initialized")
	}

	if m.Counter == nil {
		logger.Log.Info("counter local storage not initialized")
		return errors.New("counter local storage not initialized")
	}
	return nil
}

// SaveMetrics periodically saves metrics from local storage to provided file.
func (m *LocalStorage) SaveAllMetricsToFile(FlagStoreInterval int, FlagStorePath string, fname string) error {
	pauseDuration := time.Duration(FlagStoreInterval) * time.Second
	for {
		time.Sleep(pauseDuration)
		metrics, err := m.GetAll(context.Background())
		if err != nil {
			return err
		}

		data, err := json.MarshalIndent(metrics, "", "  ")
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
