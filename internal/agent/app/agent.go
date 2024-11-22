// Package agentcmd reads metrics from chan and calls relevant functions for sending it accordingly.
package agentcmd

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	config "github.com/igortoigildin/go-metrics-altering/config/agent"
	agent "github.com/igortoigildin/go-metrics-altering/internal/agent/sendMetrics"
	"github.com/igortoigildin/go-metrics-altering/internal/models"
	adapter "github.com/igortoigildin/go-metrics-altering/pkg/interceptors/logging"
	pb "github.com/igortoigildin/go-metrics-altering/pkg/metrics_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.uber.org/zap"
)

// SendMetrics reads metrics from metricsChan and sends it server address as defined by agent config.
func SendMetrics(metricsChan <-chan models.Metrics, cfg *config.ConfigAgent) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	opts := []logging.Option{
		logging.WithLogOnEvents(logging.PayloadSent),
	}

	conn, err := grpc.Dial(cfg.FlagRunPortGRPC, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			logging.UnaryClientInterceptor(adapter.InterceptorLogger(logger), opts...),
		))
	if err != nil {
		logger.Fatal("fail to dial grpc server")
	}
	defer conn.Close()

	m := pb.NewMetricsClient(conn)

	for metric := range metricsChan {

		switch metric.MType {
		case config.CountType:
			// http
			err := agent.SendURLCounter(cfg, int(*metric.Delta))
			if err != nil {
				logger.Error("unexpected sending url counter metric error:", zap.Error(err))
			}
			err = agent.SendJSONCounter(int(*metric.Delta), cfg)
			if err != nil {
				logger.Error("unexpected sending json counter metric error:", zap.Error(err))
			}

			// gRPC
			counterMetric := pb.CounterMetric{
				Name:  "counter",
				Value: *metric.Delta,
			}
			resp, err := m.AddCounterMetric(context.Background(), &pb.AddCounterRequest{
				Metric: &counterMetric,
			})
			if err != nil {
				logger.Error("error", zap.Error(err))
			}
			if resp.Error != "" {
				logger.Error(resp.Error)
			}

		case config.GaugeType:
			err := agent.SendURLGauge(cfg, *metric.Value, metric.ID)
			if err != nil {
				logger.Error("unexpected sending url gauge metric error:", zap.Error(err))
			}
			err = agent.SendJSONGauge(metric.ID, cfg, *metric.Value)
			if err != nil {
				logger.Error("unexpected sending json gauge metric error:", zap.Error(err))
			}

			// gRPC
			gaugeMetric := pb.GaugeMetric{
				Name:  "gauge",
				Value: *metric.Value,
			}
			resp, err := m.AddGaugeMetric(context.Background(), &pb.AddGaugeRequest{
				Metric: &gaugeMetric,
			})
			if err != nil {
				logger.Error("error", zap.Error(err))
			}
			if resp.Error != "" {
				logger.Error(resp.Error)
			}
		}
	}
}
