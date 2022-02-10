package rest

import (
	"errors"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/middleware"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/service/server"
)

type Handler struct {
	Server *server.Server
	Router *chi.Mux
}

func NewHandler(srv *server.Server) (*Handler, error) {
	if srv == nil {
		return nil, errors.New("invalid srv value: nil")
	}

	router := chi.NewRouter()

	h := &Handler{
		Server: srv,
		Router: router,
	}

	h.Router.Use(chimw.Recoverer)
	h.Router.Use(middleware.GzipDecoder())
	h.Router.Use(middleware.GzipEncoder())

	h.Router.Route("/update/{metricType}/{metricName}/{metricValue}", func(r chi.Router) {
		r.Use(middleware.MetricTypeValidator)
		r.Post("/", h.updateMetricWithURL)
	})

	h.Router.Post("/update/", h.updateMetricWithBody)

	h.Router.Route("/value/{metricType}/{metricName}", func(r chi.Router) {
		r.Use(middleware.MetricTypeValidator)
		r.Get("/", h.getMetricWithURL)
	})

	h.Router.Post("/value/", h.getMetricWithBody)

	h.Router.Get("/", h.getMetricList)

	return h, nil
}
