package storage

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	config "github.com/igortoigildin/go-metrics-altering/config/server"
	"github.com/igortoigildin/go-metrics-altering/internal/logger"
	"github.com/igortoigildin/go-metrics-altering/internal/models"
	"go.uber.org/zap"
)

const (
	GaugeType = "gauge"
	CountType = "counter"
	PollCount = "PollCount"
)

var (
	ErrMetricNotFound = errors.New("metric not found")
	ErrSQLExecution   = errors.New("sql query execution failed")
)

type Repository struct {
	conn *sql.DB
}

func NewRepository(conn *sql.DB) *Repository {
	return &Repository{
		conn: conn,
	}
}

func InitPostgresRepo(c context.Context, cfg *config.ConfigServer) *Repository {
	dbDSN := cfg.FlagDBDSN
	conn, err := sql.Open("pgx", dbDSN)
	if err != nil {

		logger.Log.Fatal("error while connecting to DB", zap.Error(err))
		//
		for n, t := 1, 1; n <= 3; n++ {
			time.Sleep(time.Duration(t) * time.Second)
			if _, err = sql.Open("pgx", dbDSN); err == nil {
				logger.Log.Info("Conn opened successfully")
				break
			}
			t += 2
		}
		//

	}
	rep := NewRepository(conn)
	ctx, cancel := context.WithCancel(c)
	defer cancel()
	// check connection
	if err = rep.conn.PingContext(ctx); err != nil {
		logger.Log.Fatal("error while connecting to DB", zap.Error(err))
	}
	// start Tx
	tx, err := rep.conn.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.Fatal("error while starting Tx", zap.Error(err))
	}
	defer tx.Rollback()
	tx.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS counters (id SERIAL PRIMARY KEY, name TEXT NOT NULL,"+
		"type TEXT NOT NULL, value bigint);")
	tx.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS gauges (id SERIAL PRIMARY KEY, name TEXT NOT NULL,"+
		"type TEXT NOT NULL, value DOUBLE PRECISION);")
	tx.Commit()
	return rep
}

func (rep *Repository) Exist(ctx context.Context, metricType string, metricName string) bool {
	switch metricType {
	case GaugeType:
		var exists bool
		row := rep.conn.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM gauges WHERE name = $1)", metricName)
		if err := row.Scan(&exists); err != nil {
			logger.Log.Fatal("error while checking existence", zap.Error(err))
		}
		return exists
	case CountType:
		var exists bool
		row := rep.conn.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM counters WHERE name = $1)", metricName)
		if err := row.Scan(&exists); err != nil {
			logger.Log.Fatal("error while checking existence", zap.Error(err))
		}
		return exists
	}
	return false
}

func (rep *Repository) Add(ctx context.Context, metricType string, metricName string, metricValue any) error {
	tx, err := rep.conn.Begin()
	if err != nil {
		logger.Log.Debug("error while strarting tx", zap.Error(err))
	}
	switch metricType {
	case GaugeType:
		_, err := tx.ExecContext(ctx, "INSERT INTO gauges(name, type, value) VALUES($1, $2, $3)", metricName, GaugeType, metricValue)
		if err != nil {
			tx.Rollback()
			logger.Log.Fatal("error while saving gauge metric to the db", zap.Error(err))
			return ErrSQLExecution
		}
	case CountType:
		_, err := tx.ExecContext(ctx, "INSERT INTO counters(name, type, value) VALUES($1, $2, $3)", metricName, CountType, metricValue)
		if err != nil {
			tx.Rollback()
			logger.Log.Fatal("error while saving counter metric to the db", zap.Error(err))
			return ErrSQLExecution
		}
	}
	return tx.Commit()
}

func (rep *Repository) Update(ctx context.Context, metricType string, metricName string, metricValue any) error {
	tx, err := rep.conn.Begin()
	if err != nil {
		logger.Log.Debug("error while strarting tx", zap.Error(err))
	}
	switch metricType {
	case GaugeType:
		_, err := tx.ExecContext(ctx, "UPDATE gauges SET value = $1 WHERE name = $2", metricValue, metricName)
		if err != nil {
			tx.Rollback()
			logger.Log.Fatal("error while updating counter metric", zap.Error(err))
			return ErrSQLExecution
		}
	case CountType:
		_, err := tx.ExecContext(ctx, "UPDATE counters SET value = value + $1 WHERE name = $2", metricValue, metricName)
		if err != nil {
			tx.Rollback()
			logger.Log.Fatal("error while saving counter metric to the db", zap.Error(err))
			return ErrSQLExecution
		}
	}
	return tx.Commit()
}

func (rep *Repository) Get(ctx context.Context, metricType string, metricName string) (models.Metrics, error) {
	var metric models.Metrics
	switch metricType {
	case GaugeType:
		var metric models.Metrics
		err := rep.conn.QueryRowContext(ctx, "SELECT name, type, value FROM gauges WHERE name = $1", metricName).Scan(
			&metric.ID, &metric.MType, &metric.Value)
		switch {
		case err == sql.ErrNoRows:
			logger.Log.Fatal("no rows selected", zap.Error(err))
			return metric, ErrMetricNotFound
		case err != nil:
			logger.Log.Fatal("error while obtaining metrics", zap.Error(err))
			return metric, err
		default:
			return metric, nil
		}
	case CountType:
		var metric models.Metrics
		err := rep.conn.QueryRowContext(ctx, "SELECT name, type, value FROM counters WHERE name = $1", metricName).Scan(
			&metric.ID, &metric.MType, &metric.Delta)
		switch {
		case err == sql.ErrNoRows:
			logger.Log.Fatal("no rows selected", zap.Error(err))
			return metric, ErrMetricNotFound
		case err != nil:
			logger.Log.Fatal("error while obtaining metrics", zap.Error(err))
			return metric, err
		default:
			return metric, nil
		}
	}
	return metric, nil
}

func (rep *Repository) Ping(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := rep.conn.PingContext(ctx); err != nil {
		logger.Log.Info("error", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rep.conn.Close()
	w.WriteHeader(http.StatusOK)
}

func (rep *Repository) GetAll(ctx context.Context) (map[string]any, error) {
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
