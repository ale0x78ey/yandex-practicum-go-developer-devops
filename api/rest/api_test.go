package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/service/server"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/storage"
)

func TestNewHandler(t *testing.T) {
	_, err := NewHandler(nil)
	assert.NotNil(t, err)

	srv, err := NewHandler(&server.Server{})
	assert.Nil(t, err)
	assert.NotNil(t, srv)
}

func newTestHandler(t *testing.T, metricStorage storage.MetricStorage) *Handler {
	srv, err := server.NewServer(metricStorage)
	require.NoError(t, err)

	h, err := NewHandler(srv)
	require.NoError(t, err)

	return h
}
