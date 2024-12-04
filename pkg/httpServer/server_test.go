package httpserver

import (
	"context"
	"net/http"
	"testing"
	"time"

	config "github.com/igortoigildin/go-metrics-altering/config/server"
	server "github.com/igortoigildin/go-metrics-altering/internal/server/http/api"
	"github.com/igortoigildin/go-metrics-altering/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	cfg := config.ConfigServer{}
	storage, _ := storage.New(&cfg)
	r := server.Router(context.Background(), &cfg, storage)
	s := New(r)
	assert.IsType(t, Server{}, *s)
}

func TestServer_GracefulShutdown(t *testing.T) {
	srv := http.Server{}
	s := Server{
		server:          &srv,
		ShutdownTimeout: 1 * time.Second,
	}
	err := s.GracefulShutdown()
	assert.NoError(t, err)
}
