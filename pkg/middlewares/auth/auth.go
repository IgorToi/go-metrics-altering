// Package auth provides middleware for basic authorization.
package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"

	config "github.com/igortoigildin/go-metrics-altering/config/server"
	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
)

// Auth middleware checks whether IP request is in trusted subnet.
func Auth(next http.HandlerFunc, cfg *config.ConfigServer) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cfg.FlagTrustedSubnet != "" {
			ipStr := r.Header.Get("X-Real-IP")

			// check whether X-Real-IP is the correct IP-address
			ip := net.ParseIP(ipStr)
			if ip == nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			isTrusted, err := isIPInTrustedSubnet(ipStr, cfg.FlagTrustedSubnet)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if !isTrusted {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}
		next(w, r)
	})
}

func isIPInTrustedSubnet(ipStr string, trustedSubnet string) (bool, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		logger.Log.Error("invalid IP address")
		return false, fmt.Errorf("invalid IP address: %s", ipStr)
	}

	// converting trusted subnet to net.IPNet
	_, trustedNet, err := net.ParseCIDR(trustedSubnet)
	if err != nil {
		return false, fmt.Errorf("invalid trusted subnet: %s", trustedSubnet)
	}

	// chesk whether ip is in trusted subnet
	return trustedNet.Contains(ip), nil
}

// ValidMAC reports whether messageMAC is a valid HMAC tag for message.
func ValidMAC(message, messageMAC, key []byte) (bool, error) {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(message))
	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal(messageMAC, []byte(expectedMAC)), nil
}
