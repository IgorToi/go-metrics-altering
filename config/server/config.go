package config

import (
	"encoding/json"
	"errors"
	"flag"
	"html/template"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	GaugeType  = "gauge"
	CountType  = "counter"
	PollCount  = "PollCount"
	timeout    = 10
	configPath = "config/server/configServer.json"
)

const defaultSrvConfig = `{
    "address": "localhost:8080",
    "restore": true,
    "store_interval": 1,
    "store_file": "/path/to/file.db", 
    "database_dsn": "",
    "crypto_key": "/path/to/key.pem"
}`

var errCfgVarEmpty = errors.New("configs variable not set")

type ConfigServer struct {
	FlagRunAddr       string `json:"address"`
	Template          *template.Template
	FlagLogLevel      string `json:"log_level"`
	FlagStoreInterval int    `json:"store_interval"`
	FlagStorePath     string `json:"store_file"`
	FlagRestore       bool   `json:"restore"`
	FlagDBDSN         string `json:"database_dsn"`
	FlagHashKey       string `json:"hash_key"`
	ContextTimout     time.Duration
	FlagCryptoKey     string `json:"crypto_key"`
	FlagConfigName    string `json:"config_name"`
	FlagRSAEncryption bool
}

func LoadConfig() (*ConfigServer, error) {
	// ps := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",`localhost`, `postgres`, `XXXXX`, `metrics`)
	cfg := new(ConfigServer)

	if err := os.WriteFile(configPath, []byte(defaultSrvConfig), 0666); err != nil {
		log.Println(err)
	}

	configFile, err := os.Open(configPath)
	if err != nil {
		log.Println("error while opening config.json", err)
	}
	defer configFile.Close()

	err = json.NewDecoder(configFile).Decode(&cfg)
	if err != nil {
		log.Println("error while decoding data from config.json", err)
	}

	flag.StringVar(&cfg.FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&cfg.FlagLogLevel, "l", "info", "log level")
	flag.IntVar(&cfg.FlagStoreInterval, "i", 1, "metrics backup interval")
	flag.StringVar(&cfg.FlagStorePath, "f", "/tmp/metrics-db.json", "metrics backup storage path")
	flag.BoolVar(&cfg.FlagRestore, "r", false, "true if load from backup is needed")
	flag.StringVar(&cfg.FlagDBDSN, "d", "", "string with DB DSN")
	flag.StringVar(&cfg.FlagHashKey, "k", "", "hash key")
	flag.StringVar(&cfg.FlagCryptoKey, "crypto-key", "keys", "path to private key")
	flag.StringVar(&cfg.FlagConfigName, "c", "configServer.json", "name of the config with json data")
	flag.BoolVar(&cfg.FlagRSAEncryption, "rsa-bool", false, "whether communication should be encrypted using rsa keys")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.FlagRunAddr = envRunAddr
	}

	if envRSAKey := os.Getenv("CRYPTO_KEY"); envRSAKey != "" {
		cfg.FlagHashKey = envRSAKey
	}

	if envCofigName := os.Getenv("CONFIG"); envCofigName != "" {
		cfg.FlagConfigName = envCofigName
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		cfg.FlagLogLevel = envLogLevel
	}

	if envHashKey := os.Getenv("KEY"); envHashKey != "" {
		cfg.FlagHashKey = envHashKey
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

	if envDBDSN := os.Getenv("DATABASE_DSN"); envDBDSN != "" {
		cfg.FlagDBDSN = envDBDSN
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
	if !cfg.validate() {
		return nil, errCfgVarEmpty
	}

	// rename config.json as needed
	err = os.Rename(configPath, "config/server/"+cfg.FlagConfigName)
	if err != nil {
		log.Println("error while renaming config.json", err)
	}

	cfg.ContextTimout = timeout * time.Second
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
