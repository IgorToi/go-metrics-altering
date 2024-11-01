package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"os"

	config "github.com/igortoigildin/go-metrics-altering/config/server"
	server "github.com/igortoigildin/go-metrics-altering/internal/server/api"
	storage "github.com/igortoigildin/go-metrics-altering/internal/server/api"
	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
	"go.uber.org/zap"
)

var buildVersion string = "N/A"
var buildDate string = "N/A"
var buildCommit string = "N/A"

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Println("error while logading config", err)
		return
	}

	if err = logger.Initialize(cfg.FlagLogLevel); err != nil {
		log.Println("error while initializing logger", err)
		return
	}

	/////////

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		logger.Log.Error("error while generating RSA private key", zap.Error(err))
		return
	}

	publicKey := &privateKey.PublicKey

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	err = os.WriteFile("keys/private.pem", privateKeyPEM, 0644)
	if err != nil {
		logger.Log.Error("error while writing RSA private key to the file", zap.Error(err))
		return
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		logger.Log.Error("error while converting a public key to PKIX", zap.Error(err))
		return
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	err = os.WriteFile("keys/public.pem", publicKeyPEM, 0644)
	if err != nil {
		logger.Log.Error("error while writing RSA public key to the file", zap.Error(err))
		return
	}

	//////////

	storage := storage.New(cfg)
	r := server.Router(ctx, cfg, storage)

	logger.Log.Info("Starting server on", zap.String("address", cfg.FlagRunAddr))

	err = http.ListenAndServe(cfg.FlagRunAddr, r)
	if err != nil {
		logger.Log.Error("cannot start the server", zap.Error(err))
		return
	}
}
