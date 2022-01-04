package rest

import (
	"net/http"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
	"github.com/go-chi/chi/v5"
)

func withMetricTypeValidator(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")
		if err := model.MetricType(metricType).Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusNotImplemented)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
