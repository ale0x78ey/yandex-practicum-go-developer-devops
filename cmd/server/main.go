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
	// TODO: https://github.com/spf13/viper
	shutdownTimeout = 3 * time.Second
	host            = "0.0.0.0"
	port            = "8080"
)

func runServer(ctx context.Context) error {
	srv := server.NewServer()
	api := rest.Init(srv)
	if api == nil {
		return errors.New("API wasn't created")
	}

	api.InitMiddleware()
	api.InitMetric()

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: api.Routes.Root,
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				ctx, _ := context.WithTimeout(context.Background(), shutdownTimeout)
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

	if err := runServer(ctx); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
