package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
	storagemock "github.com/ale0x78ey/yandex-practicum-go-developer-devops/storage/mock"
)

func TestUpdateMetricFromURL(t *testing.T) {
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
	defer mockCtrl.Finish()

	metricStorer := storagemock.NewMockMetricStorer(mockCtrl)
	h := newTestHandler(t, metricStorer)
	server := httptest.NewServer(h.Router)
	defer server.Close()

	metric1 := model.MetricFromGauge("metric1", model.Gauge(123.45))
	metric2 := model.MetricFromCounter("metric2", model.Counter(123))

	gomock.InOrder(
		metricStorer.EXPECT().SaveMetric(gomock.Any(), metric1).Return(nil),
		metricStorer.EXPECT().IncrMetric(gomock.Any(), metric2).Return(nil),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, _ := doRequest(t, server, http.MethodPost, tt.path, nil)
			assert.Equal(t, tt.want.code, statusCode)
		})
	}
}

func TestUpdateMetricFromBody(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name string
		path string
		metric model.Metric
		want want
	}{
		{
			name: "Valid gauge metric1",
			path: "/update",
			metric: model.MetricFromGauge("metric1", model.Gauge(123.45)),
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "Valid counter metric2",
			path: "/update",
			metric: model.MetricFromCounter("metric2", model.Counter(123)),
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "Invalid MType",
			path: "/update",
			metric: model.Metric{ID: "metric3", MType: model.MetricType("abcdef")},
			want: want{
				code: http.StatusNotImplemented,
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	metricStorer := storagemock.NewMockMetricStorer(mockCtrl)
	h := newTestHandler(t, metricStorer)
	server := httptest.NewServer(h.Router)
	defer server.Close()

	metric1 := model.MetricFromGauge("metric1", model.Gauge(123.45))
	metric2 := model.MetricFromCounter("metric2", model.Counter(123))

	gomock.InOrder(
		metricStorer.EXPECT().SaveMetric(gomock.Any(), metric1).Return(nil),
		metricStorer.EXPECT().IncrMetric(gomock.Any(), metric2).Return(nil),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.metric)
			require.NoError(t, err)
			statusCode, _ := doRequest(t, server, http.MethodPost, tt.path, &data)
			assert.Equal(t, tt.want.code, statusCode)
		})
	}
}

func TestGetMetric(t *testing.T) {
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
			path: "/value/abrakadabra/metric2",
			want: want{
				code: http.StatusNotImplemented,
				body: "unkown MType: abrakadabra\n",
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	metricStorer := storagemock.NewMockMetricStorer(mockCtrl)
	h := newTestHandler(t, metricStorer)
	server := httptest.NewServer(h.Router)
	defer server.Close()

	metric1 := model.MetricFromGauge("metric1", model.Gauge(123.45))
	metric2 := model.MetricFromCounter("metric2", model.Counter(123))

	gomock.InOrder(
		metricStorer.EXPECT().LoadMetric(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Return(&metric1, nil),

		metricStorer.EXPECT().LoadMetric(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Return(&metric2, nil),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, body := doRequest(t, server, http.MethodGet, tt.path, nil)
			assert.Equal(t, tt.want.code, statusCode)
			assert.Equal(t, tt.want.body, body)
		})
	}
}

func TestGetMetricList(t *testing.T) {
}
