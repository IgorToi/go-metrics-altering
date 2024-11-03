package storage

import (
	"testing"

	config "github.com/igortoigildin/go-metrics-altering/config/server"
	local "github.com/igortoigildin/go-metrics-altering/internal/storage/inmemory"
	psql "github.com/igortoigildin/go-metrics-altering/internal/storage/postgres"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	cfg, _ := config.LoadConfig()
	s := New(cfg)

	_, ok := s.(*local.LocalStorage)
	assert.True(t, ok)

	cfg.FlagDBDSN = "temp"
	p := New(cfg)

	_, ok = p.(*psql.PGStorage)
	assert.True(t, ok)
}
