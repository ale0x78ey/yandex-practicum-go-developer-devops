package server

import (
	"context"
	"fmt"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/storage"
)

type Server struct {
	storage.MetricStorer
}

func (s *Server) SaveMetric(ctx context.Context, metric *model.Metric) error {
	switch metric.Type {
	case model.MetricTypeGauge:
		if _, err := model.GaugeFromString(metric.StringValue); err != nil {
			return err
		}
		return s.MetricStorer.SaveMetric(ctx, metric)

	case model.MetricTypeCounter:
		newValue, err := model.CounterFromString(metric.StringValue)
		if err != nil {
			return err
		}

		oldValueString, err := s.MetricStorer.LoadMetric(ctx, metric.Type, metric.Name)
		if err != nil {
			return s.MetricStorer.SaveMetric(ctx, metric)
		}

		oldValue, err := model.CounterFromString(oldValueString)
		if err != nil {
			return err
		}

		metric.StringValue = (newValue + oldValue).String()
		return s.MetricStorer.SaveMetric(ctx, metric)

	default:
		return fmt.Errorf("unknown metricType: %v", metric.Type)
	}
}

func NewServer(metricStorer storage.MetricStorer) *Server {
	if metricStorer == nil {
		return nil
	}
	srv := &Server{
		MetricStorer: metricStorer,
	}
	return srv
}
