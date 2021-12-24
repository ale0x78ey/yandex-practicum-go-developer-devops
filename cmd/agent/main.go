package main

import (
	"context"
	"log"
	"math/rand"
	"os/signal"
	"syscall"
	"time"

	mag "github.com/ale0x78ey/yandex-practicum-go-developer-devops/service/agent"
)

const (
	pollInterval        = 2 * time.Second
	reportInterval      = 10 * time.Second
	serverHost          = "127.0.0.1"
	serverPort          = "80"
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
		PollInterval:        pollInterval,
		ReportInterval:      reportInterval,
		ServerHost:          serverHost,
		ServerPort:          serverPort,
		MaxIdleConns:        maxIdleConns,
		MaxIdleConnsPerHost: maxIdleConnsPerHost,
	})
	if err != nil {
		log.Fatalf("Failed to create an agent: %v", err)
	}

	if err := agent.Run(ctx); err != nil {
		log.Fatalf("Failed to run an agent: %v", err)
	}
}
