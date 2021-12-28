package server

import (
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/storage"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/storage/psql"
)

type Server struct {
	MetricStorer storage.MetricStorer
}

func NewServer() *Server {
	metricStorer := psql.NewMetricStorer()
	srv := &Server{
		MetricStorer: metricStorer,
	}
	return srv
}
