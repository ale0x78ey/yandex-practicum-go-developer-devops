package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/api/rest"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/service/server"
	storagefile "github.com/ale0x78ey/yandex-practicum-go-developer-devops/storage/file"
)

const (
	shutdownTimeout = 5 * time.Second
)

type config struct {
	ShutdownTimeout time.Duration
	ServerAddress   string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	StoreInterval   time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile       string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore         bool          `env:"RESTORE" envDefault:"true"`
}

type restServer struct {
	config     config
	server     *server.Server
	httpServer *http.Server
}

func newRestServer(config config) (*restServer, error) {
	s, err := server.NewServer(storagefile.NewMetricStorer())
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
	}

	return server, nil
}

func (s restServer) Run(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				ctx, cancel := context.WithTimeout(
					context.Background(), s.config.ShutdownTimeout)
				defer cancel()
				if err := s.httpServer.Shutdown(ctx); err != nil {
					log.Fatalf("HTTP Server failed: %v", err)
				}
				if err := s.server.Flush(ctx); err != nil {
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

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	config := config{
		ShutdownTimeout: shutdownTimeout,
	}
	if err := env.Parse(&config); err != nil {
		log.Fatalf("Failed to parse config options: %v", err)
	}

	server, err := newRestServer(config)
	if err != nil {
		log.Fatalf("Failed to create a server: %v", err)
	}

	if err := server.Run(ctx); err != nil {
		log.Fatalf("Failed to run a server: %v", err)
	}
}
