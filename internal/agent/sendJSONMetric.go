package agent

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	config "github.com/igortoigildin/go-metrics-altering/config/agent"
	"github.com/igortoigildin/go-metrics-altering/internal/models"
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

	////
	publicKeyPEM, err := os.ReadFile("public.pem")
	if err != nil {
		logger.Log.Info("error while reading rsa public key:", zap.Error(err))
		return err
	}
	publicKeyBlock, _ := pem.Decode(publicKeyPEM)
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		logger.Log.Info("error while parsing a public key in PKIX:", zap.Error(err))
		return err
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), compressedRequest.Bytes())
	if err != nil {
		panic(err)
	}
	/////

	_, err = req.SetBody(ciphertext).Post(cfg.URL + updEndpoint)
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

	logger.Log.Info("sent JSON gauge metric:", zap.Float64(metricName, value))
	return nil
}

// SendJSONCounter accepts and sends gauge metrics in JSON format to predefined by config server address.
func SendJSONCounter(counter int, cfg *config.ConfigAgent) error {
	agent := resty.New()

	metric := models.CounterConstructor(int64(counter))
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
	_, err = req.SetBody(compressedRequest.Bytes()).Post(cfg.URL + updEndpoint)
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

	logger.Log.Info("sent JSON counter metric:", zap.Int("conuter", counter))
	return nil
}
