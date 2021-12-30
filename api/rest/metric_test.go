package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/service/server"
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
	srv := server.NewServer(metricStorer)
	if srv == nil {
		t.Fatalf("srv wasn't created")
	}

	r := chi.NewRouter()

	r.Route("/update/{metricType}/{metricName}/{metricValue}",
		func(r chi.Router) {
			r.Use(withServer(srv))
			r.Use(withMetricTypeValidator)
			r.Get("/", updateMetric)
		})

	server := httptest.NewServer(r)
	defer server.Close()

	// TODO: pass params from tests.
	metricStorer.EXPECT().SaveMetric(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode := doRequest(t, server, http.MethodGet, tt.path)
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
			path: "/get/gauge/testMetricName",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "Invalid MetricType",
			path: "/get/abcdef/testMetricName",
			want: want{
				code: http.StatusNotImplemented,
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	metricStorer := storagemock.NewMockMetricStorer(mockCtrl)
	srv := server.NewServer(metricStorer)
	if srv == nil {
		t.Fatalf("srv wasn't created")
	}

	r := chi.NewRouter()

	r.Route("/get/{metricType}/{metricName}",
		func(r chi.Router) {
			r.Use(withServer(srv))
			r.Use(withMetricTypeValidator)
			r.Get("/", getMetric)
		})

	server := httptest.NewServer(r)
	defer server.Close()

	// TODO: pass params from tests.
	metricStorer.EXPECT().LoadMetric(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return("", nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode := doRequest(t, server, http.MethodGet, tt.path)
			// TODO: Check not only statusCode.
			assert.Equal(t, tt.want.code, statusCode)
		})
	}
}
