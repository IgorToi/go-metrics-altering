package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
)

func Decrypt(privateKeyPEM []byte, input []byte) ([]byte, error) {

	privateKeyBlock, _ := pem.Decode(privateKeyPEM)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		logger.Log.Error("error while parsing PKCS1PrivateKey")
		return nil, err
	}

	res, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, input)
	if err != nil {
		logger.Log.Error("error while decrypting PKCS1PrivateKey")
		return nil, err
	}

	return res, nil
}
