package main

import (
	"log"

	config "github.com/IgorToi/go-metrics-altering/internal/config/server_config"
	server "github.com/IgorToi/go-metrics-altering/internal/server"
)

func main() {	
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err := server.Run(cfg); err != nil {
		log.Fatal(err)
	}
}



