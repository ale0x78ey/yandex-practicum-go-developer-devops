package rest

import (
	"errors"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/config"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/middleware"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/service/server"
)

type Handler struct {
	Config *config.Config
	Server *server.Server
	Router *chi.Mux
}

func NewHandler(cfg *config.Config, srv *server.Server) (*Handler, error) {
	if cfg == nil {
		return nil, errors.New("invalid cfg value: nil")
	}
	if srv == nil {
		return nil, errors.New("invalid srv value: nil")
	}

	router := chi.NewRouter()

	h := &Handler{
		Config: cfg,
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

	h.Router.Get("/ping", h.getStorageStatus)

	return h, nil
}
