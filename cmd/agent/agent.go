package main

import (
	"log"

	config "github.com/igortoigildin/go-metrics-altering/config/agent"
	"github.com/igortoigildin/go-metrics-altering/internal/agent/app"
	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Println("error while logading config", err)
		return
	}
	if err := logger.Initialize(cfg.FlagLogLevel); err != nil {
		log.Println("error while initializing logger", err)
		return
	}

	app.Run(cfg)
}
