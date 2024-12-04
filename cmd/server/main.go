package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	config "github.com/igortoigildin/go-metrics-altering/config/server"
	grpcapp "github.com/igortoigildin/go-metrics-altering/internal/server/grpc/app"
	server "github.com/igortoigildin/go-metrics-altering/internal/server/http/api"
	"github.com/igortoigildin/go-metrics-altering/internal/storage"
	httpServer "github.com/igortoigildin/go-metrics-altering/pkg/httpServer"
	"go.uber.org/zap"

	"github.com/igortoigildin/go-metrics-altering/pkg/crypt"
	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("error while logading config", err)
	}

	if err = logger.Initialize(cfg.FlagLogLevel); err != nil {
		log.Fatal("error while initializing logger", err)
	}

	err = crypt.InitRSAKeys(cfg)
	if err != nil {
		logger.Log.Error("error while generating rsa keys", zap.Error(err))
	}

	storage, err := storage.New(cfg)
	if err != nil {
		logger.Log.Fatal("failed to init storage", zap.Error(err))
	}

	// gRPC
	application := grpcapp.New(cfg, storage)

	go func() {
		err := application.GRPCServer.MustRun()
		if err != nil {
			logger.Log.Fatal("error while initializing grpc app", zap.Error(err))
		}
	}()

	// http
	r := server.Router(context.Background(), cfg, storage)

	logger.Log.Info("Starting server on", zap.String("address", cfg.FlagRunAddrHTTP))

	// HTTP server
	srv, _ := httpServer.Address(cfg.FlagRunAddrHTTP)
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
	err = httpSrv.GracefulShutdown()
	if err != nil {
		logger.Log.Error("error:", zap.Error(err))
	}

	logger.Log.Info("Graceful server shutdown complete...")
}
