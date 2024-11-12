package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/fs"
	"os"

	config "github.com/igortoigildin/go-metrics-altering/config/server"
)

const (
	publicKey  = "public.pem"
	privateKey = "private.pem"
	keysDir    = "keys"
)

func InitRSAKeys(cfg *config.ConfigServer) error {
	privateKeyPEM, publicKeyPEM, err := GenerateRSAKeys(cfg)
	if err != nil {
		return err
	}

	err = os.MkdirAll("keys", 0777)
	if err != nil {
		return err
	}

	err = saveKey("/"+privateKey, privateKeyPEM, 0777)
	if err != nil {
		return err
	}

	err = saveKey("/"+publicKey, publicKeyPEM, 0777)
	if err != nil {
		return err
	}
	return nil
}

// InitRSAKeys generates and saves private and public rsa keys.
func GenerateRSAKeys(cfg *config.ConfigServer) ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	publicKey := &privateKey.PublicKey

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, nil, err
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return privateKeyPEM, publicKeyPEM, nil
}

func saveKey(name string, data []byte, perm fs.FileMode) error {
	err := os.WriteFile(keysDir+name, data, perm)
	if err != nil {
		return err
	}
	return nil
}
