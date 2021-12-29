package rest

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/service/server"
)

func metricTypeValidator(next http.Handler) http.Handler {
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

func metricNameValidator(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "metricName")
		if err := model.MetricName(metricType).Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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
			r.Post("/", updateMetric)
		})

	api.Routes.Metric.Route("/value/{metricType}/{metricName}",
		func(r chi.Router) {
			r.Use(metricTypeValidator)
			r.Get("/", getMetric)
		})
}

func updateMetric(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	srv := ctx.Value(ContextServerKey).(*server.Server)
	metricType := model.MetricType(chi.URLParam(r, "metricType"))
	metricName := chi.URLParam(r, "metricName")
	value := chi.URLParam(r, "metricValue")

	if err := srv.SaveMetric(ctx, metricType, metricName, value); err != nil {
		errCode := http.StatusInternalServerError
		http.Error(w, err.Error(), errCode)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getMetric(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	srv := ctx.Value(ContextServerKey).(*server.Server)
	metricType := model.MetricType(chi.URLParam(r, "metricType"))
	metricName := chi.URLParam(r, "metricName")

	value, err := srv.LoadMetric(ctx, metricType, metricName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, value)
}
