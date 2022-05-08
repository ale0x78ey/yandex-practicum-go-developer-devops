package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/api/rest"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/config"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/service/server"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/storage"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/storage/db"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/storage/file"
)

func NewMetricStorager(ctx context.Context, cfg *config.Config) (storage.MetricStorage, error) {
	if cfg == nil {
		return nil, errors.New("invalid cfg value: nil")
	}

	if cfg.DB != nil && cfg.DB.DSN != "" {
		return db.NewMetricStorage(ctx, *cfg.DB)
	}

	if cfg.StoreFile != nil && cfg.StoreFile.StoreFile != "" {
		return file.NewMetricStorage(*cfg.StoreFile)
	}

	return nil, errors.New("missing config options for a metric storager")
}

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	cfg := config.LoadServerConfig()
	metricStorage, err := NewMetricStorager(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to create a metric storage: %v", err)
	}
	defer metricStorage.Close()

	if cfg.Server == nil {
		log.Fatalf("Missing config for server")
	}
	s, err := server.NewServer(*cfg.Server, metricStorage)
	if err != nil {
		log.Fatalf("Failed to create a server: %v", err)
	}

	h, err := rest.NewHandler(cfg, s)
	if err != nil {
		log.Fatalf("Failed to create a handler: %v", err)
	}

	httpServer := &http.Server{
		Addr:    cfg.Http.ServerAddress,
		Handler: h.Router,
	}

	go func() {
		if err := s.Run(ctx); err != nil {
			log.Fatalf("Failed in a running server: %v", err)
		}
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Fatalf("HTTP Server failed: %v", err)
		}
	}()

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
