package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
	"go.uber.org/zap"
)

func Encrypt(publicKeyPEM []byte, msg []byte) ([]byte, error) {
	publicKeyBlock, _ := pem.Decode(publicKeyPEM)
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		logger.Log.Info("error while parsing a public key in PKIX:", zap.Error(err))
		return nil, err
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), msg)
	if err != nil {
		logger.Log.Info("error while EncryptPKCS1v15 encrypting:", zap.Error(err))
		return nil, err
	}

	return ciphertext, nil
}
