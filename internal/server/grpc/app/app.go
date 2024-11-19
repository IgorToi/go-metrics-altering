package grpcapp

import (
	"database/sql"

	"honnef.co/go/tools/config"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	db *sql.DB,
	config *config.Config,
) *App {
	storage := postgres.NewRepository(db)

	grpcApp := grpcapp.New(config.Port, *storage, config.Ip)

	return &App{
		GRPCServer: grpcApp,
	}
}