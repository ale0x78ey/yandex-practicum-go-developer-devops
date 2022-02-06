package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
)

const (
	defaultShutdownTimeout = 3 * time.Second
	defaultServerAddress   = "127.0.0.1:8080"
	defaultInitStore       = true
	defaultStoreInterval   = 300 * time.Second
	defaultStoreFile       = "/tmp/devops-metrics-db.json"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	config := restServerConfig{
		ShutdownTimeout: defaultShutdownTimeout,
	}

	flag.StringVar(&config.ServerAddress, "a", defaultServerAddress, "ADDRESS")
	flag.BoolVar(&config.InitStore, "r", defaultInitStore, "RESTORE")
	flag.DurationVar(&config.StoreInterval, "i", defaultStoreInterval, "STORE_INTERVAL")
	flag.StringVar(&config.StoreFile, "f", defaultStoreFile, "STORE_FILE")
	flag.Parse()

	if err := env.Parse(&config); err != nil {
		log.Fatalf("Failed to parse REST server config options: %v", err)
	}

	server, err := newRestServer(config)
	if err != nil {
		log.Fatalf("Failed to create a server: %v", err)
	}

	if err := server.Run(ctx); err != nil {
		log.Fatalf("Failed to run a server: %v", err)
	}
}
