package server

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	agentConfig "github.com/IgorToi/go-metrics-altering/internal/config/agent_config"
	"github.com/IgorToi/go-metrics-altering/internal/logger"
	"github.com/IgorToi/go-metrics-altering/templates"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

var t *template.Template

type (
	// info struct about reply
		responseData	struct {
		status 		int
		size		int
	}
	loggingResponseWriter	struct {
		http.ResponseWriter
		responseData	*responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// WithLogging adds code to regester info regarding request and returns new http.Handler
func WithLogging(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		responseData := &responseData{
			status: 0,
			size: 0,
		}
		lw := loggingResponseWriter {
			ResponseWriter: w,	//using original http.ResponseWriter
			responseData: responseData,
		}
		h.ServeHTTP(&lw, r) //using updated http.ResponseWriter
		duration := time.Since(start)
		logger.Log.Info("got incoming HTTP request",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Int("status", responseData.status),
			zap.String("duration", duration.String()),
			zap.Int("size", responseData.size),
		)
	})
}	

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
	
	r.Post("/update/{metricType}/{metricName}/{metricValue}",	WithLogging(http.HandlerFunc(memory.UpdateHandle)))
	
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
