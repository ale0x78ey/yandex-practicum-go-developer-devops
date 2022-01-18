package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/service/agent"
)

const (
	maxIdleConns        = 15
	maxIdleConnsPerHost = 15
	retryCount          = 1
	retryWaitTime       = 100 * time.Millisecond
	retryMaxWaitTime    = 900 * time.Millisecond
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

	config := agent.Config{
		MaxIdleConns:        maxIdleConns,
		MaxIdleConnsPerHost: maxIdleConnsPerHost,
		RetryCount:          retryCount,
		RetryWaitTime:       retryWaitTime,
		RetryMaxWaitTime:    retryMaxWaitTime,
	}
	if err := env.Parse(&config); err != nil {
		log.Fatalf("Failed to parse config options: %v", err)
	}

	flag.StringVar(&config.ServerAddress, "a", config.ServerAddress, "ADDRESS")
	flag.DurationVar(&config.ReportInterval, "r", config.ReportInterval, "REPORT_INTERVAL")
	flag.DurationVar(&config.PollInterval, "p", config.PollInterval, "POLL_INTERVAL")
	flag.Parse()

	agent, err := agent.NewAgent(config)
	if err != nil {
		log.Fatalf("Failed to create an agent: %v", err)
	}

	if err := agent.Run(ctx); err != nil {
		log.Fatalf("Failed to run an agent: %v", err)
	}
}
