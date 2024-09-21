package main

import (
	"context"
	"net/http"

	config "github.com/igortoigildin/go-metrics-altering/config/server"
	db "github.com/igortoigildin/go-metrics-altering/internal/server/api/saveMetricsToDB"
	mem "github.com/igortoigildin/go-metrics-altering/internal/server/api/saveMetricsToMemory"
	"github.com/igortoigildin/go-metrics-altering/internal/storage"
	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Log.Fatal("error while logading config", zap.Error(err))
	}

	if err := logger.Initialize(cfg.FlagLogLevel); err != nil {
		logger.Log.Fatal("error while initializing logger", zap.Error(err))
	}

	logger.Log.Info("Running server", zap.String("address", cfg.FlagRunAddr))

	switch cfg.FlagDBDSN {
	// start server with option of saving metrics in memory
	case "":
		http.ListenAndServe(cfg.FlagRunAddr, mem.MetricRouter(cfg, ctx))
		if err != nil {
			logger.Log.Fatal("cannot start the server", zap.Error(err))
		}
	default:
		// start server with option of saving metrics in db
		repo := storage.InitPostgresRepo(ctx, cfg)

		http.ListenAndServe(cfg.FlagRunAddr, db.RouterDB(ctx, cfg, repo))
		if err != nil {
			logger.Log.Fatal("cannot start the server", zap.Error(err))
		}
	}
}
