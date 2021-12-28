package psql

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
)

type MetricStorer struct {
	rw       sync.RWMutex
	gauges   map[string]model.Gauge
	counters map[string]model.Counter
}

func NewMetricStorer() *MetricStorer {
	return &MetricStorer{
		gauges:   make(map[string]model.Gauge),
		counters: make(map[string]model.Counter),
	}
}

func (s *MetricStorer) SaveMetricGauge(
	ctx context.Context,
	// metricName model.MetricName,
	metricName string,
	value model.Gauge,
) error {
	s.rw.Lock()
	defer s.rw.Unlock()

	// s.gauges[metricName.String()] = value
	s.gauges[metricName] = value
	return nil
}

func (s *MetricStorer) SaveMetricCounter(
	ctx context.Context,
	// metricName model.MetricName,
	metricName string,
	value model.Counter,
) error {
	s.rw.Lock()
	defer s.rw.Unlock()

	// s.counters[metricName.String()] = value
	s.counters[metricName] = value
	return nil
}

func (s *MetricStorer) LoadMetricGauge(
	ctx context.Context,
	// metricName model.MetricName,
	metricName string,
) (model.Gauge, error) {
	s.rw.RLock()
	defer s.rw.RUnlock()

	var value model.Gauge
	// if value, ok := s.gauges[metricName.String()]; ok {
	if value, ok := s.gauges[metricName]; ok {
		return value, nil
	}

	return value, errors.New(fmt.Sprintf("error to load %v", metricName))
}

func (s *MetricStorer) LoadMetricCounter(
	ctx context.Context,
	// metricName model.MetricName,
	metricName string,
) (model.Counter, error) {
	s.rw.RLock()
	defer s.rw.RUnlock()

	var value model.Counter
	// if value, ok := s.counters[metricName.String()]; ok {
	if value, ok := s.counters[metricName]; ok {
		return value, nil
	}

	return value, errors.New(fmt.Sprintf("error to load %v", metricName))
}
