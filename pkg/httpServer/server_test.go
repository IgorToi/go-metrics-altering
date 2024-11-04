package httpserver

import (
	"context"
	"testing"

	config "github.com/igortoigildin/go-metrics-altering/config/server"
	server "github.com/igortoigildin/go-metrics-altering/internal/server/api"
	"github.com/igortoigildin/go-metrics-altering/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	cfg := config.ConfigServer{}
	storage := storage.New(&cfg)
	r := server.Router(context.Background(), &cfg, storage)
	s := New(r)
	assert.IsType(t, Server{}, *s)
}
