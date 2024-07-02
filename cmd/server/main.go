package main

import (
	"fmt"
	"log"
	"net/http"

	serverConfig "github.com/IgorToi/go-metrics-altering/internal/config/server_config"
	httpServer "github.com/IgorToi/go-metrics-altering/internal/server"
)

func main() {	
	cfg, err := serverConfig.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Running server on", cfg.FlagRunAddr)
	http.ListenAndServe(cfg.FlagRunAddr, httpServer.MetricRouter(cfg))
}

