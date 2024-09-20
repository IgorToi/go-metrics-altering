package api

import (
	"sync"
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
