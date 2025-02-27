package psql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/igortoigildin/go-metrics-altering/internal/models"
	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
	"go.uber.org/zap"
)

type Strategy interface {
	Update(ctx context.Context, metricType string, metricName string, metricValue any) error
	Get(ctx context.Context, metricType string, metricName string) (models.Metrics, error)
}

type Count struct {
	conn *sql.DB
}

func (c *Count) Update(ctx context.Context, metricType string, metricName string, metricValue any) error {
	_, err := c.conn.ExecContext(ctx, `INSERT INTO counters(name, type, value) VALUES($1, $2, $3) ON CONFLICT (name) DO UPDATE SET value = counters.value + $3`, metricName, metricType, metricValue)
	return err
}

func (c *Count) Get(ctx context.Context, metricType string, metricName string) (models.Metrics, error) {
	var metric models.Metrics
	err := c.conn.QueryRowContext(ctx, "SELECT name, type, value FROM counters WHERE name = ?", metricName).Scan(
		&metric.ID, &metric.MType, &metric.Delta)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return metric, err
	case err != nil:
		logger.Log.Info("error while obtaining metrics", zap.Error(err))
		return metric, err
	}
	return metric, nil
}

type Gauge struct {
	conn *sql.DB
}

func (g *Gauge) Update(ctx context.Context, metricType string, metricName string, metricValue any) error {
	_, err := g.conn.ExecContext(ctx, "INSERT INTO gauges(name, type, value) VALUES($1, $2, $3)"+
		"ON CONFLICT (name) DO UPDATE SET value = $3", metricName, metricType, metricValue)
	return err
}

func (g *Gauge) Get(ctx context.Context, metricType string, metricName string) (models.Metrics, error) {
	var metric models.Metrics
	err := g.conn.QueryRowContext(ctx, "SELECT name, type, value FROM gauges WHERE name = $1", metricName).Scan(
		&metric.ID, &metric.MType, &metric.Value)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		logger.Log.Info("no rows selected", zap.Error(err))
		return metric, err
	case err != nil:
		return metric, err
	}
	return metric, nil
}
