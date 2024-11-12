package storage

import (
	"testing"

	config "github.com/igortoigildin/go-metrics-altering/config/server"
	local "github.com/igortoigildin/go-metrics-altering/internal/storage/inmemory"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	cfg := config.ConfigServer{}
	cfg.FlagStorePath = "temp"
	s := New(&cfg)

	_, ok := s.(*local.LocalStorage)
	assert.True(t, ok)
}
