package storagefile

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
)

type MetricStorage struct {
	sync.RWMutex

	storeFile string
	metrics   map[string]model.Metric
}

func NewMetricStorage(storeFile string, initStore bool) (*MetricStorage, error) {
	storage := &MetricStorage{
		storeFile: storeFile,
		metrics:   make(map[string]model.Metric),
	}

	if initStore {
		file, err := os.OpenFile(storeFile, os.O_RDONLY|os.O_CREATE, 0666)
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

func (s *MetricStorage) SaveMetric(ctx context.Context, metric model.Metric) error {
	s.Lock()
	defer s.Unlock()

	s.metrics[metric.ID] = metric
	return nil
}

func (s *MetricStorage) IncrMetric(ctx context.Context, metric model.Metric) error {
	s.Lock()
	defer s.Unlock()

	if oldMetric, ok := s.metrics[metric.ID]; ok {
		if oldMetric.MType != metric.MType {
			return fmt.Errorf("different metric types for metric %s", metric.ID)
		}
		if metric.Value != nil {
			*metric.Value += *oldMetric.Value
		}
		if metric.Delta != nil {
			*metric.Delta += *oldMetric.Delta
		}
	}

	s.metrics[metric.ID] = metric

	return nil
}

func (s *MetricStorage) LoadMetric(
	ctx context.Context,
	metricType model.MetricType,
	metricName string,
) (*model.Metric, error) {
	s.RLock()
	defer s.RUnlock()

	if metric, ok := s.metrics[metricName]; ok && metric.MType == metricType {
		return &metric, nil
	}

	return nil, nil
}

func (s *MetricStorage) LoadMetricList(ctx context.Context) ([]model.Metric, error) {
	s.RLock()
	defer s.RUnlock()

	metrics := make([]model.Metric, 0, len(s.metrics))
	for _, metric := range s.metrics {
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

func (s *MetricStorage) Flush(ctx context.Context) error {
	file, err := os.OpenFile(s.storeFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
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
