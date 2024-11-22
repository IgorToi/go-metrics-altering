package metricsgrpc

import (
	"context"

	"github.com/igortoigildin/go-metrics-altering/internal/models"
	metrics "github.com/igortoigildin/go-metrics-altering/pkg/metrics_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	gauge   = "gauge"
	counter = "counter"
)

//go:generate go run github.com/vektra/mockery/v2@v2.45.0 --name=Storage
type Storage interface {
	Update(ctx context.Context, metricType string, metricName string, metricValue any) error
	Get(ctx context.Context, metricType string, metricName string) (models.Metrics, error)
	GetAll(ctx context.Context) (map[string]any, error)
	Ping(ctx context.Context) error
}

type ServerAPI struct {
	metrics.UnimplementedMetricsServer
	Storage Storage
}

func Register(gRPC *grpc.Server, storage Storage) {
	metrics.RegisterMetricsServer(gRPC, &ServerAPI{Storage: storage})
}

func (s *ServerAPI) AddGaugeMetric(ctx context.Context, req *metrics.AddGaugeRequest) (*metrics.AddGaugeResponse, error) {
	err := s.Storage.Update(ctx, gauge, req.Metric.Name, req.Metric.Value)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return nil, nil
}

func (s *ServerAPI) AddCounterMetric(ctx context.Context, req *metrics.AddCounterRequest) (*metrics.AddCounterResponse, error) {
	err := s.Storage.Update(ctx, counter, req.Metric.Name, req.Metric.Value)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return nil, nil
}
