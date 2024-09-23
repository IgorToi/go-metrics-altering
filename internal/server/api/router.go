package api

import (
	"context"
	"net/http"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	config "github.com/igortoigildin/go-metrics-altering/config/server"
	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
	auth "github.com/igortoigildin/go-metrics-altering/pkg/middlewares/auth"
	compress "github.com/igortoigildin/go-metrics-altering/pkg/middlewares/compress"
	logging "github.com/igortoigildin/go-metrics-altering/pkg/middlewares/logging"
	timeout "github.com/igortoigildin/go-metrics-altering/pkg/middlewares/timeout"
	"github.com/igortoigildin/go-metrics-altering/templates"
	"go.uber.org/zap"
)

var t *template.Template

func Router(ctx context.Context, cfg *config.ConfigServer, s Storage, m LocalStorage) chi.Router {
	t = templates.ParseTemplate()
	r := chi.NewRouter()

	// if flag restore is true - load metrics from the local file
	if cfg.FlagRestore {
		err := m.LoadAllFromFile(cfg.FlagStorePath)
		if err != nil {
			logger.Log.Info("error loading metrics from the file", zap.Error(err))
		}
	}

	// start goroutine to save metrics locally
	go m.SaveAllMetrics(cfg.FlagStoreInterval, cfg.FlagStorePath, cfg.FlagStorePath)

	// v1 handlers to recieve metrics via URL
	r.Get("/value/{metricType}/{metricName}", logging.WithLogging(compress.GzipMiddleware(auth.Auth(http.HandlerFunc(ValueHandle(m)), cfg))))
	r.Post("/update/{metricType}/{metricName}/{metricValue}", logging.WithLogging(compress.GzipMiddleware(auth.Auth(http.HandlerFunc(UpdateHandle(m)), cfg))))

	// v2 handlers to receive metrics in JSON
	r.Get("/ping", timeout.Timeout(cfg.ContextTimout, logging.WithLogging(compress.GzipMiddleware(http.HandlerFunc(ping(s))))))
	r.Get("/", timeout.Timeout(cfg.ContextTimout, logging.WithLogging(compress.GzipMiddleware(auth.Auth(http.HandlerFunc(getAllmetrics(s)), cfg)))))
	r.Post("/updates/", timeout.Timeout(cfg.ContextTimout, logging.WithLogging(compress.GzipMiddleware(auth.Auth(http.HandlerFunc(updates(s)), cfg)))))

	// Check whether metrics should be saved to DB or locally.
	switch cfg.FlagDBDSN {
	case "":
		r.Post("/update/", logging.WithLogging(compress.GzipMiddleware(auth.Auth(http.HandlerFunc(UpdateHandler(m)), cfg))))
		r.Post("/value/", logging.WithLogging(compress.GzipMiddleware(auth.Auth(http.HandlerFunc(ValueHandler(m)), cfg))))
	default:
		r.Post("/value/", timeout.Timeout(cfg.ContextTimout, logging.WithLogging(compress.GzipMiddleware(auth.Auth(http.HandlerFunc(getMetric(s)), cfg)))))
		r.Post("/update/", timeout.Timeout(cfg.ContextTimout, logging.WithLogging(compress.GzipMiddleware(auth.Auth(http.HandlerFunc(updateMetric(s)), cfg)))))
	}

	r.Mount("/debug", middleware.Profiler())

	return r
}
