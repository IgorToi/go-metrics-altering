package agent

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	config "github.com/igortoigildin/go-metrics-altering/config/agent"
	"github.com/igortoigildin/go-metrics-altering/internal/agent/memory"
	"github.com/igortoigildin/go-metrics-altering/internal/logger"
	"github.com/igortoigildin/go-metrics-altering/internal/models"
	"go.uber.org/zap"
)


const urlParams = "/update/{metricType}/{metricName}/{metricValue}"

var (
	ErrConnectionFailed = errors.New("connection failed")
)

func RunAgent() {
    cfg := Initialize()
    var wg sync.WaitGroup
    jobs := make(chan models.Metrics, 33)
    memoryStats := memory.NewMemoryStats()

    wg.Add(1)
    go func() {
        defer wg.Done()
        go memoryStats.UpdateMetrics(cfg)  
    }()

    wg.Add(1)
    go func() {
        defer wg.Done()
        memoryStats.UpdateVirtualMemoryStat(cfg)  
    }()


    for w := 1; w <= cfg.FlagRateLimit; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            worker(jobs, cfg)
        }()
    }

    wg.Add(1)
    go func() {
        defer wg.Done()
        memoryStats.ReadMetrics(cfg, jobs)
    }()
    
    wg.Wait()
    


    // go MemoryStats.SendJSONMetrics(cfg)  // v2 - metrics in body json
    // go MemoryStats.SendBatchMetrics(cfg) // v3 - sending batchs of metrics json
    // MemoryStats.SendURLmetrics(cfg)      // v1 - metrics in url path
}




func worker(jobs <-chan models.Metrics, cfg *config.ConfigAgent) {
    for {
        time.Sleep(cfg.PauseDuration)
        for j := range jobs {
            switch j.MType {
            case config.CountType:
                err := sendURLCounter(cfg, int(*j.Delta))
                if err != nil {
                    logger.Log.Info("unexpected sending url counter metric error:", zap.Error(err))
                }
                err = SendJSONCounter(int(*j.Delta), cfg)
                if err != nil {
                    logger.Log.Info("unexpected sending json counter metric error:", zap.Error(err))
                }
            case config.GaugeType:
                err := SendURLGauge(cfg, *j.Value, j.ID)
                if err != nil {
                    logger.Log.Info("unexpected sending url gauge metric error:", zap.Error(err))
                }
                err = SendJSONGauge(j.ID, cfg, *j.Value)
                if err != nil {
                    logger.Log.Info("unexpected sending json gauge metric error:", zap.Error(err))
                }
            }
        }
    }
}


func Initialize() *config.ConfigAgent {
    cfg, err := config.LoadConfig()
    if err != nil {
        logger.Log.Fatal("error while logading config", zap.Error(err))
    }
    if err := logger.Initialize(cfg.FlagLogLevel); err != nil {
        logger.Log.Fatal("error while initializing logger", zap.Error(err))
    }
    return cfg
}


func sendURLCounter(cfg *config.ConfigAgent, counter int) error {
    agent := resty.New()
    req := agent.R().SetPathParams(map[string]string{
        "metricType":  config.CountType,
        "metricName":  config.PollCount,
        "metricValue": strconv.Itoa(counter),
    }).SetHeader("Content-Type", "text/plain")
    // signing metric value with sha256 and setting header accordingly
    if cfg.FlagHashKey != "" {
        key := []byte(cfg.FlagHashKey)
        h := hmac.New(sha256.New, key)
        h.Write(nil)
        dst := h.Sum(nil)
        req.SetHeader("HashSHA256", fmt.Sprintf("%x", dst))
    }
    _, err := req.Post(cfg.URL + "/update/{metricType}/{metricName}/{metricValue}")
    if err != nil {
        switch {
        case os.IsTimeout(err):
            for _, delay := range []time.Duration{time.Second, 2 * time.Second, 3 * time.Second} {
                time.Sleep(delay)
                if _, err = req.Post(cfg.URL + urlParams); err == nil {
                    break
                }
                logger.Log.Info("timeout error, server not reachable:", zap.Error(err))
            }
            return ErrConnectionFailed
        default:
            logger.Log.Info("unexpected sending metric error via URL:", zap.Error(err))
            return err
        }
    }
    return nil
}

