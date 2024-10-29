package agent

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	config "github.com/igortoigildin/go-metrics-altering/config/agent"
	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
	"go.uber.org/zap"
)

// SendJSONCounter accepts and sends counter metrics in URL format to predefined by config server address.
func sendURLCounter(cfg *config.ConfigAgent, counter int) error {
	path := fmt.Sprintf("/update/%s/%s/%s", config.GaugeType, config.PollCount, strconv.Itoa(counter))
	r, err := http.NewRequest("POST", cfg.URL + path, nil)
	if err != nil {
		logger.Log.Error("error while preparing http request", zap.Error(err))
	}

	// signing metric value with sha256 and setting header accordingly
	if cfg.FlagHashKey != "" {
		key := []byte(cfg.FlagHashKey)
		h := hmac.New(sha256.New, key)
		h.Write(nil)
		dst := h.Sum(nil)
		r.Header = http.Header{
		"HashSHA256": {fmt.Sprintf("%x", dst)},
		}
	}

	client := http.Client{}
	_, err = client.Do(r)
	if err != nil {
		switch {
		case os.IsTimeout(err):
			for _, delay := range []time.Duration{time.Second, 2 * time.Second, 3 * time.Second} {
				time.Sleep(delay)
				if _, err = client.Do(r); err == nil {
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

	logger.Log.Info("sent url counter metric:", zap.Int("conuter", counter))
	return nil
}

// SendURLGauge accepts and sends gauge metrics in URL format to predefined by config server address.
func SendURLGauge(cfg *config.ConfigAgent, value float64, metricName string) error {
	path := fmt.Sprintf("/update/%s/%s/%s", config.GaugeType, metricName, strconv.FormatFloat(value, 'f', 6, 64))
	r, err := http.NewRequest("POST", cfg.URL + path, nil)
	if err != nil {
		logger.Log.Error("error while preparing http request", zap.Error(err))
	}

	// signing metric value with sha256 and setting header accordingly
	if cfg.FlagHashKey != "" {
		key := []byte(cfg.FlagHashKey)
		h := hmac.New(sha256.New, key)
		h.Write(nil)
		dst := h.Sum(nil)
		r.Header = http.Header{
		"HashSHA256": {fmt.Sprintf("%x", dst)},
		}
	}

	client := http.Client{}
	_, err = client.Do(r)
	if err != nil {
		switch {
		case os.IsTimeout(err):
			for _, delay := range []time.Duration{time.Second, 2 * time.Second, 3 * time.Second} {
				time.Sleep(delay)
				if _, err = client.Do(r); err == nil {
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

	logger.Log.Info("sent url gauge metric:", zap.Float64(metricName, value))
	return nil
}
