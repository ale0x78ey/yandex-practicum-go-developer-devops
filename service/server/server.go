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

func (s *Server) SaveMetric(
	ctx context.Context,
	metricType model.MetricType,
	metricName string,
	value string,
) error {
	switch metricType {
	case model.MetricTypeGauge:
		if _, err := model.GaugeFromString(value); err != nil {
			return err
		}
		return s.MetricStorer.SaveMetric(ctx, metricType, metricName, value)

	case model.MetricTypeCounter:
		newValue, err := model.CounterFromString(value)
		if err != nil {
			return err
		}

		oldValueString, err := s.MetricStorer.LoadMetric(ctx, metricType, metricName)
		if err != nil {
			return s.MetricStorer.SaveMetric(ctx, metricType, metricName, value)
		}

		oldValue, err := model.CounterFromString(oldValueString)
		if err != nil {
			return err
		}

		newValue += oldValue
		return s.MetricStorer.SaveMetric(ctx, metricType, metricName, newValue.String())

	default:
		return fmt.Errorf("unknown metricType: %v", metricType)
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