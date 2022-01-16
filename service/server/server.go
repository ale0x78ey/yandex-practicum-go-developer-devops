package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/storage"
)

type Server struct {
	storage.MetricStorage
}

func NewServer(metricStorage storage.MetricStorage) (*Server, error) {
	if metricStorage == nil {
		return nil, errors.New("invalid metricStorage value: nil")
	}

	srv := &Server{
		MetricStorage: metricStorage,
	}

	return srv, nil
}

func (s *Server) PushMetric(ctx context.Context, metric model.Metric) error {
	switch metric.MType {
	case model.MetricTypeGauge:
		return s.MetricStorage.SaveMetric(ctx, metric)
	case model.MetricTypeCounter:
		return s.MetricStorage.IncrMetric(ctx, metric)
	default:
		return fmt.Errorf("unknown metricType: %v", metric.MType)
	}
}
