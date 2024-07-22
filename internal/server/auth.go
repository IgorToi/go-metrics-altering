package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"net/http"

	config "github.com/igortoigildin/go-metrics-altering/config/server"
	"github.com/igortoigildin/go-metrics-altering/internal/logger"
	"go.uber.org/zap"
)

func auth(next http.HandlerFunc, cfg *config.ConfigServer) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := []byte(cfg.FlagHashKey)
		if len(key) > 0 {
			messageMAC := []byte(r.Header.Get("HashSHA256"))
			var message []byte
			dec := json.NewDecoder(r.Body)
			if err := dec.Decode(&message); err != nil {
				logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if !ValidMAC(message, messageMAC, key) {
				logger.Log.Info("wrong hash")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		next(w, r)
	})
}

// ValidMAC reports whether messageMAC is a valid HMAC tag for message.
func ValidMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}