func SendURLGauge(cfg *config.ConfigAgent, value float64, metricName string) error {
    agent := resty.New()
    req := agent.R().SetPathParams(map[string]string{
        "metricType":  config.GaugeType,
        "metricName":  metricName,
        "metricValue": strconv.FormatFloat(value, 'f', 6, 64),
    }).SetHeader("Content-Type", "text/plain")
    // signing metric value with sha256 and setting header accordingly
    if cfg.FlagHashKey != "" {
        key := []byte(cfg.FlagHashKey)
        h := hmac.New(sha256.New, key)
        h.Write(nil)
        dst := h.Sum(nil)
        req.SetHeader("HashSHA256", fmt.Sprintf("%x", dst))
    }
    _, err := req.Post(cfg.URL + "/update/{metricType}/{metricName}/{metricValue}")
    if err != nil {
        switch {
        case os.IsTimeout(err):
            for _, delay := range []time.Duration{time.Second, 2 * time.Second, 3 * time.Second} {
                time.Sleep(delay)
                if _, err = req.Post(cfg.URL + urlParams); err == nil {
                    break
                }
                logger.Log.Info("timeout error, server not reachable:", zap.Error(err))
            }
            return ErrConnectionFailed
        default:
            logger.Log.Info("unexpected sending metric error via URL:", zap.Error(err))
            return err
        }
    }
    return nil
}

func SendJSONGauge(metricName string, cfg *config.ConfigAgent, value float64) error {
    agent := resty.New()
    metric := models.Metrics{
        ID:    metricName,
        MType: config.GaugeType,
        Value: &value,
    }
    req := agent.R().SetHeader("Content-Type", "application/json").SetHeader("Content-Encoding", "gzip").
        SetHeader("Accept-Encoding", "gzip")
    
    metricsJSON, err := json.Marshal(metric)
    if err != nil {
        logger.Log.Info("marshalling json error:", zap.Error(err))
        return err
    }
    // signing metric value with sha256 and setting header accordingly
    if cfg.FlagHashKey != "" {
        key := []byte(cfg.FlagHashKey)
        h := hmac.New(sha256.New, key)
        h.Write(metricsJSON)
        dst := h.Sum(nil)
        req.SetHeader("HashSHA256", fmt.Sprintf("%x", dst))
    }
    req.Method = resty.MethodPost
    var compressedRequest bytes.Buffer
    writer := gzip.NewWriter(&compressedRequest)
    _, err = writer.Write(metricsJSON)
    if err != nil {
        logger.Log.Info("error while compressing request:", zap.Error(err))
        return err
    }
    err = writer.Close()
    if err != nil {
        logger.Log.Info("error while closing gzip writer:", zap.Error(err))
        return err
    }
    _, err = req.SetBody(compressedRequest.Bytes()).Post(cfg.URL + "/update/")
    if err != nil {
        // send again n times if timeout error
        switch {
        case os.IsTimeout(err):
            for _, delay := range []time.Duration{time.Second, 2 * time.Second, 3 * time.Second} {
                time.Sleep(delay)
                if _, err = req.Post(req.URL); err == nil {
                    break
                }
                logger.Log.Info("timeout error, server not reachable:", zap.Error(err))
            }
            return ErrConnectionFailed
        default:
            logger.Log.Info("unexpected sending metric error:", zap.Error(err))
            return err
        }
    }
    return nil
}

func SendJSONCounter(counter int, cfg *config.ConfigAgent) error {
    agent := resty.New()
    valueDelta := int64(counter)
    metric := models.Metrics{
        ID:    config.PollCount,
        MType: config.CountType,
        Delta: &valueDelta,
    }
    req := agent.R().SetHeader("Content-Type", "application/json").SetHeader("Content-Encoding", "gzip").
        SetHeader("Accept-Encoding", "gzip")
    metricJSON, err := json.Marshal(metric)
    if err != nil {
        logger.Log.Info("marshalling json error:", zap.Error(err))
        return err
    }
    // signing metric value with sha256 and setting header accordingly
    if cfg.FlagHashKey != "" {
        key := []byte(cfg.FlagHashKey)
        h := hmac.New(sha256.New, key)
        h.Write(metricJSON)
        dst := h.Sum(nil)
        req.SetHeader("HashSHA256", fmt.Sprintf("%x", dst))
    }
    var compressedRequest bytes.Buffer
    writer := gzip.NewWriter(&compressedRequest)
    _, err = writer.Write(metricJSON)
    if err != nil {
        logger.Log.Info("error while compressing request:", zap.Error(err))
        return err
    }
    err = writer.Close()
    if err != nil {
        logger.Log.Info("error while closing gzip writer:", zap.Error(err))
        return err
    }
    _, err = req.SetBody(compressedRequest.Bytes()).Post(cfg.URL + "/update/")
    if err != nil {
        // send again n times if timeout error
        switch {
        case os.IsTimeout(err):
            for _, delay := range []time.Duration{time.Second, 2 * time.Second, 3 * time.Second} {
                time.Sleep(delay)
                if _, err = req.Post(req.URL); err == nil {
                    break
                }
                logger.Log.Info("timeout error, server not reachable:", zap.Error(err))
            }
            return ErrConnectionFailed
        default:
            logger.Log.Info("unexpected sending metric error via URL:", zap.Error(err))
            return err
        }
    }
    return nil
}



