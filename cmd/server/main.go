package main

import (
	"log"

	config "github.com/igortoigildin/go-metrics-altering/config/server"
	"github.com/igortoigildin/go-metrics-altering/internal/server/app"
	"go.uber.org/zap"

	"github.com/igortoigildin/go-metrics-altering/pkg/crypt"
	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Println("error while logading config", err)
		return
	}

	if err = logger.Initialize(cfg.FlagLogLevel); err != nil {
		log.Println("error while initializing logger", err)
		return
	}

	err = crypt.InitRSAKeys(cfg)
	if err != nil {
		logger.Log.Error("error while generating rsa keys", zap.Error(err))
		return
	}

	app.Run(cfg)
}
