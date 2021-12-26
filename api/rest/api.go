package rest

import (
	"github.com/go-chi/chi/v5"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/service/server"
)

type Routes struct {
	Root   *chi.Mux
	Metric *chi.Mux
}

type API struct {
	srv    *server.Server
	Routes *Routes
}

func Init(srv *server.Server) *API {
	if srv == nil {
		return nil
	}

	api := &API{
		srv:    srv,
		Routes: &Routes{},
	}

	api.Routes.Metric = chi.NewRouter()
	if api.Routes.Metric == nil {
		return nil
	}

	// It's better to mount the API with own prefixes
	// but these are the requirements of the task.
	api.Routes.Root = api.Routes.Metric

	return api
}
