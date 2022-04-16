package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/api/rest"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/config"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/service/server"
	storagefile "github.com/ale0x78ey/yandex-practicum-go-developer-devops/storage/file"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	cfg := config.LoadServerConfig()
	if cfg.StoreFile == nil {
		log.Fatalf("Missing config for storing in file")
	}
	metricStorage, err := storagefile.NewMetricStorage(
		cfg.StoreFile.StoreFile,
		cfg.StoreFile.InitStore,
	)
	if err != nil {
		log.Fatalf("Failed to create metric storage: %v", err)
	}

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
			log.Fatalf("Failed to run a server: %v", err)
		}
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Fatalf("HTTP Server failed: %v", err)
		}
	}()

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
