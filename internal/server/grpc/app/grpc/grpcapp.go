package grpcapp

import (
	"errors"
	"fmt"
	"net"
	"net/netip"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/realip"
	config "github.com/igortoigildin/go-metrics-altering/config/server"
	server "github.com/igortoigildin/go-metrics-altering/internal/server/grpc/rpcserver"
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
	config *config.ConfigServer,
	storage server.Storage,
) *App {

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	opts1 := []logging.Option{
		logging.WithLogOnEvents(logging.PayloadReceived),
	}

	// Define list of trusted peers from which we accept forwarded-for and
	// real-ip headers.
	trustedPeers := []netip.Prefix{
		netip.MustParsePrefix(config.FlagTrustedSubnet),
	}

	// Define headers to look for in the incoming request.
	headers := []string{realip.XForwardedFor, realip.XRealIp}
	// Consider that there is one proxy in front,
	// so the real client ip will be rightmost - 1 in the csv list of X-Forwarded-For
	// Optionally you can specify TrustedProxies
	opts2 := []realip.Option{
		realip.WithTrustedPeers(trustedPeers),
		realip.WithHeaders(headers),
		realip.WithTrustedProxiesCount(1),
	}

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		logging.UnaryServerInterceptor(adapter.InterceptorLogger(logger), opts1...),
		realip.UnaryServerInterceptorOpts(opts2...),
	))

	server.Register(gRPCServer, storage)

	return &App{
		GRPCServer: gRPCServer,
		port:       config.FlagRunAddrGRPC,
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
