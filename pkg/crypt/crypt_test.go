package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	config "github.com/igortoigildin/go-metrics-altering/config/server"
	"github.com/stretchr/testify/assert"
)

func TestGenerateRSAKeys(t *testing.T) {
	cfg := config.ConfigServer{}
	_, _, err := GenerateRSAKeys(&cfg)
	assert.NoError(t, err)
}

func TestEncrypt(t *testing.T) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	publicKey := &privateKey.PublicKey
	publicKeyBytes, _ := x509.MarshalPKIXPublicKey(publicKey)
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	_, err := Encrypt(publicKeyPEM, []byte("test data"))
	assert.NoError(t, err)
}

func TestDecrypt(t *testing.T) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	publicKey := &privateKey.PublicKey
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	publicKeyBytes, _ := x509.MarshalPKIXPublicKey(publicKey)
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	res, err := Encrypt(publicKeyPEM, []byte("test data"))
	assert.NoError(t, err)

	_, err = Decrypt(privateKeyPEM, res)
	assert.NoError(t, err)
}
