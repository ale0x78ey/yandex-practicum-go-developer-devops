package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestMetricTypeValidator(t *testing.T) {
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
			path: "/smth/gauge",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "Invalid MetricType",
			path: "/smth/abrakadabra",
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}

	r := chi.NewRouter()

	r.Route("/smth/{metricType}", func(r chi.Router) {
		r.Use(metricTypeValidator)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	server := httptest.NewServer(r)
	defer server.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := doRequest(t, server, http.MethodGet, tt.path)
			assert.Equal(t, tt.want.code, response.StatusCode)
		})
	}
}

func TestMetricNameValidator(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name string
		path string
		want want
	}{
		{
			name: "Valid MetricName",
			path: "/smth/Alloc",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "Invalid MetricName",
			path: "/smth/abrakadabra",
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}

	r := chi.NewRouter()

	r.Route("/smth/{metricName}", func(r chi.Router) {
		r.Use(metricNameValidator)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	server := httptest.NewServer(r)
	defer server.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := doRequest(t, server, http.MethodGet, tt.path)
			assert.Equal(t, tt.want.code, response.StatusCode)
		})
	}
}

func TestUpdateMetric(t *testing.T) {
}

func TestGetMetric(t *testing.T) {
}
