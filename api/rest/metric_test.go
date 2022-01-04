package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
	storagemock "github.com/ale0x78ey/yandex-practicum-go-developer-devops/storage/mock"
)

func TestUpdateMetric(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name string
		path string
		want want
	}{
		{
			name: "Valid MetricType and value",
			path: "/update/gauge/testMetricName/123",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "Invalid MetricType",
			path: "/update/abcdef/testMetricName/123",
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

	// TODO: pass params from tests.
	metricStorer.EXPECT().SaveMetric(
		gomock.Any(),
		gomock.Any(),
	).Return(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode := doRequest(t, server, http.MethodPost, tt.path)
			assert.Equal(t, tt.want.code, statusCode)
		})
	}
}

func TestGetMetric(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name string
		path string
		want want
	}{
		{
			name: "Valid MetricType",
			path: "/value/gauge/testMetricName",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "Invalid MetricType",
			path: "/value/abcdef/testMetricName",
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

	// TODO: pass params from tests.
	metric := model.Metric{}
	metricStorer.EXPECT().LoadMetric(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(&metric, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode := doRequest(t, server, http.MethodGet, tt.path)
			// TODO: Check not only statusCode.
			assert.Equal(t, tt.want.code, statusCode)
		})
	}
}

func TestGetMetricList(t *testing.T) {
}
