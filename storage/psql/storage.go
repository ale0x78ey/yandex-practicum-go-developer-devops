package psql

import (
	"context"
	"fmt"
	"sync"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
)

type MetricStorer struct {
	rw      sync.RWMutex
	metrics map[string]string
}

func NewMetricStorer() *MetricStorer {
	return &MetricStorer{
		metrics: make(map[string]string),
	}
}

func (s *MetricStorer) SaveMetric(
	ctx context.Context,
	metricType model.MetricType,
	metricName string,
	value string,
) error {
	s.rw.Lock()
	defer s.rw.Unlock()

	s.metrics[metricName] = value
	return nil
}

func (s *MetricStorer) LoadMetric(
	ctx context.Context,
	metricType model.MetricType,
	metricName string,
) (string, error) {
	s.rw.RLock()
	defer s.rw.RUnlock()

	if value, ok := s.metrics[metricName]; ok {
		return value, nil
	}

	return "", fmt.Errorf("%v not found", metricName)
}
