package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/storage"
)

type Server struct {
	storage.MetricStorer
}

func NewServer(metricStorer storage.MetricStorer) (*Server, error) {
	if metricStorer == nil {
		return nil, errors.New("invalid metricStorer value: nil")
	}

	srv := &Server{
		MetricStorer: metricStorer,
	}

	return srv, nil
}

func (s *Server) PushMetric(ctx context.Context, metric model.Metric) error {
	switch metric.MType {
	case model.MetricTypeGauge:
		return s.MetricStorer.SaveMetric(ctx, metric)
	case model.MetricTypeCounter:
		return s.MetricStorer.IncrMetric(ctx, metric)
	default:
		return fmt.Errorf("unknown metricType: %v", metric.MType)
	}
}

func (s *Server) Flush(ctx context.Context) error {
	return nil
}
