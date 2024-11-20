// Package agent accumulates, runtime metrics
// and sends it to predefined server every poll interval.

package agent

import (
	"context"

	config "github.com/igortoigildin/go-metrics-altering/config/agent"
	"github.com/igortoigildin/go-metrics-altering/internal/models"
	"github.com/igortoigildin/go-metrics-altering/pkg/logger"
	pb "github.com/igortoigildin/go-metrics-altering/pkg/metrics_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// SendMetrics reads metrics from metricsChan and sends it server address as defined by agent config.
func SendMetrics(metricsChan <-chan models.Metrics, cfg *config.ConfigAgent) {
	conn, err := grpc.Dial(cfg.FlagRunPortGRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Fatal("fail to dial grpc server")
	}
	defer conn.Close()

	m := pb.NewMetricsClient(conn)

	for metric := range metricsChan {

		switch metric.MType {
		case config.CountType:

			// http
			// err := sendURLCounter(cfg, int(*metric.Delta))
			// if err != nil {
			// 	logger.Log.Error("unexpected sending url counter metric error:", zap.Error(err))
			// }
			// err = SendJSONCounter(int(*metric.Delta), cfg)
			// if err != nil {
			// 	logger.Log.Error("unexpected sending json counter metric error:", zap.Error(err))
			// }

			// gRPC
			counterMetric := pb.CounterMetric{
				Name: "counter",
				Value: *metric.Delta,
			}
			resp, err := m.AddCounterMetric(context.Background(), &pb.AddCounterRequest{
				Metric: &counterMetric,
			})
			if err != nil {
				logger.Log.Error("error", zap.Error(err))
			}
			if resp.Error != "" {
				logger.Log.Error(resp.Error)
       		}

			

		case config.GaugeType:
			// err := SendURLGauge(cfg, *metric.Value, metric.ID)
			// if err != nil {
			// 	logger.Log.Error("unexpected sending url gauge metric error:", zap.Error(err))
			// }
			// err = SendJSONGauge(metric.ID, cfg, *metric.Value)
			// if err != nil {
			// 	logger.Log.Error("unexpected sending json gauge metric error:", zap.Error(err))
			// }

			// gRPC
			gaugeMetric := pb.GaugeMetric{
				Name: "gauge",
				Value: *metric.Value,
			}
			resp, err := m.AddGaugeMetric(context.Background(), &pb.AddGaugeRequest{
				Metric: &gaugeMetric,
			})
			if err != nil {
				logger.Log.Error("error", zap.Error(err))
			}
			if resp.Error != "" {
				logger.Log.Error(resp.Error)
       		}
		}
	}
}
