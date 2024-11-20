package main

import (
	"log"

	config "github.com/igortoigildin/go-metrics-altering/config/server"
	grpcapp "github.com/igortoigildin/go-metrics-altering/internal/server/grpc/app/grpc"
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

	grpcapp.New()

	go grpcapp.Run(cfg, storage)

	httpapp.Run(cfg, storage)
}
