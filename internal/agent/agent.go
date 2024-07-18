package agent

import (
	"errors"

	config "github.com/igortoigildin/go-metrics-altering/config/agent"
	"github.com/igortoigildin/go-metrics-altering/internal/logger"
	"go.uber.org/zap"
)

var (
    ErrConnectionFailed = errors.New("connection failed")
)

func RunAgent() {
    cfg, err := config.LoadConfig()
    if err != nil {
        logger.Log.Fatal("error while logading config", zap.Error(err))
    }
    if err := logger.Initialize(cfg.FlagLogLevel); err != nil {
        logger.Log.Fatal("error while initializing logger", zap.Error(err))
    }
    go cfg.UpdateMetrics()      // updating metrics in memory every pollInterval
    sendJSONMetrics(cfg)        // v2 - metrics in body json
    go sendBatchMetrics(cfg)    // v3 - sending batchs of metrics json
    err = sendURLmetrics(cfg)   // v1 - metrics in url path
    if err != nil {
        logger.Log.Debug("error while sending URLpath metric", zap.Error(err))
    }
}

