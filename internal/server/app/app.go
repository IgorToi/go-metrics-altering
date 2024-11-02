package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	config "github.com/igortoigildin/go-metrics-altering/config/server"
	server "github.com/igortoigildin/go-metrics-altering/internal/server/api"
	storage "github.com/igortoigildin/go-metrics-altering/internal/server/api"
	httpServer "github.com/igortoigildin/go-metrics-altering/pkg/httpServer"
	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
	"go.uber.org/zap"
)

var buildVersion string = "N/A"
var buildDate string = "N/A"
var buildCommit string = "N/A"


func Run(cfg *config.ConfigServer) {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	ctx := context.Background()
	
	// Postgres Storage initialization
	storage := storage.New(cfg)

	// Handlers
	r := server.Router(ctx, cfg, storage)

	logger.Log.Info("Starting server on", zap.String("address", cfg.FlagRunAddr))

	// HTTP server
	httpSrv := httpServer.New(r, httpServer.Port(cfg.FlagRunAddr))

	// waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	
	select {
	case s := <-interrupt:
		logger.Log.Info("Received: ", zap.String("signal", s.String()))
	case err := <-httpSrv.Notify():
		logger.Log.Error("error:", zap.Error(err))
	}

	// graceful shutdown
	err := httpSrv.GracefulShutdown()
	if err != nil {
		logger.Log.Error("error:", zap.Error(err))
	}
}