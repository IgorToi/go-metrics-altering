package config

import (
	"flag"
	"html/template"
	"os"
)

type ConfigServer struct {
	FlagRunAddr 	string
	Template 		*template.Template
}

func LoadConfig() (*ConfigServer, error) {
	cfg := new(ConfigServer)
	var err error 
	flag.StringVar(&cfg.FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
        cfg.FlagRunAddr = envRunAddr
	}
	return cfg, err
}


