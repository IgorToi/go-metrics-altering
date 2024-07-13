package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

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
type Repository struct {
	DB 	*sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func InitPostgresRepo(c context.Context, cfg *config.ConfigServer) *Repository {
	dbDSN := cfg.FlagDBDSN
	db, err := sql.Open("pgx", dbDSN)
	if err != nil {
		logger.Log.Fatal("error while connecting to DB", zap.Error(err))
	}
	ctx, cancel := context.WithCancel(c)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		logger.Log.Fatal("error while connecting to DB", zap.Error(err))
	}
	_, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS counters (id SERIAL PRIMARY KEY, name TEXT NOT NULL," +
	"type TEXT NOT NULL, value int);")
	if err != nil {
		logger.Log.Fatal("error while creating counters table", zap.Error(err))
	}
	_, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS gauges (id SERIAL PRIMARY KEY, name TEXT NOT NULL," +
	"type TEXT NOT NULL, value DOUBLE PRECISION);")
	if err != nil {
		logger.Log.Fatal("error while creating gauges table", zap.Error(err))
	}
	rep := NewRepository(db)
	return rep
}

func (rep *Repository) Exist(ctx context.Context, metricType string, metricName string) (bool) {
	switch metricType {
	case GaugeType:
		var check bool
		err := rep.DB.QueryRowContext(ctx, "SELECT EXISTS(SELECT type FROM gauges WHERE name = $1)", metricName).Scan(&check)

		switch {
		case err == sql.ErrNoRows:
			fmt.Println("NOT EXIST " + metricName )
			return false
		case err != nil:
			logger.Log.Fatal("query error:", zap.Error(err))
			return false
		default:
			return true
		}
	case CountType:
		var check bool
		err := rep.DB.QueryRowContext(ctx, "SELECT EXISTS(SELECT type FROM counters WHERE name = $1)", metricName).Scan(&check)
		switch {
		case err == sql.ErrNoRows:
			fmt.Println("NOT EXIST " + metricName )
			return false
		case err != nil:
			logger.Log.Fatal("query error:", zap.Error(err))
			return false
		default:
			return true
		}
	}
	return false
}

func (rep *Repository) Add(ctx context.Context, metricType string, metricName string, metricValue any) error {
	switch metricType {
	case GaugeType:
		result, err := rep.DB.ExecContext(ctx, "INSERT INTO gauges(name, type, value) VALUES($1, $2, $3)", metricName, GaugeType, metricValue)
		if err != nil {
			logger.Log.Fatal("error while saving gauge metric to the db", zap.Error(err))
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			logger.Log.Fatal("error while saving gauge metric to the db", zap.Error(err))
		}
		if rows != 1 {
			log.Fatalf("expected to affect (ADD) 1 row, affected %d", rows)
		}
	case CountType:
		result, err := rep.DB.ExecContext(ctx, "INSERT INTO counters(name, type, value) VALUES($1, $2, $3)", metricName, CountType, metricValue)
		if err != nil {
			logger.Log.Fatal("error while saving counter metric to the db", zap.Error(err))
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			logger.Log.Fatal("error while saving counter metric to the db", zap.Error(err))
		}
		if rows != 1 {
			log.Fatalf("expected to affect (ADD) 1 row, affected %d", rows)
		}
	}
	return nil
}

func (rep *Repository) Update(ctx context.Context, metricType string, metricName string, metricValue any) error {
	switch metricType {
	case GaugeType:
		result, err := rep.DB.ExecContext(ctx, "UPDATE gauges SET value = $1 WHERE name = $2", metricValue, metricName)
		if err != nil {
			logger.Log.Fatal("error while updating counter metric to the db", zap.Error(err))
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			logger.Log.Fatal("error while updating counter metric to the db", zap.Error(err))
		}
		if rows != 1 {
			log.Fatalf("expected to affect (UPDATE) 1 row, affected %d", rows)
		}

	case CountType:
		result, err := rep.DB.ExecContext(ctx, "UPDATE counters SET value = $1 WHERE name = $2", metricValue, metricName)
		if err != nil {
			logger.Log.Fatal("error while saving counter metric to the db", zap.Error(err))
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			logger.Log.Fatal("error while saving counter metric to the db", zap.Error(err))
		}
		if rows != 1 {
			log.Fatalf("expected to affect (UPDATE) 1 row, affected %d", rows)
		}
	}
	return nil
}

func (rep *Repository) Get(ctx context.Context, metricType string, metricName string) (models.Metrics, error) {
	var metric models.Metrics
	switch metricType {
	case GaugeType:
		var metric models.Metrics
		// to be checked
			err := rep.DB.QueryRowContext(ctx, "SELECT name, type, value FROM gauges WHERE name = $1",metricName).Scan(
			&metric.ID, &metric.MType, &metric.Value)
		switch {
		case err == sql.ErrNoRows:
			logger.Log.Fatal("no rows", zap.Error(err))
			return metric, err
		case err != nil:
			logger.Log.Fatal("error while obtaining metrics", zap.Error(err))
			return metric, err
		default:
			return metric, nil
		}
	case CountType:
		var metric models.Metrics
		err := rep.DB.QueryRowContext(ctx, "SELECT name, type, value FROM counters WHERE name = $1",metricName).Scan(
			&metric.ID, &metric.MType, &metric.Delta)
		switch {
		case err == sql.ErrNoRows:
			logger.Log.Fatal("no rows", zap.Error(err))
			return metric, err
		case err != nil:
			logger.Log.Fatal("error while obtaining metrics", zap.Error(err))
			return metric, err
		default:
			return metric, nil
		}
	}
	return metric, nil
}

// func (rep *Repository) Ping(ctx context.Context, metricType string, metricName string) (error) {
// 	if err := rep.DB.PingContext(ctx); err != nil {
//         logger.Log.Info("error", zap.Error(err))
// 		return err
//     } 
// 	return nil
// }

func (rep *Repository) Ping(w http.ResponseWriter, r *http.Request ) {
	ctx := r.Context()
	if err := rep.DB.PingContext(ctx); err != nil {
		fmt.Println(err)
        logger.Log.Info("error", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
    }
	rep.DB.Close()
	w.WriteHeader(http.StatusOK)
}

func (rep *Repository) GetAll(ctx context.Context) (map[string]any, error) {
	metrics := make(map[string]any, 33)
	rows, err := rep.DB.QueryContext(ctx, "SELECT name, value FROM gauges WHERE type = $1", GaugeType)
	if err != nil {
        return nil, err
    }
	defer rows.Close()
	for rows.Next() {
		var name 	string
		var value 	any
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
	rows, err = rep.DB.QueryContext(ctx, "SELECT name, value FROM counters WHERE type = $1", CountType)
	if err != nil {
        return nil, err
    }
	defer rows.Close()
	for rows.Next() {
		var name 	string
		var value 	any
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







