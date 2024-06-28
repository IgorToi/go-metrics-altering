package server

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	agentConfig "github.com/IgorToi/go-metrics-altering/internal/config/agent_config"
	"github.com/IgorToi/go-metrics-altering/internal/logger"
	"github.com/IgorToi/go-metrics-altering/internal/models"
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
	
	r.Post("/update/",	WithLogging(http.HandlerFunc(memory.UpdateHandler)))
	r.Get("/value/{metricType}/{metricName}", WithLogging(http.HandlerFunc(memory.ValueHandle)))
	r.Post("/value/", WithLogging(http.HandlerFunc(memory.ValueHandler)))
	r.Get("/value/{metricType}/{metricName}", WithLogging(http.HandlerFunc(memory.ValueHandle)))
	r.Get("/", memory.InformationHandle)
	return r
}

func (m *MemStorage) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Log.Debug("got request with bad method", zap.String("method", r.Method))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	logger.Log.Debug("decoding request")
	var req models.Metrics
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if req.MType != agentConfig.GaugeType && req.MType != agentConfig.CountType {
		logger.Log.Debug("usupported request type", zap.String("type", req.MType))
		w.WriteHeader(http.StatusUnprocessableEntity)
        return
	}
	resp := models.Metrics{
		ID: req.ID,
		MType: req.MType,
		Delta: req.Delta,
		Value: req.Value,
	}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		return
	}
	logger.Log.Debug("sending HTTP 200 response")
}

func (m *MemStorage) ValueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Log.Debug("got request with bad method", zap.String("method", r.Method))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	logger.Log.Debug("decoding request")
	var req models.Metrics
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp := models.Metrics{
		ID:		req.ID,
		MType: 	req.MType,
	}
	switch req.MType {
	case agentConfig.GaugeType:
		if !m.CheckIfGaugeMetricPresent(req.ID) {
			logger.Log.Debug("usupported metric name", zap.String("name", req.ID))
			w.WriteHeader(http.StatusNotFound)
			return
		}
		value := m.GetGaugeMetricFromMemory(req.ID)
		resp.Value = &value
	case agentConfig.CountType:
		if !m.CheckIfCountMetricPresent(req.ID) {
			logger.Log.Debug("usupported metric name", zap.String("name", req.ID))
			w.WriteHeader(http.StatusNotFound)
			return
		}
		delta :=  m.GetCountMetricFromMemory(req.ID)
		resp.Delta = &delta
	default:
		logger.Log.Debug("usupported request type", zap.String("type", req.MType))
		w.WriteHeader(http.StatusUnprocessableEntity)
        return
	}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		return
	}
	logger.Log.Debug("sending HTTP 200 response")
}

//without body
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
//without body
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
