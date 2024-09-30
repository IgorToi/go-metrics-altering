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
	rm       sync.RWMutex
	Gauge    map[string]float64
	Counter  map[string]int64
	strategy Strategy
}

func InitLocalStorage() *LocalStorage {
	var m LocalStorage
	m.Counter = make(map[string]int64)
	m.Counter[pollCount] = 0
	m.Gauge = make(map[string]float64)
	return &m
}

func (m *LocalStorage) SetStrategy(metricType string) {
	if metricType == config.CountType {
		count := count{}
		m.strategy = &count
	} else {
		gauge := gauge{}
		m.strategy = &gauge
	}
}

func (m *LocalStorage) Update(ctx context.Context, metricType string, metricName string, metricValue any) error {
	m.SetStrategy(metricType)
	return m.strategy.Update(metricType, metricName, metricValue)
}

func (m *LocalStorage) Get(ctx context.Context, metricType string, metricName string) (models.Metrics, error) {
	m.SetStrategy(metricType)
	return m.strategy.Get(metricType, metricName)
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
