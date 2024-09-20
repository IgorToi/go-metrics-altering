package api

import (
	"context"
	"net/http"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	config "github.com/igortoigildin/go-metrics-altering/config/server"
	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
	"github.com/igortoigildin/go-metrics-altering/pkg/middlewares/auth"
	"github.com/igortoigildin/go-metrics-altering/pkg/middlewares/compress"
	"github.com/igortoigildin/go-metrics-altering/pkg/middlewares/logging"
	"github.com/igortoigildin/go-metrics-altering/templates"
	"go.uber.org/zap"
)

var t *template.Template

func MetricRouter(cfg *config.ConfigServer, ctx context.Context) chi.Router {
	var m = InitStorage()
	if cfg.FlagRestore {
		err := m.Load(cfg.FlagStorePath) // if flag restore is true - load metrics from the local file
		if err != nil {
			logger.Log.Info("error loading metrics from the file", zap.Error(err))
		}
	}
	// start goroutine to save metrics locally
	go m.saveMetrics(cfg)
	// parse template
	t = templates.ParseTemplate()
	r := chi.NewRouter()

	r.Mount("/debug", middleware.Profiler())

	DB := InitRepo(ctx, cfg)
	r.Get("/ping", logging.WithLogging(compress.GzipMiddleware(http.HandlerFunc(DB.Ping))))
	// v1
	r.Get("/value/{metricType}/{metricName}", logging.WithLogging(compress.GzipMiddleware(auth.Auth(http.HandlerFunc(m.ValueHandle), cfg))))
	r.Get("/", logging.WithLogging(compress.GzipMiddleware(http.HandlerFunc(auth.Auth(m.InformationHandle, cfg)))))

	r.Route("/", func(r chi.Router) {
		// v1
		r.Post("/update/{metricType}/{metricName}/{metricValue}", logging.WithLogging(compress.GzipMiddleware(auth.Auth(http.HandlerFunc(m.UpdateHandle), cfg))))
		// v2
		r.Post("/update/", logging.WithLogging(compress.GzipMiddleware(auth.Auth(http.HandlerFunc(m.UpdateHandler), cfg))))
		r.Post("/value/", logging.WithLogging(compress.GzipMiddleware(auth.Auth(http.HandlerFunc(m.ValueHandler), cfg))))
	})
	return r
}
