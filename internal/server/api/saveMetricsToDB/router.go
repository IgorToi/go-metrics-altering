package api

import (
	"context"
	"net/http"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	config "github.com/igortoigildin/go-metrics-altering/config/server"
	"github.com/igortoigildin/go-metrics-altering/internal/storage"
	auth "github.com/igortoigildin/go-metrics-altering/pkg/middlewares/auth"
	compress "github.com/igortoigildin/go-metrics-altering/pkg/middlewares/compress"
	logging "github.com/igortoigildin/go-metrics-altering/pkg/middlewares/logging"
	"github.com/igortoigildin/go-metrics-altering/templates"
)

var t *template.Template

func RouterDB(ctx context.Context, cfg *config.ConfigServer) chi.Router {
	repo := storage.InitPostgresRepo(ctx, cfg)
	app := newApp(repo, cfg)
	t = templates.ParseTemplate()
	r := chi.NewRouter()
	r.Get("/ping", logging.WithLogging(compress.GzipMiddleware(http.HandlerFunc(app.Ping))))
	r.Get("/", logging.WithLogging(compress.GzipMiddleware(auth.Auth(http.HandlerFunc(app.getAllmetrics), cfg))))
	r.Post("/update/", logging.WithLogging(compress.GzipMiddleware(auth.Auth(http.HandlerFunc(app.updateMetric), cfg))))
	r.Post("/updates/", logging.WithLogging(compress.GzipMiddleware(auth.Auth(http.HandlerFunc(app.updates), cfg))))
	r.Post("/value/", logging.WithLogging(compress.GzipMiddleware(auth.Auth(http.HandlerFunc(app.getMetric), cfg))))
	r.Mount("/debug", middleware.Profiler())

	return r
}
