package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
)

const (
	shutdownTimeout = 5 * time.Second
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	serverConfig := restServerConfig{
		ShutdownTimeout: shutdownTimeout,
	}
	if err := env.Parse(&serverConfig); err != nil {
		log.Fatalf("Failed to parse REST server config options: %v", err)
	}

	server, err := newRestServer(serverConfig)
	if err != nil {
		log.Fatalf("Failed to create a server: %v", err)
	}

	if err := server.Run(ctx); err != nil {
		log.Fatalf("Failed to run a server: %v", err)
	}
}
