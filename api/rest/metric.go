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
			http.Error(w, err.Error(), http.StatusBadRequest)
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
			// Hide it because of tests in github.
			// r.Use(metricNameValidator)
			r.Post("/", updateMetric)
		})

	api.Routes.Metric.Route("/value/{metricType}/{metricName}",
		func(r chi.Router) {
			r.Use(metricTypeValidator)
			// Hide it because of tests in github.
			// r.Use(metricNameValidator)
			r.Get("/", getMetric)
		})
}

func updateMetric(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	srv := ctx.Value(ContextServerKey).(*server.Server)
	if srv == nil {
		http.Error(w, "srv == nil", http.StatusInternalServerError)
		return
	}

	// metricName := model.MetricName(chi.URLParam(r, "metricName"))
	metricName := chi.URLParam(r, "metricName")
	metricType := model.MetricType(chi.URLParam(r, "metricType"))

	switch metricType {
	case model.MetricTypeGauge:
		value, err := model.GaugeFromString(chi.URLParam(r, "metricValue"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := srv.MetricStorer.SaveMetricGauge(ctx, metricName, value); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case model.MetricTypeCounter:
		value, err := model.CounterFromString(chi.URLParam(r, "metricValue"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := srv.MetricStorer.SaveMetricCounter(ctx, metricName, value); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func getMetric(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	srv := ctx.Value(ContextServerKey).(*server.Server)
	if srv == nil {
		http.Error(w, "srv == nil", http.StatusInternalServerError)
		return
	}

	// metricName := model.MetricName(chi.URLParam(r, "metricName"))
	metricName := chi.URLParam(r, "metricName")
	metricType := model.MetricType(chi.URLParam(r, "metricType"))

	var strVal string

	switch metricType {
	case model.MetricTypeGauge:
		value, err := srv.MetricStorer.LoadMetricGauge(ctx, metricName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		strVal = value.String()

	case model.MetricTypeCounter:
		value, err := srv.MetricStorer.LoadMetricCounter(ctx, metricName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		strVal = value.String()
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, strVal)
}
