package storage

import (
	"context"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
)

type MetricStorer interface {
	SaveMetricGauge(
		ctx context.Context,
		metricName model.MetricName,
		value model.Gauge,
	) error

	SaveMetricCounter(
		ctx context.Context,
		metricName model.MetricName,
		value model.Counter,
	) error

	LoadMetricGauge(
		ctx context.Context,
		metricName model.MetricName,
	) (model.Gauge, error)

	LoadMetricCounter(
		ctx context.Context,
		metricName model.MetricName,
	) (model.Counter, error)
}
