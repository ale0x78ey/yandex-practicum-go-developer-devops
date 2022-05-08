//go:generate mockgen -source=interface.go -destination=./mock/storage.go -package=storagemock
package storage

import (
	"context"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
)

type MetricStorage interface {
	SaveMetric(ctx context.Context, metric model.Metric) error
	IncrMetric(ctx context.Context, metric model.Metric) error
	LoadMetric(ctx context.Context, metric model.Metric) (*model.Metric, error)

	SaveMetricList(ctx context.Context, metrics []model.Metric) error
	IncrMetricList(ctx context.Context, metrics []model.Metric) error
	LoadMetricList(ctx context.Context) ([]model.Metric, error)

	Flush(ctx context.Context) error
	Close()

	Heartbeat(ctx context.Context) error
}
