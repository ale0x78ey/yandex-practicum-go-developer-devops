package file

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
)

type Config struct {
	InitStore bool   `env:"RESTORE"`
	StoreFile string `env:"STORE_FILE"`
}

type metricsMap map[model.MetricName]model.Metric
type metricsMapMap map[model.MetricType]metricsMap

type MetricStorage struct {
	sync.RWMutex

	config  Config
	metrics metricsMapMap
}

func NewMetricStorage(config Config) (*MetricStorage, error) {
	storage := &MetricStorage{
		config:  config,
		metrics: make(metricsMapMap),
	}

	if config.InitStore {
		file, err := os.OpenFile(config.StoreFile, os.O_RDONLY|os.O_CREATE, 0644)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		if err := json.NewDecoder(file).Decode(&storage.metrics); err != nil && err != io.EOF {
			return nil, err
		}
	}

	return storage, nil
}

func (s *MetricStorage) saveMetric(ctx context.Context, metric model.Metric) error {
	metrics, ok := s.metrics[metric.MType]
	if !ok {
		metrics = make(metricsMap)
		s.metrics[metric.MType] = metrics
	}

	metrics[metric.ID] = metric

	return nil
}

func (s *MetricStorage) SaveMetric(ctx context.Context, metric model.Metric) error {
	s.Lock()
	defer s.Unlock()

	return s.saveMetric(ctx, metric)
}

func (s *MetricStorage) IncrMetric(ctx context.Context, metric model.Metric) error {
	s.Lock()
	defer s.Unlock()

	m, err := s.loadMetric(ctx, metric)
	if err != nil {
		return err
	}

	if m != nil {
		if metric.Value != nil {
			*metric.Value += *m.Value
		}
		if metric.Delta != nil {
			*metric.Delta += *m.Delta
		}
	}

	return s.saveMetric(ctx, metric)
}

func (s *MetricStorage) loadMetric(
	ctx context.Context,
	metric model.Metric,
) (*model.Metric, error) {
	if metrics, ok := s.metrics[metric.MType]; ok {
		if m, ok := metrics[metric.ID]; ok {
			return &m, nil
		}
	}

	return nil, nil
}

func (s *MetricStorage) LoadMetric(
	ctx context.Context,
	metric model.Metric,
) (*model.Metric, error) {
	s.RLock()
	defer s.RUnlock()

	return s.loadMetric(ctx, metric)
}

func (s *MetricStorage) count() int {
	count := 0
	for _, metrics := range s.metrics {
		count += len(metrics)
	}
	return count
}

func (s *MetricStorage) SaveMetricList(ctx context.Context, metrics []model.Metric) error {
	for _, metric := range metrics {
		if err := s.SaveMetric(ctx, metric); err != nil {
			return err
		}
	}
	return nil
}

func (s *MetricStorage) IncrMetricList(ctx context.Context, metrics []model.Metric) error {
	for _, metric := range metrics {
		if err := s.IncrMetric(ctx, metric); err != nil {
			return err
		}
	}
	return nil
}

func (s *MetricStorage) LoadMetricList(ctx context.Context) ([]model.Metric, error) {
	s.RLock()
	defer s.RUnlock()

	metrics := make([]model.Metric, 0, s.count())

	for _, metricsByMetricType := range s.metrics {
		for _, m := range metricsByMetricType {
			metrics = append(metrics, m)
		}
	}

	return metrics, nil
}

func (s *MetricStorage) Flush(ctx context.Context) error {
	file, err := os.OpenFile(s.config.StoreFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	s.RLock()
	defer s.RUnlock()

	if err := json.NewEncoder(file).Encode(&s.metrics); err != nil {
		return nil
	}

	return nil
}

func (s *MetricStorage) Close() {
}

func (s *MetricStorage) Heartbeat(ctx context.Context) error {
	return nil
}
