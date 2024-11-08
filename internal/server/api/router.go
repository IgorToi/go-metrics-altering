package api

import (
	"context"
	"net/http"
	"text/template"

	config "github.com/igortoigildin/go-metrics-altering/config/server"
	"github.com/igortoigildin/go-metrics-altering/pkg/middlewares/compress"
	"github.com/igortoigildin/go-metrics-altering/pkg/middlewares/logging"
	"github.com/igortoigildin/go-metrics-altering/pkg/middlewares/timeout"
	"github.com/igortoigildin/go-metrics-altering/templates"
)

var t *template.Template

func Router(ctx context.Context, cfg *config.ConfigServer, storage Storage) *http.ServeMux {
	t = templates.ParseTemplate()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /value/{metricType}/{metricName}", logging.WithLogging(compress.GzipMiddleware((http.HandlerFunc(valuePathHandler(storage))))))
	mux.HandleFunc("POST /update/{metricType}/{metricName}/{metricValue}", logging.WithLogging(compress.GzipMiddleware((http.HandlerFunc(updatePathHandler(storage))))))
	mux.HandleFunc("GET /ping", timeout.Timeout(cfg.ContextTimout, logging.WithLogging(compress.GzipMiddleware(http.HandlerFunc(ping(storage))))))
	mux.HandleFunc("GET /", timeout.Timeout(cfg.ContextTimout, logging.WithLogging(compress.GzipMiddleware((http.HandlerFunc(getAllmetrics(storage)))))))
	mux.HandleFunc("POST /updates/", timeout.Timeout(cfg.ContextTimout, logging.WithLogging(compress.GzipMiddleware((http.HandlerFunc(updates(storage)))))))
	mux.HandleFunc("POST /value/", timeout.Timeout(cfg.ContextTimout, logging.WithLogging(compress.GzipMiddleware((http.HandlerFunc(getMetric(storage)))))))
	mux.HandleFunc("POST /update/", timeout.Timeout(cfg.ContextTimout, logging.WithLogging(compress.GzipMiddleware((http.HandlerFunc(updateMetric(storage)))))))

	return mux
}
