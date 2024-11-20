package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"os"
	"strconv"
	"time"
)

const (
	GaugeType = "gauge"
	CountType = "counter"
	PollCount = "PollCount"
	timeout   = 10
	cfgName   = "config/server/configServer.json"
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
	FlagRunAddrHTTP       string `json:"address"`
	FlagRunAddrGRPC       string `json:"address"`
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
	FlagTrustedSubnet string `json:"trusted_subnet"`
}

func LoadConfig() (*ConfigServer, error) {
	// ps := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",`localhost`, `postgres`, `XXXXX`, `metrics`)
	cfg := new(ConfigServer)

	if err := os.WriteFile(cfgName, []byte(defaultSrvConfig), 0777); err != nil {
		fmt.Println(err)
	}

	configFile, err := os.Open(cfgName)
	if err != nil {
		fmt.Println("error while opening config.json", err)
	}
	defer configFile.Close()

	err = json.NewDecoder(configFile).Decode(&cfg)
	if err != nil {
		fmt.Println("error while decoding data from config.json", err)
	}

	flag.StringVar(&cfg.FlagRunAddrHTTP, "a", ":8080", "address and port to run HTTP server")
	flag.StringVar(&cfg.FlagRunAddrGRPC, "a", ":8081", "address and port to run gRPC server")
	flag.StringVar(&cfg.FlagLogLevel, "l", "info", "log level")
	flag.IntVar(&cfg.FlagStoreInterval, "i", 1, "metrics backup interval")
	flag.StringVar(&cfg.FlagStorePath, "f", "/tmp/metrics-db.json", "metrics backup storage path")
	flag.BoolVar(&cfg.FlagRestore, "r", false, "true if load from backup is needed")
	flag.StringVar(&cfg.FlagDBDSN, "d", "", "string with DB DSN")
	flag.StringVar(&cfg.FlagHashKey, "k", "", "hash key")
	flag.StringVar(&cfg.FlagCryptoKey, "crypto-key", "keys", "path to private key")
	flag.StringVar(&cfg.FlagConfigName, "c", "configServer.json", "name of the config with json data")
	flag.BoolVar(&cfg.FlagRSAEncryption, "rsa-bool", false, "whether communication should be encrypted using rsa keys")
	flag.StringVar(&cfg.FlagTrustedSubnet, "t", "127.0.0.0/8", "trusted_subnet")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS_HTTP"); envRunAddr != "" {
		cfg.FlagRunAddrHTTP = envRunAddr
	}

	if envRunAddr := os.Getenv("ADDRESS_GRPC"); envRunAddr != "" {
		cfg.FlagRunAddrGRPC = envRunAddr
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
		v, err := strconv.Atoi(envStorageInterval)
		if err != nil {
			return nil, err
		}
		cfg.FlagStoreInterval = v
	}

	if envStorePath := os.Getenv("FILE_STORAGE_PATH"); envStorePath != "" {
		cfg.FlagStorePath = envStorePath
	}

	if envTrustedSubnet := os.Getenv("TRUSTED_SUBNET"); envTrustedSubnet != "" {
		cfg.FlagTrustedSubnet = envTrustedSubnet
	}

	if envDBDSN := os.Getenv("DATABASE_DSN"); envDBDSN != "" {
		cfg.FlagDBDSN = envDBDSN
	}

	if envFlagRestore := os.Getenv("RESTORE"); envFlagRestore != "" {
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
