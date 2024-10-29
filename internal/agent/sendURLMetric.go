package agent

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	config "github.com/igortoigildin/go-metrics-altering/config/agent"
	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
	"go.uber.org/zap"
)

//const urlParams = "/update/{metricType}/{metricName}/{metricValue}"

// SendJSONCounter accepts and sends counter metrics in URL format to predefined by config server address.
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
	_, _ = req.Post(cfg.URL + "/update/{metricType}/{metricName}/{metricValue}")
	// if err != nil {
	// 	switch {
	// 	case os.IsTimeout(err):
	// 		for _, delay := range []time.Duration{time.Second, 2 * time.Second, 3 * time.Second} {
	// 			time.Sleep(delay)
	// 			if _, err = req.Post(cfg.URL + urlParams); err == nil {
	// 				break
	// 			}
	// 			logger.Log.Info("timeout error, server not reachable:", zap.Error(err))
	// 		}
	// 		return ErrConnectionFailed
	// 	default:
	// 		logger.Log.Info("unexpected sending metric error via URL:", zap.Error(err))
	// 		return err
	// 	}
	// }

	logger.Log.Info("sent url counter metric:", zap.Int("conuter", counter))
	return nil
}

// SendURLGauge accepts and sends gauge metrics in URL format to predefined by config server address.
func SendURLGauge(cfg *config.ConfigAgent, value float64, metricName string) error {
	// agent := resty.New()
	// req := agent.R().SetPathParams(map[string]string{
	// 	"metricType":  config.GaugeType,
	// 	"metricName":  metricName,
	// 	"metricValue": strconv.FormatFloat(value, 'f', 6, 64),
	// }).SetHeader("Content-Type", "text/plain")

	if metricName == "" {
		logger.Log.Info("metric data not complete")
		return errors.New("metric data not complete")
	}

	
	
	

	// signing metric value with sha256 and setting header accordingly
	// if cfg.FlagHashKey != "" {
	// 	key := []byte(cfg.FlagHashKey)
	// 	h := hmac.New(sha256.New, key)
	// 	h.Write(nil)
	// 	dst := h.Sum(nil)
	// 	r.Header.Add("HashSHA256", fmt.Sprintf("%x", dst))
	// }

	urlPath := fmt.Sprintf("%s/update/{metricType}/{metricName}/{metricValue}", cfg.URL, )

	var r http.Request
	req, err := http.NewRequest("POST", cfg.URL + "/update/{metricType}/{metricName}/{metricValue}", nil)
	if err != nil {
		fmt.Println(err)
	}

	r.SetPathValue("metricType",  config.GaugeType)
	r.SetPathValue("metricName",  metricName)
	r.SetPathValue("metricValue",  strconv.FormatFloat(value, 'f', 6, 64))

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	// if err != nil {
	// 	switch {
	// 	case os.IsTimeout(err):
	// 		for _, delay := range []time.Duration{time.Second, 2 * time.Second, 3 * time.Second} {
	// 			time.Sleep(delay)
	// 			if _, err = req.Post(cfg.URL + urlParams); err == nil {
	// 				break
	// 			}
	// 			logger.Log.Info("timeout error, server not reachable:", zap.Error(err))
	// 		}
	// 		return ErrConnectionFailed
	// 	default:
	// 		logger.Log.Info("unexpected sending metric error via URL:", zap.Error(err))
	// 		return err
	// 	}
	// }

	logger.Log.Info("sent url gauge metric:", zap.Float64(metricName, value))
	return nil
}
