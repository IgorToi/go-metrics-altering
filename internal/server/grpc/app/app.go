package grpcapp

import (
	config "github.com/igortoigildin/go-metrics-altering/config/server"
	grpcapp "github.com/igortoigildin/go-metrics-altering/internal/server/grpc/app/grpc"
	"github.com/igortoigildin/go-metrics-altering/internal/storage"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(config *config.ConfigServer, storage storage.Storage) *App {

	grpcApp := grpcapp.New(config.FlagRunAddrGRPC, storage)

	return &App{
		GRPCServer: grpcApp,
	}
}
