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
)

const (
	shutdownTimeout = 3 * time.Second
	host            = "0.0.0.0"
	port            = "80"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	restServer := rest.NewServer()
	if restServer == nil {
		log.Fatal("Server failed: REST Server wasn't created")
	}

	router := restServer.NewRouter()
	if router == nil {
		log.Fatal("Server failed: Router wasn't created")
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: router,
	}

	go func() {
		err := httpServer.ListenAndServe()
		if err == http.ErrServerClosed {
			log.Print(err)
			return
		}
		if err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

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
}
