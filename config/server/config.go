package config

import (
	"errors"
	"flag"
	"html/template"
	"os"
	"strconv"
)

const (
	GaugeType = "gauge"
	CountType = "counter"
	PollCount = "PollCount"
)

type ConfigServer struct {
	FlagRunAddr       string
	Template          *template.Template
	FlagLogLevel      string
	FlagStoreInterval int
	FlagStorePath     string
	FlagRestore       bool
}

func LoadConfig() (*ConfigServer, error) {
	cfg := new(ConfigServer)
	var err error
	flag.StringVar(&cfg.FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&cfg.FlagLogLevel, "l", "info", "log level")
	flag.IntVar(&cfg.FlagStoreInterval, "i", 1, "metrics backup interval")
	flag.StringVar(&cfg.FlagStorePath, "f", "/tmp/metrics-db.json", "metrics backup storage path")
	flag.BoolVar(&cfg.FlagRestore, "r", true, "true if load from backup is needed")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.FlagRunAddr = envRunAddr
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		cfg.FlagLogLevel = envLogLevel
	}
	if envStorageInterval := os.Getenv("STORE_INTERVAL"); envStorageInterval != "" {
		// parse string env variable
		v, err := strconv.Atoi(envStorageInterval)
		if err != nil {
			return nil, err
		}
		cfg.FlagStoreInterval = v
	}
	if envStorePath := os.Getenv("FILE_STORAGE_PATH"); envStorePath != "" {
		cfg.FlagStorePath = envStorePath
	}
	if envFlagRestore := os.Getenv("RESTORE"); envFlagRestore != "" {
		// parse bool env variable
		v, err := strconv.ParseBool(envFlagRestore)
		if err != nil {
			return nil, err
		}
		cfg.FlagRestore = v
	}
	// check if any config variables is empty
	var cfgVarEmpty = errors.New("configs variable not set")
	if !cfg.validate() {
		return nil, cfgVarEmpty
	}
	return cfg, err
}

func (cfg *ConfigServer) validate() bool {
	if cfg.FlagRunAddr == "" {
		return false
	}
	if cfg.FlagLogLevel == "" {
		return false
	}
	if cfg.FlagStorePath == "" {
		return false
	}
	return true
}
