package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/api/rest"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/service/server"
)

const (
	shutdownTimeout = 3 * time.Second
	host            = "0.0.0.0"
	port            = "80"
)

type config struct {
	ShutdownTimeout time.Duration
	Host            string
	Port            string
}

func runServer(ctx context.Context, config *config) error {
	if config == nil {
		return errors.New("invalid config=nil")
	}

	srv := server.NewServer()
	api := rest.Init(srv)
	if api == nil {
		return errors.New("API wasn't created")
	}

	api.InitMiddleware()
	api.InitMetric()

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.Host, config.Port),
		Handler: api.Routes.Root,
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				ctx, _ := context.WithTimeout(context.Background(), config.ShutdownTimeout)
				if err := httpServer.Shutdown(ctx); err != nil {
					log.Fatalf("Server failed: %v", err)
				}
				return
			}
		}
	}()

	err := httpServer.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	config := &config{
		ShutdownTimeout: shutdownTimeout,
		Host:            host,
		Port:            port,
	}

	if err := runServer(ctx, config); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
