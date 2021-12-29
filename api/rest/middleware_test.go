package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/service/server"
)

func TestWithMetricTypeValidator(t *testing.T) {
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
				code: http.StatusNotImplemented,
			},
		},
	}

	r := chi.NewRouter()

	r.Route("/smth/{metricType}", func(r chi.Router) {
		r.Use(withMetricTypeValidator)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	server := httptest.NewServer(r)
	defer server.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode := doRequest(t, server, http.MethodGet, tt.path)
			assert.Equal(t, tt.want.code, statusCode)
		})
	}
}

func TestWithMetricNameValidator(t *testing.T) {
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
		r.Use(withMetricNameValidator)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	server := httptest.NewServer(r)
	defer server.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode := doRequest(t, server, http.MethodGet, tt.path)
			assert.Equal(t, tt.want.code, statusCode)
		})
	}
}

func TestWithServer(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name string
		path string
		want want
	}{
		{
			name: "add server to context",
			path: "/srv",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "do not add server to context",
			path: "/nil",
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}

	r := chi.NewRouter()

	check := func(w http.ResponseWriter, r *http.Request) {
		if srv := r.Context().Value(ContextServerKey); srv != nil {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	r.Route("/srv", func(r chi.Router) {
		r.Use(withServer(&server.Server{}))
		r.Get("/", check)
	})

	r.Get("/nil", check)

	server := httptest.NewServer(r)
	defer server.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode := doRequest(t, server, http.MethodGet, tt.path)
			assert.Equal(t, tt.want.code, statusCode)
		})
	}
}
