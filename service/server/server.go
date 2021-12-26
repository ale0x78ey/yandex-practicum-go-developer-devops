package server

import (
	"github.com/go-chi/chi/v5"
)

type Routers struct {
	Root   *chi.Mux
	Metric *chi.Mux
}

type Server struct {
	Routers Routers
}

func NewServer() *Server {
	srv := &Server{}
	return srv
}
