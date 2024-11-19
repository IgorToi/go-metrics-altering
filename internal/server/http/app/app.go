package httpapp

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	config "github.com/igortoigildin/go-metrics-altering/config/server"
	server "github.com/igortoigildin/go-metrics-altering/internal/server/http/api"
	storage "github.com/igortoigildin/go-metrics-altering/internal/storage"
	httpServer "github.com/igortoigildin/go-metrics-altering/pkg/httpServer"
	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

var buildVersion string = "N/A"
var buildDate string = "N/A"
var buildCommit string = "N/A"

func Run(cfg *config.ConfigServer, storage storage.Storage) {
	printInfo()

	ctx := context.Background()

	// Handlers
	r := server.Router(ctx, cfg, storage)

	logger.Log.Info("Starting server on", zap.String("address", cfg.FlagRunAddr))

	// HTTP server
	srv, _ := httpServer.Address(cfg.FlagRunAddr)
	httpSrv := httpServer.New(r, srv)

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

	logger.Log.Info("Graceful server shutdown complete...")
}

func printInfo() error {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	return nil
}
