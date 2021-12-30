package server

import (
	"testing"

	"github.com/stretchr/testify/assert"

	storagemock "github.com/ale0x78ey/yandex-practicum-go-developer-devops/storage/mock"
)

func TestServer_SaveMetric(t *testing.T) {
}

func TestServer_LoadMetric(t *testing.T) {
}

func TestNewServer(t *testing.T) {
	metricStorer := storagemock.NewMockMetricStorer(nil)
	assert.NotNil(t, NewServer(metricStorer))
}
