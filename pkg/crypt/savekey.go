package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"

	config "github.com/igortoigildin/go-metrics-altering/config/server"
	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
	"go.uber.org/zap"
)

// InitRSAKeys generates and saves private and public rsa keys.
func InitRSAKeys(cfg *config.ConfigServer) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		logger.Log.Error("error while generating RSA private key", zap.Error(err))
		return err
	}

	publicKey := &privateKey.PublicKey

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	err = os.WriteFile(cfg.FlagCryptoKey+"/private.pem", privateKeyPEM, 0644)
	if err != nil {
		logger.Log.Error("error while writing RSA private key to the file", zap.Error(err))
		return err
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		logger.Log.Error("error while converting a public key to PKIX", zap.Error(err))
		return err
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	err = os.WriteFile(cfg.FlagCryptoKey+"/public.pem", publicKeyPEM, 0644)
	if err != nil {
		logger.Log.Error("error while writing RSA public key to the file", zap.Error(err))
		return err
	}

	return nil
}
