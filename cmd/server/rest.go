package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/api/rest"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/service/server"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/storage"
	storagefile "github.com/ale0x78ey/yandex-practicum-go-developer-devops/storage/file"
)

type restServerConfig struct {
	ShutdownTimeout time.Duration
	ServerAddress   string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	InitStore       bool          `env:"RESTORE" envDefault:"true"`
	StoreInterval   time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile       string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
}

type restServer struct {
	config        restServerConfig
	server        *server.Server
	httpServer    *http.Server
	metricStorage storage.MetricStorage
}

func newRestServer(config restServerConfig) (*restServer, error) {
	metricStorage, err := storagefile.NewMetricStorage(config.StoreFile, config.InitStore)
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
	if s.config.StoreInterval <= 0 {
		return fmt.Errorf("invalid non-positive StoreInterval=%v", s.config.StoreInterval)
	}

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
