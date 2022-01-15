package rest

import (
	"errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

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
	if router == nil {
		return nil, errors.New("router wasn't created")
	}

	h := &Handler{
		Server: srv,
		Router: router,
	}

	h.initMiddleware()
	h.initMetric()

	return h, nil
}

func (h *Handler) initMiddleware() {
	h.Router.Use(middleware.Recoverer)
}

func (h *Handler) initMetric() {
	h.Router.Route("/update/{metricType}/{metricName}/{metricValue}",
		func(r chi.Router) {
			r.Use(withMTypeValidator)
			r.Post("/", h.updateMetricWithURL)
		})

	h.Router.Post("/update", h.updateMetricWithBody)

	h.Router.Route("/value/{metricType}/{metricName}",
		func(r chi.Router) {
			r.Use(withMTypeValidator)
			r.Get("/", h.getMetricWithURL)
		})

	h.Router.Post("/value", h.getMetricWithBody)

	h.Router.Get("/", h.getMetricList)
}
