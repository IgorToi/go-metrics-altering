package main

import (
	"context"
	"net/http"

	config "github.com/igortoigildin/go-metrics-altering/config/server"
	server "github.com/igortoigildin/go-metrics-altering/internal/server/api"
	local "github.com/igortoigildin/go-metrics-altering/internal/storage/local"
	psql "github.com/igortoigildin/go-metrics-altering/internal/storage/postgres"
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

	DBstorage := psql.InitPostgresRepo(ctx, cfg)
	localStorage := local.InitLocalStorage()

	logger.Log.Info("Running server", zap.String("address", cfg.FlagRunAddr))

	http.ListenAndServe(cfg.FlagRunAddr, server.Router(ctx, cfg, DBstorage, localStorage))
	if err != nil {
		logger.Log.Fatal("cannot start the server", zap.Error(err))
	}
}
