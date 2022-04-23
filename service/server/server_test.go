package server

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
	storagemock "github.com/ale0x78ey/yandex-practicum-go-developer-devops/storage/mock"
)

func TestNewServer(t *testing.T) {
	cfg := Config{
		StoreInterval: 1*time.Second,
	}
	_, err := NewServer(cfg, nil)
	assert.NotNil(t, err)

	metricStorage := storagemock.NewMockMetricStorage(nil)
	srv, err := NewServer(cfg, metricStorage)
	assert.Nil(t, err)
	assert.NotNil(t, srv)
}

func TestServer_PushMetric(t *testing.T) {
	tests := []struct {
		name    string
		metric  model.Metric
		wantErr bool
	}{
		{
			name:    "gauge metric1",
			metric:  model.MetricFromGauge("metric1", model.Gauge(0)),
			wantErr: false,
		},
		{
			name:    "counter metric2",
			metric:  model.MetricFromCounter("metric2", model.Counter(0)),
			wantErr: false,
		},
		{
			name: "invalid metricType for metric3",
			metric: model.Metric{
				ID:    "metric3",
				MType: model.MetricType("abrakadabra"),
			},
			wantErr: true,
		},
	}

	mockCtrl := gomock.NewController(t)

	cfg := Config{
		StoreInterval: 1*time.Second,
	}
	metricStorage := storagemock.NewMockMetricStorage(mockCtrl)
	srv, err := NewServer(cfg, metricStorage)
	assert.Nil(t, err)

	gomock.InOrder(
		metricStorage.EXPECT().SaveMetric(gomock.Any(), gomock.Any()).Return(nil),
		metricStorage.EXPECT().IncrMetric(gomock.Any(), gomock.Any()).Return(nil),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := srv.PushMetric(context.Background(), tt.metric)
			if !tt.wantErr {
				require.NoError(t, err)
				return
			}
			assert.Error(t, err)
		})
	}
}
