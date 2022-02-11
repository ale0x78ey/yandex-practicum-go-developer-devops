package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/config"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/service/server"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/storage"
)

func TestNewHandler(t *testing.T) {
	cfg := config.LoadServerConfig()
	_, err := NewHandler(cfg, nil)
	assert.NotNil(t, err)

	srv := &server.Server{}
	_, err = NewHandler(nil, srv)
	assert.NotNil(t, err)

	h, err := NewHandler(cfg, srv)
	assert.Nil(t, err)
	assert.NotNil(t, h)
}

func newTestHandler(t *testing.T, metricStorage storage.MetricStorage) *Handler {
	srv, err := server.NewServer(&server.Config{}, metricStorage)
	require.NoError(t, err)

	h, err := NewHandler(&config.Config{}, srv)
	require.NoError(t, err)

	return h
}
