package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/caarlos0/env/v6"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/api/rest"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/service/server"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/storage"
	storagefile "github.com/ale0x78ey/yandex-practicum-go-developer-devops/storage/file"
)

type restServerConfig struct {
	ShutdownTimeout time.Duration
	ServerAddress   string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	StoreInterval   time.Duration `env:"STORE_INTERVAL" envDefault:"30s"`
}

type restServer struct {
	config        restServerConfig
	server        *server.Server
	httpServer    *http.Server
	metricStorage storage.MetricStorage
}

func newRestServer(config restServerConfig) (*restServer, error) {
	var metricStorageConfig storagefile.Config
	if err := env.Parse(&metricStorageConfig); err != nil {
		log.Fatalf("Failed to parse metric storage config options: %v", err)
	}

	metricStorage, err := storagefile.NewMetricStorage(metricStorageConfig)
	if err != nil {
		log.Fatalf("Failed to create metric storage: %v", err)
	}

	s, err := server.NewServer(metricStorage)
	if err != nil {
		return nil, err
	}

	h, err := rest.NewHandler(s)
	if err != nil {
		return nil, err
	}

	server := &restServer{
		config: config,
		server: s,
		httpServer: &http.Server{
			Addr:    config.ServerAddress,
			Handler: h.Router,
		},
		metricStorage: metricStorage,
	}

	return server, nil
}

func (s restServer) Run(ctx context.Context) error {
	go func() {
		storeTicker := time.NewTicker(s.config.StoreInterval)
		defer storeTicker.Stop()

		for {
			select {
			case <-storeTicker.C:
				if err := s.metricStorage.Flush(ctx); err != nil {
					log.Fatalf("Failed to flush: %v", err)
				}
			case <-ctx.Done():
				ctx, cancel := context.WithTimeout(
					context.Background(), s.config.ShutdownTimeout)
				defer cancel()
				if err := s.httpServer.Shutdown(ctx); err != nil {
					log.Fatalf("HTTP Server failed: %v", err)
				}
				if err := s.metricStorage.Flush(ctx); err != nil {
					log.Fatalf("Failed to flush: %v", err)
				}
				break
			}
		}
	}()

	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}