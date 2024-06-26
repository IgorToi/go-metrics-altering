package config

import (
	"flag"
	"html/template"
	"os"
)

type ConfigServer struct {
	FlagRunAddr 	string
	Template 		*template.Template
	FlagLogLevel	string
}

func LoadConfig() (*ConfigServer, error) {
	cfg := new(ConfigServer)
	var err error 
	flag.StringVar(&cfg.FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&cfg.FlagLogLevel, "l", "info", "log level")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
        cfg.FlagRunAddr = envRunAddr
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		cfg.FlagLogLevel = envLogLevel
	}
	return cfg, err
}


