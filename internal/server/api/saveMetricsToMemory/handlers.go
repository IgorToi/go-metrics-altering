package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	config "github.com/igortoigildin/go-metrics-altering/config/server"
	"github.com/igortoigildin/go-metrics-altering/internal/models"
	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
	processjson "github.com/igortoigildin/go-metrics-altering/pkg/processJSON"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

func (rep *DB) Ping(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := rep.conn.PingContext(ctx); err != nil {
		logger.Log.Info("error", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rep.conn.Close()
	w.WriteHeader(http.StatusOK)
}

// v2 with requst body
func (m *MemStorage) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Log.Debug("got request with bad method", zap.String("method", r.Method))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.Metrics
	err := processjson.ReadJSON(r, req)
	if err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		processjson.SendJSONError(w, http.StatusBadRequest, "badly-formed JSON")
		return
	}

	if req.MType != config.GaugeType && req.MType != config.CountType {
		logger.Log.Debug("usupported request type", zap.String("type", req.MType))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	switch req.MType {
	case config.GaugeType:
		m.UpdateGaugeMetric(req.ID, *req.Value)
	case config.CountType:
		m.Counter[config.PollCount] += *req.Delta
	}

	var resp models.Metrics
	delta := m.Counter[config.PollCount]
	if delta == 0 {
		resp = models.Metrics{
			ID:    req.ID,
			MType: req.MType,
			Value: req.Value,
		}
	} else {
		resp = models.Metrics{
			ID:    req.ID,
			MType: req.MType,
			Value: req.Value,
			Delta: &delta,
		}
	}

	err = processjson.WriteJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		return
	}
}

// v2 with request body
func (m *MemStorage) ValueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Log.Debug("got request with bad method", zap.String("method", r.Method))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.Metrics
	err := processjson.ReadJSON(r, req)
	if err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		processjson.SendJSONError(w, http.StatusBadRequest, "badly-formed JSON")
		return
	}

	resp := models.Metrics{
		ID:    req.ID,
		MType: req.MType,
	}
	switch req.MType {
	case config.GaugeType:
		if !m.CheckIfGaugeMetricPresent(req.ID) {
			logger.Log.Info("usupported metric name", zap.String("name", req.ID))
			w.WriteHeader(http.StatusNotFound)
			return
		}
		value := m.GetGaugeMetricFromMemory(req.ID)
		resp.Value = &value
	case config.CountType:
		delta := m.Counter[config.PollCount]
		resp.Delta = &delta
	default:
		logger.Log.Debug("usupported request type", zap.String("type", req.MType))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	w.Header().Add("Content-Encoding", "gzip")
	err = processjson.WriteJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		return
	}
}

// v1
func (m *MemStorage) UpdateHandle(rw http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	if metricName == "" {
		logger.Log.Info("metricName not provided")
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	metricValue := chi.URLParam(r, "metricValue")
	switch metricType {
	case config.GaugeType:
		metricValueConverted, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			logger.Log.Debug("error parsing metric value to float", zap.Error(err))
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		m.UpdateGaugeMetric(metricName, metricValueConverted)
	case config.CountType:
		metricValueConverted, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			logger.Log.Debug("error parsing metric value to int", zap.Error(err))
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

// v1
func (m *MemStorage) ValueHandle(rw http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	switch metricType {
	case config.GaugeType:
		if !m.CheckIfGaugeMetricPresent(metricName) {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		_, err := rw.Write([]byte(strconv.FormatFloat(m.GetGaugeMetricFromMemory(metricName), 'f', -1, 64)))
		if err != nil {
			logger.Log.Debug("error occured while sending response", zap.Error(err))
		}
	case config.CountType:
		if !m.CheckIfCountMetricPresent(metricName) {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		_, err := rw.Write([]byte(strconv.FormatInt(m.GetCountMetricFromMemory(metricName), 10)))
		if err != nil {
			logger.Log.Debug("error occured while sending response", zap.Error(err))
		}
	default:
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
}

func (m *MemStorage) InformationHandle(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	rw.Header().Add("Content-Encoding", "gzip")
	if err := t.Execute(rw, ConvertToSingleMap(m.Gauge, m.Counter)); err != nil {
		logger.Log.Debug("error executing template", zap.Error(err))
	}
}
