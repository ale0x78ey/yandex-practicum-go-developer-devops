package rest

import (
	"net/http"

	"log"
	"path"

	"github.com/go-chi/chi/v5"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
)

func metricTypeValidator(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricType")
		if err := model.MetricType(metricType).Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("metricTypeValidator: %s", err)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func metricNameValidator(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricName")
		if err := model.MetricName(metricType).Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("metricNameValidator: %s", err)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (api *API) InitMetric() {
	api.Routes.Metric.Route("/update/{metricType}/{metricName}/{metricValue}",
		func(r chi.Router) {
			r.Use(metricTypeValidator)
			r.Use(metricNameValidator)
			r.Post("/", updateMetric)
		})

	api.Routes.Metric.Route("/value/{metricType}/{metricName}",
		func(r chi.Router) {
			r.Use(metricTypeValidator)
			r.Use(metricNameValidator)
			r.Get("/", getMetric)
		})
}

func updateMetric(w http.ResponseWriter, r *http.Request) {
	log.Printf("updateMetric: %v", path.Base(r.URL.Path))
}

func getMetric(w http.ResponseWriter, r *http.Request) {
}
