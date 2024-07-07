package main

import (
	"net/http"

	config "github.com/igortoigildin/go-metrics-altering/config/server"
	"github.com/igortoigildin/go-metrics-altering/internal/logger"
	"github.com/igortoigildin/go-metrics-altering/internal/server"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Log.Fatal("error while logading config", zap.Error(err))
	}
	logger.Log.Info("Running server on", zap.String("port", cfg.FlagRunAddr))
	http.ListenAndServe(cfg.FlagRunAddr, server.MetricRouter(cfg))
}
