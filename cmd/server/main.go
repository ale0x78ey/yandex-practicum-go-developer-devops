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
	shutdownTimeout = 3 * time.Second
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
		ShutdownTimeout: shutdownTimeout,
	}
	if err := env.Parse(&config); err != nil {
		log.Fatalf("Failed to parse REST server config options: %v", err)
	}

	log.Printf("!!!config1!!! restore=%v", config.InitStore)

	flag.StringVar(&config.ServerAddress, "a", config.ServerAddress, "ADDRESS")
	flag.BoolVar(&config.InitStore, "r", config.InitStore, "RESTORE")
	flag.DurationVar(&config.StoreInterval, "i", config.StoreInterval, "STORE_INTERVAL")
	flag.StringVar(&config.StoreFile, "f", config.StoreFile, "STORE_FILE")
	flag.Parse()

	log.Printf("!!!config2!!! restore=%v", config.InitStore)

	server, err := newRestServer(config)
	if err != nil {
		log.Fatalf("Failed to create a server: %v", err)
	}

	if err := server.Run(ctx); err != nil {
		log.Fatalf("Failed to run a server: %v", err)
	}
}
