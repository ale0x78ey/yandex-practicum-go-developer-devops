package server

import (
	"testing"

	"github.com/stretchr/testify/assert"

	storagemock "github.com/ale0x78ey/yandex-practicum-go-developer-devops/storage/mock"
)

func TestNewServer(t *testing.T) {
	_, err := NewServer(nil)
	assert.NotNil(t, err)

	metricStorer := storagemock.NewMockMetricStorer(nil)
	srv, err := NewServer(metricStorer)
	assert.Nil(t, err)
	assert.NotNil(t, srv)
}

func TestServer_SaveMetric(t *testing.T) {
}

func TestServer_LoadMetric(t *testing.T) {
}
