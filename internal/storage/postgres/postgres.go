package storage

import (
	"context"
	"database/sql"
	"errors"

	config "github.com/igortoigildin/go-metrics-altering/config/server"
	"github.com/igortoigildin/go-metrics-altering/internal/models"
	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
	"go.uber.org/zap"
)

//TODO

const (
	GaugeType = "gauge"
	CountType = "counter"
	PollCount = "PollCount"
)

type PGStorage struct {
	conn *sql.DB
}

func NewPGStorage(conn *sql.DB) *PGStorage {
	return &PGStorage{
		conn: conn,
	}
}

func InitPostgresRepo(ctx context.Context, cfg *config.ConfigServer) *PGStorage {
	dbDSN := cfg.FlagDBDSN
	conn, err := sql.Open("pgx", dbDSN)
	if err != nil {
		logger.Log.Debug("error while connecting to DB", zap.Error(err))
	}
	rep := NewPGStorage(conn)
	_, err = rep.conn.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS counters (id SERIAL, name TEXT NOT NULL,"+
		"type TEXT NOT NULL, value bigint, primary key(name));")
	if err != nil {
		logger.Log.Debug("error while creating counters table", zap.Error(err))
	}
	_, err = rep.conn.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS gauges (id SERIAL, name TEXT NOT NULL,"+
		"type TEXT NOT NULL, value DOUBLE PRECISION, primary key(name));")
	if err != nil {
		logger.Log.Debug("error while creating gauges table", zap.Error(err))
	}
	return rep
}

func (rep *PGStorage) Ping(ctx context.Context) error {
	err := rep.conn.PingContext(ctx)
	if err != nil {
		logger.Log.Info("connection to the database not alive", zap.Error(err))
	}
	return err
}

func (rep *PGStorage) Update(ctx context.Context, metricType string, metricName string, metricValue any) error {
	switch metricType {
	case GaugeType:
		_, err := rep.conn.ExecContext(ctx, "INSERT INTO gauges(name, type, value) VALUES($1, $2, $3)"+
			"ON CONFLICT (name) DO UPDATE SET value = $3", metricName, GaugeType, metricValue)
		if err != nil {
			logger.Log.Debug("error while saving gauge metric to the db", zap.Error(err))
			return err
		}
	case CountType:
		_, err := rep.conn.ExecContext(ctx, "INSERT INTO counters(name, type, value) VALUES($1, $2, $3)"+
			"ON CONFLICT (name) DO UPDATE SET value = counters.value + $3", metricName, CountType, metricValue)
		if err != nil {
			logger.Log.Debug("error while saving counter metric to the db", zap.Error(err))
			return err
		}
	}
	return nil
}

func (rep *PGStorage) Get(ctx context.Context, metricType string, metricName string) (models.Metrics, error) {
	var metric models.Metrics
	switch metricType {
	case GaugeType:
		var metric models.Metrics
		err := rep.conn.QueryRowContext(ctx, "SELECT name, type, value FROM gauges WHERE name = $1", metricName).Scan(
			&metric.ID, &metric.MType, &metric.Value)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			logger.Log.Debug("no rows selected", zap.Error(err))
			return metric, err
		case err != nil:
			logger.Log.Debug("error while obtaining metrics", zap.Error(err))
			return metric, err
		default:
			return metric, nil
		}
	case CountType:
		var metric models.Metrics
		err := rep.conn.QueryRowContext(ctx, "SELECT name, type, value FROM counters WHERE name = $1", metricName).Scan(
			&metric.ID, &metric.MType, &metric.Delta)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			logger.Log.Debug("no rows selected", zap.Error(err))
			return metric, err
		case err != nil:
			logger.Log.Debug("error while obtaining metrics", zap.Error(err))
			return metric, err
		default:
			return metric, nil
		}
	}
	return metric, nil
}

func (rep *PGStorage) GetAll(ctx context.Context) (map[string]any, error) {
	metrics := make(map[string]any, 33)
	rows, err := rep.conn.QueryContext(ctx, "SELECT name, value FROM gauges WHERE type = $1", GaugeType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var value any
		err = rows.Scan(&name, &value)
		if err != nil {
			return nil, err
		}
		metrics[name] = value
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	rows, err = rep.conn.QueryContext(ctx, "SELECT name, value FROM counters WHERE type = $1", CountType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var value any
		err = rows.Scan(&name, &value)
		if err != nil {
			return nil, err
		}
		metrics[name] = value
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return metrics, nil
}
