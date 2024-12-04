package sendmetrics

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	config "github.com/igortoigildin/go-metrics-altering/config/agent"
	"github.com/igortoigildin/go-metrics-altering/internal/models"
	"github.com/igortoigildin/go-metrics-altering/pkg/crypt"
	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
	"go.uber.org/zap"
)

const updEndpoint = "/update/"

var (
	ErrConnectionFailed = errors.New("connection failed")
)

// SendJSONGauge accepts and sends gauge metrics in JSON format to predefined by config server address.
func SendJSONGauge(metricName string, cfg *config.ConfigAgent, value float64) error {
	agent := resty.New()

	if metricName == "" {
		logger.Log.Info("metric data not complete")
		return errors.New("metric data not complete")
	}

	metric := models.GaugeConstructor(value, metricName)
	req := agent.R().SetHeader("Content-Type", "application/json")

	// Add X-Real-IP header as defined by agent config
	req.SetHeader("X-Real-IP", cfg.FlagRealIP)

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

	if cfg.FlagRSAEncryption {
		publicKeyPEM, err := os.ReadFile(cfg.FlagCryptoKey)
		if err != nil {
			logger.Log.Info("error while reading rsa public key:", zap.Error(err))
		}

		// encrypting using public key
		metricsJSON, err = crypt.Encrypt(publicKeyPEM, metricsJSON)
		if err != nil {
			logger.Log.Error("error while encrypting data")
		}
	}

	_, err = req.SetBody(metricsJSON).Post(cfg.URL + updEndpoint)
	if err != nil {
		// send again n times if timeout error
		switch {
		case os.IsTimeout(err):
			for _, delay := range []time.Duration{time.Second, 2 * time.Second, 3 * time.Second} {
				time.Sleep(delay)
				if _, err = req.Post(req.URL); err == nil {
					break
				}
			}
			return ErrConnectionFailed
		default:
			logger.Log.Info("unexpected sending metric error:", zap.Error(err))
			return err
		}
	}

	logger.Log.Info("sent JSON gauge metric:", zap.Float64(metricName, value))
	return nil
}

// SendJSONCounter accepts and sends gauge metrics in JSON format to predefined by config server address.
func SendJSONCounter(counter int, cfg *config.ConfigAgent) error {
	agent := resty.New()

	metric := models.CounterConstructor(int64(counter))
	req := agent.R().SetHeader("Content-Type", "application/json")

	// Add X-Real-IP header as defined by agent config
	req.SetHeader("X-Real-IP", cfg.FlagRealIP)

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

	if cfg.FlagRSAEncryption {
		publicKeyPEM, err := os.ReadFile(cfg.FlagCryptoKey)
		if err != nil {
			logger.Log.Info("error while reading rsa public key:", zap.Error(err))
			return err
		}
		// encrypting using public key
		metricJSON, err = crypt.Encrypt(publicKeyPEM, metricJSON)
		if err != nil {
			logger.Log.Error("error while encrypting data")
			return err
		}
	}

	_, err = req.SetBody(metricJSON).Post(cfg.URL + updEndpoint)
	if err != nil {
		// send again n times if timeout error
		switch {
		case os.IsTimeout(err):
			for _, delay := range []time.Duration{time.Second, 2 * time.Second, 3 * time.Second} {
				time.Sleep(delay)
				if _, err = req.Post(req.URL); err == nil {
					break
				}
			}
			return ErrConnectionFailed
		default:
			logger.Log.Info("unexpected sending metric error:", zap.Error(err))
			return err
		}
	}

	logger.Log.Info("sent JSON counter metric:", zap.Int("conuter", counter))
	return nil
}
