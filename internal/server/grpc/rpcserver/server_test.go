package metricsgrpc

import (
	"context"
	"testing"

	config "github.com/igortoigildin/go-metrics-altering/config/server"
	"github.com/igortoigildin/go-metrics-altering/internal/storage"
	pb "github.com/igortoigildin/go-metrics-altering/pkg/metrics_v1"
	"github.com/stretchr/testify/assert"
)

func TestServerAPI_AddGaugeMetric(t *testing.T) {
	gaugeMetric := pb.GaugeMetric{
		Name:  "gauge",
		Value: float64(1),
	}
	cfg := config.ConfigServer{}
	st, _ := storage.New(&cfg)
	s := ServerAPI{
		Storage: st,
	}
	_, err := s.AddGaugeMetric(context.Background(), &pb.AddGaugeRequest{
		Metric: &gaugeMetric,
	})
	assert.NoError(t, err)
}

func TestServerAPI_AddCounterMetric(t *testing.T) {
	conterMetric := pb.CounterMetric{
		Name:  "counter",
		Value: int64(1),
	}
	cfg := config.ConfigServer{}
	st, _ := storage.New(&cfg)
	s := ServerAPI{
		Storage: st,
	}

	_, err := s.AddCounterMetric(context.Background(), &pb.AddCounterRequest{
		Metric: &conterMetric,
	})
	assert.NoError(t, err)
}
