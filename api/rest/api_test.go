package rest

import (
	"testing"
	"time"

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
	srvConfig := server.Config{
		StoreInterval: 1 * time.Second,
	}
	srv, err := server.NewServer(srvConfig, metricStorage)
	require.NoError(t, err)

	hConfig := &config.Config{}
	h, err := NewHandler(hConfig, srv)
	require.NoError(t, err)

	return h
}
