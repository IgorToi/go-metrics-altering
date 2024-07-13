package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	config "github.com/igortoigildin/go-metrics-altering/config/server"
	"github.com/igortoigildin/go-metrics-altering/internal/logger"
	"github.com/igortoigildin/go-metrics-altering/internal/models"
	"github.com/igortoigildin/go-metrics-altering/internal/storage"
	"github.com/igortoigildin/go-metrics-altering/templates"
	"go.uber.org/zap"
)

type app struct {
	storage storage.Storage
}

func newApp(s storage.Storage) *app {
	return &app{storage: s}
}

func routerDB(ctx context.Context, cfg *config.ConfigServer) chi.Router {
	repo := storage.InitPostgresRepo(ctx, cfg)
	app := newApp(repo)
	t = templates.ParseTemplate()
	r := chi.NewRouter()
	r.Get("/ping", WithLogging(gzipMiddleware(http.HandlerFunc(repo.Ping))))
	r.Get("/", WithLogging(gzipMiddleware(http.HandlerFunc(app.getAllmetrics))))
	r.Post("/update/", WithLogging(gzipMiddleware(http.HandlerFunc(app.updateMetric))))
	r.Post("/value/", WithLogging(gzipMiddleware(http.HandlerFunc(app.getMetric))))
	return r
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

func (app *app) updateMetric(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.MType != config.GaugeType && req.MType != config.CountType {
		logger.Log.Debug("usupported request type", zap.String("type", req.MType))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	
	switch req.MType {
	case config.GaugeType:
	if app.storage.Exist(ctx, req.MType, req.ID) {
		err := app.storage.Update(ctx, req.MType, req.ID, req.Value)
		if err != nil {
			logger.Log.Debug("error while updating value", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		err := app.storage.Add(ctx, req.MType, req.ID, req.Value)
		if err != nil {
			logger.Log.Debug("error while adding value", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	case config.CountType:
		if app.storage.Exist(ctx, req.MType, req.ID) {
			err := app.storage.Update(ctx, req.MType, req.ID, req.Value)
			if err != nil {
				logger.Log.Debug("error while updating value", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			err := app.storage.Add(ctx, req.MType, req.ID, req.Value)
			if err != nil {
				logger.Log.Debug("error while adding value", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
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
		if app.storage.Exist(ctx, req.MType, req.ID) {
			res, err := app.storage.Get(ctx, req.MType, req.ID)
			if err != nil {
				logger.Log.Debug("error while obtaining metric", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			resp.Value = res.Value
		} else {
			logger.Log.Debug("usupported metric name", zap.String("name", req.ID))
			w.WriteHeader(http.StatusNotFound)
			return
		}
	case config.CountType:
		if app.storage.Exist(ctx, req.MType, req.ID) {
			res, err := app.storage.Get(ctx, req.MType, req.ID)
			if err != nil {
				logger.Log.Debug("error while obtaining metric", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			resp.Delta = res.Delta
		} else {
			logger.Log.Debug("usupported metric name", zap.String("name", req.ID))
			w.WriteHeader(http.StatusNotFound)
			return
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
}