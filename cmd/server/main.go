package main

import (
	"log"

	config "github.com/igortoigildin/go-metrics-altering/config/server"
	grpcapp "github.com/igortoigildin/go-metrics-altering/internal/server/grpc/app"
	httpapp "github.com/igortoigildin/go-metrics-altering/internal/server/http/app"
	"github.com/igortoigildin/go-metrics-altering/internal/storage"
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

	storage := storage.New(cfg)

	application := grpcapp.New(cfg, storage)

	go func() {
		err := application.GRPCServer.MustRun()
		if err != nil {
			logger.Log.Fatal("error while initializing grpc app", zap.Error(err))
		}
	}()

	httpapp.Run(cfg, storage)
}
