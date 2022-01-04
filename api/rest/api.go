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
			r.Use(withMetricTypeValidator)
			r.Post("/", h.updateMetric)
		})

	h.Router.Route("/value/{metricType}/{metricName}",
		func(r chi.Router) {
			r.Use(withMetricTypeValidator)
			r.Get("/", h.getMetric)
		})

	h.Router.Get("/", h.getMetricList)
}
