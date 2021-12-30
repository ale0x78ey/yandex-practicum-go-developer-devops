package rest

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/service/server"
)

func (api *API) InitMetric() {
	api.Routes.Metric.Route("/update/{metricType}/{metricName}/{metricValue}",
		func(r chi.Router) {
			r.Use(withMetricTypeValidator)
			r.Post("/", updateMetric)
		})

	api.Routes.Metric.Route("/value/{metricType}/{metricName}",
		func(r chi.Router) {
			r.Use(withMetricTypeValidator)
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

	r.Body.Close()

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

	r.Body.Close()

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, value)
}
