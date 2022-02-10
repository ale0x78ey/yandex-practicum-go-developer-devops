package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/pkg/testutils"
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
			name: "Valid Gauge MType",
			path: "/smth/gauge",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "Valid Counter MType",
			path: "/smth/counter",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "Invalid MType",
			path: "/smth/abrakadabra",
			want: want{
				code: http.StatusNotImplemented,
			},
		},
	}

	r := chi.NewRouter()

	r.Route("/smth/{metricType}", func(r chi.Router) {
		r.Use(MetricTypeValidator)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	server := httptest.NewServer(r)
	defer server.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, _ := testutils.DoRequest(t, server, http.MethodGet, tt.path, nil)
			assert.Equal(t, tt.want.code, statusCode)
		})
	}
}
