package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/storage"
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

type Server struct {
	storage.MetricStorage
	config *Config
}

func NewServer(config *Config, metricStorage storage.MetricStorage) (*Server, error) {
	if config == nil {
		return nil, errors.New("invalid config value: nil")
	}
	if metricStorage == nil {
		return nil, errors.New("invalid metricStorage value: nil")
	}

	srv := &Server{
		MetricStorage: metricStorage,
		config:        config,
	}

	return srv, nil
}

func (s *Server) ValidateMetricHash(metric model.Metric) (bool, error) {
	if s.config.Key != "" {
		hash, err := metric.ProcessHash(s.config.Key)
		if err != nil {
			return false, err
		}
		if hash != metric.Hash {
			return false, nil
		}
	}
	return true, nil
}

func (s *Server) PushMetric(ctx context.Context, metric model.Metric) error {
	switch metric.MType {
	case model.MetricTypeGauge:
		return s.MetricStorage.SaveMetric(ctx, metric)
	case model.MetricTypeCounter:
		return s.MetricStorage.IncrMetric(ctx, metric)
	default:
		return fmt.Errorf("unknown metricType: %v", metric.MType)
	}
}

func (s *Server) LoadMetric(
	ctx context.Context,
	metricType model.MetricType,
	metricName string,
) (*model.Metric, error) {
	metric, err := s.MetricStorage.LoadMetric(ctx, metricType, metricName)
	if err != nil {
		return nil, err
	}
	if s.config.Key != "" {
		hash, err := metric.ProcessHash(s.config.Key)
		if err != nil {
			return nil, err
		}
		metric.Hash = hash
	}
	return metric, nil
}

func (s *Server) LoadMetricList(ctx context.Context) ([]model.Metric, error) {
	metricList, err := s.MetricStorage.LoadMetricList(ctx)
	if err != nil {
		return nil, err
	}
	if s.config.Key != "" {
		for i := range metricList {
			metric := &metricList[i]
			hash, err := metric.ProcessHash(s.config.Key)
			if err != nil {
				return nil, err
			}
			metric.Hash = hash
		}
	}
	return metricList, nil
}

func (s *Server) Run(ctx context.Context) error {
	if s.config.StoreInterval <= 0 {
		return fmt.Errorf("invalid non-positive StoreInterval=%v", s.config.StoreInterval)
	}

	storeTicker := time.NewTicker(s.config.StoreInterval)
	defer storeTicker.Stop()

	for {
		select {
		case <-storeTicker.C:
			if err := s.Flush(ctx); err != nil {
				log.Fatalf("Failed to flush: %v", err)
			}
		case <-ctx.Done():
			ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
			defer cancel()
			if err := s.Flush(ctx); err != nil {
				log.Fatalf("Failed to flush: %v", err)
			}
			return nil
		}
	}
}
