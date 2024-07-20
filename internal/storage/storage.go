package storage

import (
	"context"

	"github.com/igortoigildin/go-metrics-altering/internal/models"
)

type Storage interface {
	Update(ctx context.Context, metricType string, metricName string, metricValue any) error
	Get(ctx context.Context, metricType string, metricName string) (models.Metrics, error)
	GetAll(ctx context.Context) (map[string]any, error)
	Ping(ctx context.Context) error
}
