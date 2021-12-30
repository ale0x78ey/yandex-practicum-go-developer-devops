//go:generate mockgen -source=interface.go -destination=./mock/storage.go -package=storagemock
package storage

import (
	"context"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
)

type MetricStorer interface {
	// TODO: explicitly use gauge and counter types
	// in value fields?
	SaveMetric(
		ctx context.Context,
		metricType model.MetricType,
		metricName, value string) error

	LoadMetric(ctx context.Context,
		metricType model.MetricType,
		metricName string) (string, error)
}
