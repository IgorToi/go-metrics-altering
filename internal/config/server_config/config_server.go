package config

import (
	"flag"
	"html/template"
	"os"
)

type ConfigServer struct {
    FlagRunAddr         string
    Template            *template.Template
    FlagLogLevel        string
    FlagStoreInterval   string
    FlagStorePath       string
    FlagRestore         string
}

func LoadConfig() (*ConfigServer, error) {
    cfg := new(ConfigServer)
    var err error 
    flag.StringVar(&cfg.FlagRunAddr, "a", ":8080", "address and port to run server")
    flag.StringVar(&cfg.FlagLogLevel, "l", "info", "log level")
    flag.StringVar(&cfg.FlagStoreInterval, "i", "1", "metrics backup interval")
    flag.StringVar(&cfg.FlagStorePath, "f", "/tmp/metrics-db.json", "metrics backup storage path")
    flag.StringVar(&cfg.FlagRestore, "r", "true", "true if load from backup is needed")
    flag.Parse()
    if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
        cfg.FlagRunAddr = envRunAddr
    }
    if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
        cfg.FlagLogLevel = envLogLevel
    }
    if envStorageInterval := os.Getenv("STORE_INTERVAL"); envStorageInterval != "" {
        cfg.FlagStoreInterval = envStorageInterval
    }
    if envStorePath := os.Getenv("FILE_STORAGE_PATH"); envStorePath != "" {
        cfg.FlagStorePath = envStorePath
    }
    if envFlagRestore := os.Getenv("RESTORE"); envFlagRestore != "" {
        cfg.FlagRestore = envFlagRestore
    }
    return cfg, err
}



