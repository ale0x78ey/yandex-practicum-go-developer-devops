package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/internal/model"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/internal/storage"
)

const (
	DefaultShutdownTimeout = 3 * time.Second
	DefaultStoreInterval   = 300 * time.Second
)

type Config struct {
	ShutdownTimeout time.Duration
	StoreInterval   time.Duration `env:"STORE_INTERVAL"`
	Key             string        `env:"KEY"`
}

func (c Config) Validate() error {
	if c.StoreInterval <= 0 {
		return fmt.Errorf("invalid non-positive StoreInterval=%v", c.StoreInterval)
	}
	return nil
}

type Server struct {
	storage.MetricStorage
	config Config
}

func NewServer(config Config, metricStorage storage.MetricStorage) (*Server, error) {
	if metricStorage == nil {
		return nil, errors.New("invalid metricStorage value: nil")
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	srv := &Server{
		MetricStorage: metricStorage,
		config:        config,
	}

	return srv, nil
}

func (s *Server) PushMetric(ctx context.Context, metric model.Metric) error {
	switch metric.MType {
	case model.MetricTypeGauge:
		return s.MetricStorage.SaveMetric(ctx, metric)
	case model.MetricTypeCounter:
		return s.MetricStorage.IncrMetric(ctx, metric)
	}
	return nil
}

func (s *Server) PushMetricList(ctx context.Context, metrics []model.Metric) error {
	gaugeMetrics := make([]model.Metric, 0, len(metrics))
	counterMetrics := make([]model.Metric, 0, len(metrics))

	for _, metric := range metrics {
		switch metric.MType {
		case model.MetricTypeGauge:
			gaugeMetrics = append(gaugeMetrics, metric)
		case model.MetricTypeCounter:
			counterMetrics = append(counterMetrics, metric)
		}
	}

	if err := s.MetricStorage.SaveMetricList(ctx, gaugeMetrics); err != nil {
		return err
	}

	if err := s.MetricStorage.IncrMetricList(ctx, counterMetrics); err != nil {
		return err
	}

	return nil
}

func (s *Server) LoadMetric(ctx context.Context, metric model.Metric) (*model.Metric, error) {
	m, err := s.MetricStorage.LoadMetric(ctx, metric)
	if err != nil {
		return nil, err
	}

	if m == nil {
		return m, nil
	}

	if err := m.UpdateHash(s.config.Key); err != nil {
		return nil, err
	}

	return m, nil
}

func (s *Server) LoadMetricList(ctx context.Context) ([]model.Metric, error) {
	metricList, err := s.MetricStorage.LoadMetricList(ctx)
	if err != nil {
		return nil, err
	}
	if s.config.Key != "" {
		for i := range metricList {
			metric := &metricList[i]
			if err := metric.UpdateHash(s.config.Key); err != nil {
				return nil, err
			}
		}
	}
	return metricList, nil
}

func (s *Server) ValidateHash(metric model.Metric) (bool, error) {
	if s.config.Key == "" {
		return true, nil
	}

	v, err := metric.ValidateHash(s.config.Key)
	if err != nil {
		return false, err
	}
	return v, nil
}

func (s *Server) Run(ctx context.Context) error {
	storeTicker := time.NewTicker(s.config.StoreInterval)
	defer storeTicker.Stop()

	for {
		select {
		case <-storeTicker.C:
			if err := s.Flush(ctx); err != nil {
				return fmt.Errorf("failed to flush: %w", err)
			}
		case <-ctx.Done():
			ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
			defer cancel()

			if err := s.Flush(ctx); err != nil {
				return fmt.Errorf("failed to flush: %w", err)
			}

			return nil
		}
	}
}
