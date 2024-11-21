package grpcapp

import (
	"errors"
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	server "github.com/igortoigildin/go-metrics-altering/internal/server/grpc"
	adapter "github.com/igortoigildin/go-metrics-altering/pkg/interceptors/logging"
	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	GRPCServer *grpc.Server
	port       string
}

func New(
	port string,
	storage server.Storage,
) *App {

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	opts := []logging.Option{
		logging.WithLogOnEvents(logging.PayloadReceived),
	}

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		logging.UnaryServerInterceptor(adapter.InterceptorLogger(logger), opts...),
	))

	server.Register(gRPCServer, storage)

	return &App{
		GRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() error {
	if err := a.Run(); err != nil {
		logger.Log.Error("failed to run grpc app", zap.Error(err))
		return errors.New("failed to start grpc app")
	}
	return nil
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	l, err := net.Listen("tcp", a.port)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.Log.Info("grpc server is running:", zap.String("addr", l.Addr().String()))

	if err := a.GRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
