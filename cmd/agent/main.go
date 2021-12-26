package main

import (
	"context"
	"log"
	"math/rand"
	"os/signal"
	"syscall"
	"time"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/service/agent"
)

const (
	pollInterval        = 2 * time.Second
	reportInterval      = 10 * time.Second
	serverHost          = "127.0.0.1"
	serverPort          = "80"
	maxIdleConns        = 15
	maxIdleConnsPerHost = 15
	retryCount          = 5
	retryWaitTime       = 10 * time.Second
	retryMaxWaitTime    = 60 * time.Second
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

	agent := agent.NewAgent(agent.Config{
		PollInterval:        pollInterval,
		ReportInterval:      reportInterval,
		ServerHost:          serverHost,
		ServerPort:          serverPort,
		MaxIdleConns:        maxIdleConns,
		MaxIdleConnsPerHost: maxIdleConnsPerHost,
		RetryCount:          retryCount,
		RetryWaitTime:       retryWaitTime,
		RetryMaxWaitTime:    retryMaxWaitTime,
	})
	if agent == nil {
		log.Fatalf("Failed to create an agent")
	}

	if err := agent.Run(ctx); err != nil {
		log.Fatalf("Failed to run an agent: %v", err)
	}
}
