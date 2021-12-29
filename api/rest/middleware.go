package rest

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/pkg"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/service/server"
)

const (
	ContextServerKey = pkg.ContextKey("Server")
)

func (api *API) InitMiddleware() {
	api.Routes.Root.Use(middleware.Recoverer)
	api.Routes.Root.Use(withServer(api.srv))
}

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

func withMetricNameValidator(next http.Handler) http.Handler {
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

func withServer(srv *server.Server) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), ContextServerKey, srv)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
