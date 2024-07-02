package server

import (
	"net/http"
	"strings"

	config "github.com/IgorToi/go-metrics-altering/internal/config/server_config"
	"github.com/IgorToi/go-metrics-altering/internal/logger"
	"go.uber.org/zap"
)

type MemStorage struct {
	Gauge     map[string]float64
	Counter   map[string]int64
}

func Run(cfg *config.ConfigServer) error {
	if err := logger.Initialize(cfg.FlagLogLevel); err != nil {
		return err
	}
	logger.Log.Info("Running server", zap.String("address", cfg.FlagRunAddr))
	return http.ListenAndServe(cfg.FlagRunAddr, MetricRouter(cfg))
}

func InitStorage() (*MemStorage) {
	var m MemStorage
	m.Counter  = make(map[string]int64)
	m.Counter["PollCount"] = 0
	m.Gauge = make(map[string]float64)
	return &m
}

func (m *MemStorage) UpdateGaugeMetric(metricName string, metricValue float64) {
	if m.Gauge == nil {
		m.Gauge = make(map[string]float64)
	}
	m.Gauge[metricName] = metricValue
}

func (m *MemStorage) UpdateCounterMetric(metricName string, metricValue int64) {
	if m.Counter == nil {
		m.Counter = make(map[string]int64)
	}
	m.Counter[metricName] += metricValue
	
}

func (m *MemStorage) GetGaugeMetricFromMemory(metricName string) float64 {
	return m.Gauge[metricName]
}

func (m *MemStorage) GetCountMetricFromMemory(metricName string) int64 {
	return m.Counter[metricName]
}

func (m *MemStorage) CheckIfGaugeMetricPresent(metricName string) bool {
	_, ok := m.Gauge[metricName]
	return ok
}

func (m *MemStorage) CheckIfCountMetricPresent(metricName string) bool {
	_, ok := m.Counter[metricName]
	return ok
}

func ConvertToSingleMap(a map[string]float64, b map[string]int64) map[string]interface{} {
	c := make(map[string]interface{})
	for i, v := range a {
		c[i] = v
	}
	for j, l := range b {
		c[j] = l
	}
	return c
}

func gzipMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ow := w
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
                return
			}
			r.Body = cr
			defer cr.Close()
		}
		h.ServeHTTP(ow, r)
	}
}

