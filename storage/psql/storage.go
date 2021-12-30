package psql

import (
	"context"
	"fmt"
	"sync"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
)

type MetricStorer struct {
	rw sync.RWMutex
	// TODO: store metricType in Postgres.
	// It's not beautiful to use map[string]map[string]string or smth else.
	metrics map[string]string
}

func NewMetricStorer() *MetricStorer {
	return &MetricStorer{
		metrics: make(map[string]string),
	}
}

func (s *MetricStorer) SaveMetric(ctx context.Context, metric *model.Metric) error {
	s.rw.Lock()
	defer s.rw.Unlock()

	s.metrics[metric.Name.String()] = metric.StringValue
	return nil
}

func (s *MetricStorer) LoadMetric(
	ctx context.Context,
	metricType model.MetricType,
	metricName model.MetricName,
) (string, error) {
	s.rw.RLock()
	defer s.rw.RUnlock()

	if value, ok := s.metrics[metricName.String()]; ok {
		return value, nil
	}

	return "", fmt.Errorf("%v not found", metricName)
}

func (s *MetricStorer) LoadMetricList(ctx context.Context) ([]*model.Metric, error) {
	s.rw.RLock()
	defer s.rw.RUnlock()

	metrics := make([]*model.Metric, 0, len(s.metrics))
	for name, value := range s.metrics {
		m := &model.Metric{Name: model.MetricName(name), StringValue: value}
		metrics = append(metrics, m)
	}

	return metrics, nil
}
