package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/api/rest"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/service/server"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/storage/psql"
)

const (
	// TODO: https://github.com/spf13/viper
	shutdownTimeout = 5 * time.Second
	host            = "0.0.0.0"
	port            = "8080"
)

func runServer(ctx context.Context) error {
	metricStorer := psql.NewMetricStorer()
	srv, err := server.NewServer(metricStorer)
	if err != nil {
		return err
	}

	h, err := rest.NewHandler(srv)
	if err != nil {
		return err
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: h.Router,
	}

	go func() {
		<-ctx.Done()
		ctx2, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := httpServer.Shutdown(ctx2); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
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
