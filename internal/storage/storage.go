package storage

import (
	"context"

	"github.com/igortoigildin/go-metrics-altering/internal/models"
)

type Storage interface {
	Exist(ctx context.Context, metricType string, metricName string) (bool)
	Add(ctx context.Context, metricType string, metricName string, metricValue any) (error)
	Update(ctx context.Context, metricType string, metricName string, metricValue any) (error)
	Get(ctx context.Context, metricType string, metricName string) (models.Metrics, error)
	// Ping(ctx context.Context, metricType string, metricName string) (error)
	GetAll(ctx context.Context) (map[string]any, error)
}