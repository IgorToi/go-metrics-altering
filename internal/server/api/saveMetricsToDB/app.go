package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	_ "net/http/pprof" // подключаем пакет pprof

	config "github.com/igortoigildin/go-metrics-altering/config/server"
	"github.com/igortoigildin/go-metrics-altering/internal/models"
	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
	"go.uber.org/zap"
)

type Storage interface {
	Update(ctx context.Context, metricType string, metricName string, metricValue any) error
	Get(ctx context.Context, metricType string, metricName string) (models.Metrics, error)
	GetAll(ctx context.Context) (map[string]any, error)
	Ping(ctx context.Context) error
}

type app struct {
	storage Storage
	cfg     *config.ConfigServer
}

func newApp(s Storage, cfg *config.ConfigServer) *app {
	return &app{
		storage: s,
		cfg:     cfg,
	}
}

func (app *app) Ping(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := app.storage.Ping(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (app *app) updates(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	if r.Method != http.MethodPost {
		logger.Log.Debug("got request with bad method", zap.String("method", r.Method))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var metrics []models.Metrics
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&metrics); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for _, metric := range metrics { // iterating through []Metrics and adding it to db one by one
		if metric.MType != config.GaugeType && metric.MType != config.CountType {
			logger.Log.Debug("usupported request type", zap.String("type", metric.MType))
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		switch metric.MType {
		case config.GaugeType:
			err := app.storage.Update(ctx, metric.MType, metric.ID, metric.Value)
			if err != nil {
				logger.Log.Debug("error while updating value", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		case config.CountType:
			err := app.storage.Update(ctx, metric.MType, metric.ID, metric.Delta)
			if err != nil {
				logger.Log.Debug("error while updating value", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
	w.Header().Set("Content-Type", "application/json")
	logger.Log.Debug("sending HTTP 200 response")
}

func (app *app) updateMetric(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	if r.Method != http.MethodPost {
		logger.Log.Debug("got request with bad method", zap.String("method", r.Method))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req models.Metrics
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Log.Info("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if req.MType != config.GaugeType && req.MType != config.CountType {
		logger.Log.Debug("usupported request type", zap.String("type", req.MType))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	switch req.MType {
	case config.GaugeType:
		err := app.storage.Update(ctx, req.MType, req.ID, req.Value)
		if err != nil {
			logger.Log.Debug("error while updating value", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case config.CountType:
		err := app.storage.Update(ctx, req.MType, req.ID, req.Delta)
		if err != nil {
			logger.Log.Debug("error while updating value", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	resp := models.Metrics{
		ID:    req.ID,
		MType: req.MType,
		Value: req.Value,
		Delta: req.Delta,
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	logger.Log.Debug("sending HTTP 200 response")
}

func (app *app) getAllmetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Add("Content-Encoding", "gzip")
	metrics, err := app.storage.GetAll(r.Context())
	if err != nil {
		logger.Log.Debug("error", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, metrics); err != nil {
		logger.Log.Debug("error executing template", zap.Error(err))
	}
}

func (app *app) getMetric(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
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
		ID:    req.ID,
		MType: req.MType,
	}
	switch req.MType {
	case config.GaugeType:
		res, err := app.storage.Get(ctx, req.MType, req.ID)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				logger.Log.Debug("error while obtaining metric", zap.Error(err))
				w.WriteHeader(http.StatusNotFound)
				return
			default:
				logger.Log.Debug("error while obtaining metric", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		resp.Value = res.Value
	case config.CountType:
		res, err := app.storage.Get(ctx, req.MType, req.ID)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				logger.Log.Debug("error while obtaining metric", zap.Error(err))
				w.WriteHeader(http.StatusNotFound)
				return
			default:
				logger.Log.Debug("error while obtaining metric", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		resp.Delta = res.Delta
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Content-Encoding", "gzip")
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		return
	}
	logger.Log.Debug("sending HTTP 200 response")
}
