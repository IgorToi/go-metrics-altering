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

type LocalStorage interface {
	UpdateGaugeMetric(metricName string, metricValue float64)
	UpdateCounterMetric(metricName string, metricValue int64)
	GetGaugeMetricFromMemory(metricName string) float64
	GetCountMetricFromMemory(metricName string) int64
	CheckIfGaugeMetricPresent(metricName string) bool
	CheckIfCountMetricPresent(metricName string) bool
	GetAllMetrics() []models.Metrics
	LoadAllFromFile(fname string) error
	SaveAllMetrics(FlagStoreInterval int, FlagStorePath string, fname string) error
}

func UpdateHandler(LocalStorage LocalStorage) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			logger.Log.Debug("got request with bad method", zap.String("method", r.Method))
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var req models.Metrics
		err := processjson.ReadJSON(r, &req)
		if err != nil {
			logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
			processjson.SendJSONError(w, http.StatusBadRequest, "badly-formed JSON")
			return
		}

		if req.MType != config.GaugeType && req.MType != config.CountType {
			logger.Log.Info("usupported request type", zap.String("type", req.MType))
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		switch req.MType {
		case config.GaugeType:
			LocalStorage.UpdateGaugeMetric(req.ID, *req.Value)
		case config.CountType:
			LocalStorage.UpdateCounterMetric(config.PollCount, *req.Delta)
		}

		delta := LocalStorage.GetCountMetricFromMemory(config.PollCount)
		req.Delta = &delta
		resp := models.Metrics{
			ID:    req.ID,
			MType: req.MType,
			Value: req.Value,
			Delta: req.Delta,
		}

		err = processjson.WriteJSON(w, http.StatusOK, resp, nil)
		if err != nil {
			logger.Log.Debug("error encoding response", zap.Error(err))
			return
		}
	})
}

func ValueHandler(LocalStorage LocalStorage) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			logger.Log.Info("got request with bad method", zap.String("method", r.Method))
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var req models.Metrics
		err := processjson.ReadJSON(r, req)
		if err != nil {
			logger.Log.Info("cannot decode request JSON body", zap.Error(err))
			processjson.SendJSONError(w, http.StatusBadRequest, "badly-formed JSON")
			return
		}

		resp := models.Metrics{
			ID:    req.ID,
			MType: req.MType,
		}
		switch req.MType {
		case config.GaugeType:
			if !LocalStorage.CheckIfGaugeMetricPresent(req.ID) {
				logger.Log.Info("usupported metric name", zap.String("name", req.ID))
				w.WriteHeader(http.StatusNotFound)
				return
			}
			value := LocalStorage.GetGaugeMetricFromMemory(req.ID)
			resp.Value = &value
		case config.CountType:
			delta := LocalStorage.GetCountMetricFromMemory(config.PollCount)
			resp.Delta = &delta
		default:
			logger.Log.Info("usupported request type", zap.String("type", req.MType))
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		w.Header().Add("Content-Encoding", "gzip")
		err = processjson.WriteJSON(w, http.StatusOK, resp, nil)
		if err != nil {
			logger.Log.Info("error encoding response", zap.Error(err))
			return
		}
	})
}

func UpdateHandle(LocalStorage LocalStorage) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")
		if metricName == "" {
			logger.Log.Info("metricName not provided")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		metricValue := chi.URLParam(r, "metricValue")
		switch metricType {
		case config.GaugeType:
			metricValueConverted, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				logger.Log.Info("error parsing metric value to float", zap.Error(err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			LocalStorage.UpdateGaugeMetric(metricName, metricValueConverted)
		case config.CountType:
			metricValueConverted, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				logger.Log.Info("error parsing metric value to int", zap.Error(err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			LocalStorage.UpdateCounterMetric(metricName, metricValueConverted)
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	})
}

// v1
func ValueHandle(LocalStorage LocalStorage) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")

		switch metricType {
		case config.GaugeType:
			if !LocalStorage.CheckIfGaugeMetricPresent(metricName) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			_, err := w.Write([]byte(strconv.FormatFloat(LocalStorage.GetGaugeMetricFromMemory(metricName), 'f', -1, 64)))
			if err != nil {
				logger.Log.Info("error occured while sending response", zap.Error(err))
			}
		case config.CountType:
			if !LocalStorage.CheckIfCountMetricPresent(metricName) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			_, err := w.Write([]byte(strconv.FormatInt(LocalStorage.GetCountMetricFromMemory(metricName), 10)))
			if err != nil {
				logger.Log.Debug("error occured while sending response", zap.Error(err))
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
	})
}
