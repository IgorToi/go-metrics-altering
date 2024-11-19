package api

import (
	"context"
	"net/http"
	"text/template"

	config "github.com/igortoigildin/go-metrics-altering/config/server"
	"github.com/igortoigildin/go-metrics-altering/pkg/middlewares/auth"
	"github.com/igortoigildin/go-metrics-altering/pkg/middlewares/compress"
	"github.com/igortoigildin/go-metrics-altering/pkg/middlewares/logging"
	"github.com/igortoigildin/go-metrics-altering/pkg/middlewares/timeout"
	"github.com/igortoigildin/go-metrics-altering/templates"
)

var t *template.Template

func Router(ctx context.Context, cfg *config.ConfigServer, storage Storage) *http.ServeMux {
	t = templates.ParseTemplate()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /value/{metricType}/{metricName}", logging.WithLogging(compress.GzipMiddleware((auth.Auth(http.HandlerFunc(valuePathHandler(storage)), cfg)))))
	mux.HandleFunc("POST /update/{metricType}/{metricName}/{metricValue}", logging.WithLogging(compress.GzipMiddleware(auth.Auth(http.HandlerFunc(updatePathHandler(storage)), cfg))))
	mux.HandleFunc("GET /ping", timeout.Timeout(cfg.ContextTimout, logging.WithLogging(compress.GzipMiddleware(auth.Auth(http.HandlerFunc(ping(storage)), cfg)))))
	mux.HandleFunc("GET /", timeout.Timeout(cfg.ContextTimout, logging.WithLogging(compress.GzipMiddleware(auth.Auth(http.HandlerFunc(getAllmetrics(storage)), cfg)))))
	mux.HandleFunc("POST /updates/", timeout.Timeout(cfg.ContextTimout, logging.WithLogging(compress.GzipMiddleware(auth.Auth(http.HandlerFunc(updates(storage)), cfg)))))
	mux.HandleFunc("POST /value/", timeout.Timeout(cfg.ContextTimout, logging.WithLogging(compress.GzipMiddleware(auth.Auth(http.HandlerFunc(getMetric(storage)), cfg)))))
	mux.HandleFunc("POST /update/", timeout.Timeout(cfg.ContextTimout, logging.WithLogging(compress.GzipMiddleware(auth.Auth(http.HandlerFunc(updateMetric(storage)), cfg)))))

	return mux
} 
