// Package auth provides middleware for basic authorization.
package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// ValidMAC reports whether messageMAC is a valid HMAC tag for message.
func ValidMAC(message, messageMAC, key []byte) (bool, error) {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(message))
	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal(messageMAC, []byte(expectedMAC)), nil
}
