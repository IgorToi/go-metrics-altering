package migrations

import (
	"database/sql"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
	"go.uber.org/zap"
)

func New(migrationsPath string, db *sql.DB) error {
	if migrationsPath == "" {
		return errors.New("missing migrations path")
	}
	instance, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Log.Fatal("migration error", zap.Error(err))
	}

	migrator, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", instance)
	if err != nil {
		logger.Log.Fatal("migration error", zap.Error(err))
	}
	if err = migrator.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Log.Fatal("migration error", zap.Error(err))
	}
	return nil
}
