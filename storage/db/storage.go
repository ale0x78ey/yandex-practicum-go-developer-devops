package db

import (
	"context"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
)

type Config struct {
	DSN string `env:"DATABASE_DSN"`
}

type MetricStorage struct {
	config Config
}

func NewMetricStorage(config Config) (*MetricStorage, error) {
	storage := &MetricStorage{
		config: config,
	}

	return storage, nil
}

func (s *MetricStorage) SaveMetric(ctx context.Context, metric model.Metric) error {
	return nil
}

func (s *MetricStorage) IncrMetric(ctx context.Context, metric model.Metric) error {
	return nil
}

func (s *MetricStorage) LoadMetric(
	ctx context.Context,
	metricType model.MetricType,
	metricName string,
) (*model.Metric, error) {
	return nil, nil
}

func (s *MetricStorage) LoadMetricList(ctx context.Context) ([]model.Metric, error) {
	metrics := make([]model.Metric, 0)
	return metrics, nil
}

func (s *MetricStorage) Flush(ctx context.Context) error {
	return nil
}
