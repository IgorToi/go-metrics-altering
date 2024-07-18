package config

import (
	"flag"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/igortoigildin/go-metrics-altering/internal/logger"
	"go.uber.org/zap"
)

const (
    PollInterval   = 2
    GaugeType      = "gauge"
    CountType      = "counter"
    PollCount      = "PollCount"
    StatusOK       = 200
    ProtocolScheme = "http://"
)

type ConfigAgent struct {
    FlagRunAddr        string
    FlagReportInterval int
    FlagPollInterval   int
    FlagLogLevel       string
    Rtm                runtime.MemStats
    Memory             map[string]float64
    Count              int
    PauseDuration      time.Duration // Time agent will wait to send metrics again
    URL                string
}

func LoadConfig() (*ConfigAgent, error) {
    cfg := new(ConfigAgent)
    cfg.Memory = make(map[string]float64)
    var err error
    flag.StringVar(&cfg.FlagRunAddr, "a", "localhost:8080", "address and port to run server")
    flag.StringVar(&cfg.FlagLogLevel, "l", "info", "log level")
    flag.IntVar(&cfg.FlagReportInterval, "r", 10, "frequency of metrics being sent to the server")
    flag.IntVar(&cfg.FlagPollInterval, "p", 2, "frequency of metrics being received from the runtime package")
    flag.Parse()
    if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
        cfg.FlagRunAddr = envRunAddr
    }
    if envRoportInterval := os.Getenv("REPORT_INTERVAL"); envRoportInterval != "" {
        cfg.FlagReportInterval, err = strconv.Atoi(envRoportInterval)
        if err != nil {
            logger.Log.Fatal("error while parsing report interval", zap.Error(err))
        }
    }
    if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
        cfg.FlagPollInterval, err = strconv.Atoi(envPollInterval)
        if err != nil {
            logger.Log.Fatal("error while parsing poll interval", zap.Error(err))
        }
    }
    if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
        cfg.FlagLogLevel = envLogLevel
    }
    cfg.PauseDuration = time.Duration(cfg.FlagReportInterval) * time.Second
    cfg.URL = ProtocolScheme + cfg.FlagRunAddr
    return cfg, err
}

func (c *ConfigAgent) UpdateMetrics() {
    PauseDuration := time.Duration(c.FlagPollInterval) * time.Second
    for {
        time.Sleep(PauseDuration)
        runtime.ReadMemStats(&c.Rtm)
        c.Memory["Alloc"] = float64(c.Rtm.Alloc)
        c.Memory["BuckHashSys"] = float64(c.Rtm.BuckHashSys)
        c.Memory["Frees"] = float64(c.Rtm.Frees)
        c.Memory["GCCPUFraction"] = float64(c.Rtm.GCCPUFraction)
        c.Memory["GCSys"] = float64(c.Rtm.GCSys)
        c.Memory["HeapAlloc"] = float64(c.Rtm.HeapAlloc)
        c.Memory["HeapIdle"] = float64(c.Rtm.HeapIdle)
        c.Memory["HeapInuse"] = float64(c.Rtm.HeapInuse)
        c.Memory["HeapObjects"] = float64(c.Rtm.HeapObjects)
        c.Memory["HeapReleased"] = float64(c.Rtm.HeapReleased)
        c.Memory["HeapSys"] = float64(c.Rtm.HeapSys)
        c.Memory["LastGC"] = float64(c.Rtm.LastGC)
        c.Memory["Lookups"] = float64(c.Rtm.Lookups)
        c.Memory["MCacheInuse"] = float64(c.Rtm.MCacheInuse)
        c.Memory["MCacheSys"] = float64(c.Rtm.MCacheSys)
        c.Memory["MSpanInuse"] = float64(c.Rtm.MSpanInuse)
        c.Memory["MSpanSys"] = float64(c.Rtm.MSpanSys)
        c.Memory["Mallocs"] = float64(c.Rtm.Mallocs)
        c.Memory["NextGC"] = float64(c.Rtm.NextGC)
        c.Memory["NumForcedGC"] = float64(c.Rtm.NumForcedGC)
        c.Memory["NumGC"] = float64(c.Rtm.NumGC)
        c.Memory["OtherSys"] = float64(c.Rtm.OtherSys)
        c.Memory["NextGC"] = float64(c.Rtm.NextGC)
        c.Memory["NumForcedGC"] = float64(c.Rtm.NumForcedGC)
        c.Memory["NumGC"] = float64(c.Rtm.NumGC)
        c.Memory["OtherSys"] = float64(c.Rtm.OtherSys)
        c.Memory["PauseTotalNs"] = float64(c.Rtm.PauseTotalNs)
        c.Memory["StackInuse"] = float64(c.Rtm.StackInuse)
        c.Memory["StackSys"] = float64(c.Rtm.StackSys)
        c.Memory["Sys"] = float64(c.Rtm.StackSys)
        c.Memory["TotalAlloc"] = float64(c.Rtm.TotalAlloc)
        c.Memory["RandomValue"] = rand.Float64()
        c.Count++
    }
}

