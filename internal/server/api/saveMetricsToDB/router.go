package api

import (
	"context"
	"net/http"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	config "github.com/igortoigildin/go-metrics-altering/config/server"
	auth "github.com/igortoigildin/go-metrics-altering/pkg/middlewares/auth"
	compress "github.com/igortoigildin/go-metrics-altering/pkg/middlewares/compress"
	logging "github.com/igortoigildin/go-metrics-altering/pkg/middlewares/logging"
	timeout "github.com/igortoigildin/go-metrics-altering/pkg/middlewares/timeout"
	"github.com/igortoigildin/go-metrics-altering/templates"
)

var t *template.Template

func RouterDB(ctx context.Context, cfg *config.ConfigServer, s Storage) chi.Router {
	// repo := storage.InitPostgresRepo(ctx, cfg)
	// app := newApp(repo, cfg)
	t = templates.ParseTemplate()
	r := chi.NewRouter()

	r.Get("/ping", timeout.Timeout(cfg.ContextTimout, logging.WithLogging(compress.GzipMiddleware(http.HandlerFunc(ping(s))))))
	r.Get("/", timeout.Timeout(cfg.ContextTimout, logging.WithLogging(compress.GzipMiddleware(auth.Auth(http.HandlerFunc(getAllmetrics(s)), cfg)))))
	r.Post("/update/", timeout.Timeout(cfg.ContextTimout, logging.WithLogging(compress.GzipMiddleware(auth.Auth(http.HandlerFunc(updateMetric(s)), cfg)))))
	r.Post("/updates/", timeout.Timeout(cfg.ContextTimout, logging.WithLogging(compress.GzipMiddleware(auth.Auth(http.HandlerFunc(updates(s)), cfg)))))
	r.Post("/value/", timeout.Timeout(cfg.ContextTimout, logging.WithLogging(compress.GzipMiddleware(auth.Auth(http.HandlerFunc(getMetric(s)), cfg)))))
	r.Mount("/debug", middleware.Profiler())

	return r
}
