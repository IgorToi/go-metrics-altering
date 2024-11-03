package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	PollInterval   = 2
	GaugeType      = "gauge"
	CountType      = "counter"
	PollCount      = "PollCount"
	StatusOK       = 200
	ProtocolScheme = "http://"
	configPath     = "config.json"
)

type ConfigAgent struct {
	FlagRunAddr        string        `json:"address"`
	FlagReportInterval int           `json:"report_interval"`
	FlagPollInterval   int           `json:"poll_interval"`
	FlagLogLevel       string        `json:"log_level"`
	FlagHashKey        string        `json:"hash_key"`
	FlagRateLimit      int           `json:"rate_limit"`
	PauseDuration      time.Duration // Time - agent will wait to send metrics again
	URL                string
	FlagCryptoKey      string `json:"crypto_key"`
	FlagConfigName     string `json:"config_name"`
	FlagRSAEncryption bool 
}

func LoadConfig() (*ConfigAgent, error) {
	cfg := new(ConfigAgent)

	configFile, err := os.OpenFile(configPath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		log.Println("error while opening config.json", err)
	}
	defer configFile.Close()

	err = json.NewDecoder(configFile).Decode(&cfg)
	if err != nil {
		log.Println("error while decoding data from config.json", err)
	}

	// var err error
	flag.StringVar(&cfg.FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&cfg.FlagLogLevel, "u", "info", "log level")
	flag.IntVar(&cfg.FlagReportInterval, "r", 10, "frequency of metrics being sent")
	flag.IntVar(&cfg.FlagPollInterval, "p", 0, "frequency of metrics being received from the runtime package")
	flag.IntVar(&cfg.FlagRateLimit, "l", 3, "rate limit")
	flag.StringVar(&cfg.FlagHashKey, "k", "", "hash key")
	flag.StringVar(&cfg.FlagCryptoKey, "crypto-key", "keys/public.pem", "path to public key")
	flag.StringVar(&cfg.FlagConfigName, "c", "config.json", "name of the config with json data")
	flag.BoolVar(&cfg.FlagRSAEncryption, "rsa-bool", false, "whether communication should be encrypted using rsa keys")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.FlagRunAddr = envRunAddr
	}

	if envCofigName := os.Getenv("CONFIG"); envCofigName != "" {
		cfg.FlagConfigName = envCofigName
	}

	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		cfg.FlagRateLimit, err = strconv.Atoi(envRateLimit)
		if err != nil {
			log.Fatal("error while parsing rate limit", err)
		}
	}

	if envRSAKey := os.Getenv("CRYPTO_KEY"); envRSAKey != "" {
		cfg.FlagHashKey = envRSAKey
	}

	if envHashValue := os.Getenv("KEY"); envHashValue != "" {
		cfg.FlagHashKey = envHashValue
	}

	if envRoportInterval := os.Getenv("REPORT_INTERVAL"); envRoportInterval != "" {
		cfg.FlagReportInterval, err = strconv.Atoi(envRoportInterval)
		if err != nil {
			log.Fatal("error while parsing report interval", err)
		}
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		cfg.FlagPollInterval, err = strconv.Atoi(envPollInterval)
		if err != nil {
			log.Fatal("error while parsing poll interval", err)
		}
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		cfg.FlagLogLevel = envLogLevel
	}

	// rename config.json as needed
	err = os.Rename(configPath, cfg.FlagConfigName)
	if err != nil {
		log.Println("error while renaming config.json", err)
	}

	cfg.PauseDuration = time.Duration(cfg.FlagReportInterval) * time.Second
	cfg.URL = ProtocolScheme + cfg.FlagRunAddr
	return cfg, err
}
