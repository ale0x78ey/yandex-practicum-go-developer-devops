package psql

import (
	"context"
	"fmt"
	"sync"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
)

type MetricStorer struct {
	sync.RWMutex

	metrics map[string]model.Metric
}

func NewMetricStorer() *MetricStorer {
	return &MetricStorer{
		metrics: make(map[string]model.Metric),
	}
}

func (s *MetricStorer) SaveMetric(ctx context.Context, metric model.Metric) error {
	s.Lock()
	defer s.Unlock()

	s.metrics[metric.Name] = metric
	return nil
}

func (s *MetricStorer) IncrMetric(ctx context.Context, metric model.Metric) error {
	s.Lock()
	defer s.Unlock()

	if oldMetric, ok := s.metrics[metric.Name]; ok {
		if oldMetric.Type != metric.Type {
			return fmt.Errorf("different metric types for metric %s", metric.Name)
		}
		metric.GaugeValue += oldMetric.GaugeValue
		metric.CounterValue += oldMetric.CounterValue
	}

	s.metrics[metric.Name] = metric

	return nil
}

func (s *MetricStorer) LoadMetric(
	ctx context.Context,
	metricType model.MetricType,
	metricName string,
) (*model.Metric, error) {
	s.RLock()
	defer s.RUnlock()

	if metric, ok := s.metrics[metricName]; ok && metric.Type == metricType {
		return &metric, nil
	}

	return nil, nil
}

func (s *MetricStorer) LoadMetricList(ctx context.Context) ([]model.Metric, error) {
	s.RLock()
	defer s.RUnlock()

	metrics := make([]model.Metric, 0, len(s.metrics))
	for _, metric := range s.metrics {
		metrics = append(metrics, metric)
	}

	return metrics, nil
}
