package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/internal/common/testutils"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/internal/model"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/internal/storage"
	storagemock "github.com/ale0x78ey/yandex-practicum-go-developer-devops/internal/storage/mock"
)

func TestNewHandler(t *testing.T) {
	_, err := NewHandler(nil)
	assert.NotNil(t, err)

	h, err := NewHandler(&Server{})
	assert.Nil(t, err)
	assert.NotNil(t, h)
}

func newTestHandler(t *testing.T, metricStorage storage.MetricStorage) *Handler {
	config := Config{
		StoreInterval: 1 * time.Second,
	}
	srv, err := NewServer(config, metricStorage)
	require.NoError(t, err)

	h, err := NewHandler(srv)
	require.NoError(t, err)

	return h
}

func TestUpdateMetricWithURL(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name string
		path string
		want want
	}{
		{
			name: "Valid gauge metric1",
			path: "/update/gauge/metric1/123.45",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "Valid counter metric2",
			path: "/update/counter/metric2/123",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "Invalid MType",
			path: "/update/abcdef/metric3/123",
			want: want{
				code: http.StatusNotImplemented,
			},
		},
	}

	mockCtrl := gomock.NewController(t)

	metricStorage := storagemock.NewMockMetricStorage(mockCtrl)
	h := newTestHandler(t, metricStorage)
	server := httptest.NewServer(h.Router)
	defer server.Close()

	metric1 := model.MetricFromGauge("metric1", model.Gauge(123.45))
	metric2 := model.MetricFromCounter("metric2", model.Counter(123))

	gomock.InOrder(
		metricStorage.EXPECT().SaveMetric(gomock.Any(), metric1).Return(nil),
		metricStorage.EXPECT().IncrMetric(gomock.Any(), metric2).Return(nil),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, _ := testutils.DoRequest(t, server, http.MethodPost, tt.path, nil)
			assert.Equal(t, tt.want.code, statusCode)
		})
	}
}

func TestUpdateMetricWithBody(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name   string
		path   string
		metric model.Metric
		want   want
	}{
		{
			name:   "Valid gauge metric1",
			path:   "/update/",
			metric: model.MetricFromGauge("metric1", model.Gauge(123.45)),
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:   "Valid counter metric2",
			path:   "/update/",
			metric: model.MetricFromCounter("metric2", model.Counter(123)),
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:   "Invalid MType",
			path:   "/update/",
			metric: model.Metric{ID: "metric3", MType: model.MetricType("abcdef")},
			want: want{
				code: http.StatusNotImplemented,
			},
		},
	}

	mockCtrl := gomock.NewController(t)

	metricStorage := storagemock.NewMockMetricStorage(mockCtrl)
	h := newTestHandler(t, metricStorage)
	server := httptest.NewServer(h.Router)
	defer server.Close()

	metric1 := model.MetricFromGauge("metric1", model.Gauge(123.45))
	metric2 := model.MetricFromCounter("metric2", model.Counter(123))

	gomock.InOrder(
		metricStorage.EXPECT().SaveMetric(gomock.Any(), metric1).Return(nil),
		metricStorage.EXPECT().IncrMetric(gomock.Any(), metric2).Return(nil),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.metric)
			require.NoError(t, err)
			statusCode, _ := testutils.DoRequest(t, server, http.MethodPost, tt.path, &data)
			assert.Equal(t, tt.want.code, statusCode)
		})
	}
}

func TestGetMetricWithURL(t *testing.T) {
	type want struct {
		code int
		body string
	}
	tests := []struct {
		name string
		path string
		want want
	}{
		{
			name: "Valid gauge metric1",
			path: "/value/gauge/metric1",
			want: want{
				code: http.StatusOK,
				body: "123.45",
			},
		},
		{
			name: "Valid counter metric2",
			path: "/value/counter/metric2",
			want: want{
				code: http.StatusOK,
				body: "123",
			},
		},
		{
			name: "Invalid MType",
			path: "/value/abrakadabra/metric3",
			want: want{
				code: http.StatusNotImplemented,
				body: "unknown MetricType: abrakadabra\n",
			},
		},
	}

	mockCtrl := gomock.NewController(t)

	metricStorage := storagemock.NewMockMetricStorage(mockCtrl)
	h := newTestHandler(t, metricStorage)
	server := httptest.NewServer(h.Router)
	defer server.Close()

	metric1 := model.MetricFromGauge("metric1", model.Gauge(123.45))
	metric2 := model.MetricFromCounter("metric2", model.Counter(123))

	gomock.InOrder(
		metricStorage.EXPECT().LoadMetric(
			gomock.Any(),
			gomock.Any(),
		).Return(&metric1, nil),

		metricStorage.EXPECT().LoadMetric(
			gomock.Any(),
			gomock.Any(),
		).Return(&metric2, nil),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, body := testutils.DoRequest(t, server, http.MethodGet, tt.path, nil)
			assert.Equal(t, tt.want.code, statusCode)
			assert.Equal(t, tt.want.body, body)
		})
	}
}

func TestGetMetricWithBody(t *testing.T) {
	type want struct {
		code int
		body string
	}
	tests := []struct {
		name   string
		path   string
		metric model.Metric
		want   want
	}{
		{
			name:   "Valid gauge metric1",
			path:   "/value/",
			metric: model.Metric{ID: "metric1", MType: model.MetricTypeGauge},
			want: want{
				code: http.StatusOK,
				body: "{\"id\":\"metric1\",\"type\":\"gauge\",\"value\":123.45}",
			},
		},
		{
			name:   "Valid counter metric2",
			path:   "/value/",
			metric: model.Metric{ID: "metric2", MType: model.MetricTypeCounter},
			want: want{
				code: http.StatusOK,
				body: "{\"id\":\"metric2\",\"type\":\"counter\",\"delta\":123}",
			},
		},
		{
			name:   "Invalid MType",
			path:   "/value/",
			metric: model.Metric{ID: "metric3", MType: model.MetricType("abrakadabra")},
			want: want{
				code: http.StatusNotImplemented,
				body: "unknown MetricType: abrakadabra\n",
			},
		},
	}

	mockCtrl := gomock.NewController(t)

	metricStorage := storagemock.NewMockMetricStorage(mockCtrl)
	h := newTestHandler(t, metricStorage)
	server := httptest.NewServer(h.Router)
	defer server.Close()

	metric1 := model.MetricFromGauge("metric1", model.Gauge(123.45))
	metric2 := model.MetricFromCounter("metric2", model.Counter(123))

	gomock.InOrder(
		metricStorage.EXPECT().LoadMetric(
			gomock.Any(),
			gomock.Any(),
		).Return(&metric1, nil),

		metricStorage.EXPECT().LoadMetric(
			gomock.Any(),
			gomock.Any(),
		).Return(&metric2, nil),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.metric)
			require.NoError(t, err)
			statusCode, body := testutils.DoRequest(t, server, http.MethodPost, tt.path, &data)
			assert.Equal(t, tt.want.code, statusCode)
			assert.Equal(t, tt.want.body, body)
		})
	}
}
