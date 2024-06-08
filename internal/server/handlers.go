package server

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	agentConfig "github.com/IgorToi/go-metrics-altering/internal/config/agent_config"
	"github.com/IgorToi/go-metrics-altering/templates"
	"github.com/go-chi/chi"
)

var t *template.Template

func ParseTemplate() *template.Template {
	t, err := template.ParseFS(templates.FS, "home.gohtml")
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func MetricRouter() chi.Router {
	memory := InitStorage()
	t = ParseTemplate()
	r := chi.NewRouter()
	r.Post("/update/{metricType}/{metricName}/{metricValue}", memory.UpdateHandle)
	r.Get("/value/{metricType}/{metricName}", memory.ValueHandle)
	r.Get("/", memory.InformationHandle)
	return r
}

func (m *MemStorage) UpdateHandle(rw http.ResponseWriter, r *http.Request) { 
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	if metricName ==  "" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	metricValue := chi.URLParam(r, "metricValue")
	switch metricType {
	case agentConfig.GaugeType:
		metricValueConverted, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		m.UpdateGaugeMetric(metricName, metricValueConverted)
	case agentConfig.CountType:
		metricValueConverted, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		m.UpdateCounterMetric(metricName, metricValueConverted)
	default:
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

func (m *MemStorage) ValueHandle(rw http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	switch metricType {
	case agentConfig.GaugeType:
		if !m.CheckIfGaugeMetricPresent(metricName) {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		rw.Write([]byte(strconv.FormatFloat(m.GetGaugeMetricFromMemory(metricName),'f', -1, 64)))
	case agentConfig.CountType:
		if !m.CheckIfCountMetricPresent(metricName) {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		rw.Write([]byte(strconv.FormatInt(m.GetCountMetricFromMemory(metricName), 10)))
	default:
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
}

func (m *MemStorage) InformationHandle(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	
    if err := t.Execute(rw, ConvertToSingleMap(m.Gauge, m.Counter)); err != nil {
		log.Printf("Failed to execute template: %v", err)
	}
}