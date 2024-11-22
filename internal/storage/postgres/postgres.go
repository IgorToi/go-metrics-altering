package psql

import (
	"context"
	"database/sql"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	config "github.com/igortoigildin/go-metrics-altering/config/server"
	"github.com/igortoigildin/go-metrics-altering/internal/models"
	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
	migrate "github.com/igortoigildin/go-metrics-altering/pkg/migrations"
	_ "github.com/lib/pq"

	"go.uber.org/zap"
)

const (
	GaugeType = "gauge"
	CountType = "counter"
	PollCount = "PollCount"
)

type PGStorage struct {
	conn     *sql.DB
	strategy Strategy
}

func New(cfg *config.ConfigServer) (*PGStorage, error) {
	db, err := sql.Open("pgx", cfg.FlagDBDSN)
	if err != nil {
		return nil, err
	}

	if err := migrate.New("migrations", db); err != nil {
		logger.Log.Fatal("could not migrate db: %s", zap.Error(err))
	}

	return &PGStorage{
		conn: db,
	}, nil
}

func (pg *PGStorage) SetStrategy(metricType string) error {
	if metricType == config.CountType {
		count := Count{
			conn: pg.conn,
		}
		pg.strategy = &count
		return nil
	}

	gauge := Gauge{
		conn: pg.conn,
	}
	pg.strategy = &gauge
	return nil
}

func (pg *PGStorage) Ping(ctx context.Context) error {
	err := pg.conn.PingContext(ctx)
	if err != nil {
		logger.Log.Info("connection to the database not alive", zap.Error(err))
	}
	return err
}

func (pg *PGStorage) Update(ctx context.Context, metricType string, metricName string, metricValue any) error {
	pg.SetStrategy(metricType)
	return pg.strategy.Update(ctx, metricType, metricName, metricValue)
}

func (pg *PGStorage) Get(ctx context.Context, metricType string, metricName string) (models.Metrics, error) {
	pg.SetStrategy(metricType)
	return pg.strategy.Get(ctx, metricType, metricName)
}

func (pg *PGStorage) GetAll(ctx context.Context) (map[string]any, error) {
	metrics := make(map[string]any, 33)
	rows, err := pg.conn.QueryContext(ctx, `SELECT name, value FROM gauges WHERE type = ?`, GaugeType)
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
	rows, err = pg.conn.QueryContext(ctx, `SELECT name, value FROM counters WHERE type = ?`, CountType)
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
