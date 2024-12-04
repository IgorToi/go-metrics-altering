package storage

import (
	"context"
	"errors"

	config "github.com/igortoigildin/go-metrics-altering/config/server"
	"github.com/igortoigildin/go-metrics-altering/internal/models"
	local "github.com/igortoigildin/go-metrics-altering/internal/storage/inmemory"
	psql "github.com/igortoigildin/go-metrics-altering/internal/storage/postgres"
)

//go:generate go run github.com/vektra/mockery/v2@v2.45.0 --name=Storage
type Storage interface {
	Update(ctx context.Context, metricType string, metricName string, metricValue any) error
	Get(ctx context.Context, metricType string, metricName string) (models.Metrics, error)
	GetAll(ctx context.Context) (map[string]any, error)
	Ping(ctx context.Context) error
}

func New(cfg *config.ConfigServer) (Storage, error) {
	if cfg.FlagDBDSN != "" {
		storage, err := psql.New(cfg)
		if err != nil {
			return nil, errors.New("failed to init storage")
		}
		return storage, nil
	}

	memory := local.New()

	if cfg.FlagRestore {
		err := memory.LoadMetricsFromFile(cfg.FlagStorePath)
		if err != nil {
			return memory, err
		}
	}
	if cfg.FlagStorePath != "" {
		go memory.SaveAllMetricsToFile(cfg.FlagStoreInterval, cfg.FlagStorePath, cfg.FlagStorePath)
	}
	return memory, nil
}
