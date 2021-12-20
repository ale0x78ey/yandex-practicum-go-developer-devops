package main

import (
	"context"
	"log"
	"math/rand"
	"os/signal"
	"syscall"
	"time"

	mag "github.com/ale0x78ey/yandex-practicum-go-developer-devops/internal/agent"
)

const (
	pollInterval        = 2 * time.Second
	reportInterval      = 10 * time.Second
	serverHost          = "127.0.0.1"
	serverPort          = "8080"
	maxIdleConns        = 100
	maxIdleConnsPerHost = 100
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	agent, err := mag.NewAgent(mag.AgentConfig{
		PollInterval: pollInterval,
	})
	if err != nil {
		log.Fatalf("Failed to create an agent: %v", err)
	}

	client, err := mag.NewClient(mag.ClientConfig{
		ServerHost:          serverHost,
		ServerPort:          serverPort,
		ReportInterval:      reportInterval,
		MaxIdleConns:        maxIdleConns,
		MaxIdleConnsPerHost: maxIdleConnsPerHost,
	})
	if err != nil {
		log.Fatalf("Failed to create an agent client: %v", err)
	}

	go func() {
		if err := client.Run(ctx); err != nil {
			log.Fatalf("Agent failed: %v", err)
		}
	}()

	if err := agent.Run(ctx, mag.MetricsConsumerFunc(client.UpdateMetrics)); err != nil {
		log.Fatalf("Agent failed: %v", err)
	}
}
