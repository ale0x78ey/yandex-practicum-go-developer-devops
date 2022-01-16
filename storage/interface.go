//go:generate mockgen -source=interface.go -destination=./mock/storage.go -package=storagemock
package storage

import (
	"context"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
)

type MetricStorage interface {
	SaveMetric(ctx context.Context, metric model.Metric) error
	IncrMetric(ctx context.Context, metric model.Metric) error

	LoadMetric(
		ctx context.Context,
		metricType model.MetricType,
		metricName string,
	) (*model.Metric, error)

	// TODO: Add offset, limit.
	LoadMetricList(ctx context.Context) ([]model.Metric, error)

	Flush(ctx context.Context) error
}
