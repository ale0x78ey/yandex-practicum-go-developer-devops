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
	defaultMaxIdleConns        = 15
	defaultMaxIdleConnsPerHost = 15
	defaultRetryCount          = 1
	defaultRetryWaitTime       = 100 * time.Millisecond
	defaultRetryMaxWaitTime    = 900 * time.Millisecond
	defaultPollInterval        = 02 * time.Second
	defaultReportInterval      = 10 * time.Second
	defaultServerAddress       = "127.0.0.1:8080"
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
		MaxIdleConns:        defaultMaxIdleConns,
		MaxIdleConnsPerHost: defaultMaxIdleConnsPerHost,
		RetryCount:          defaultRetryCount,
		RetryWaitTime:       defaultRetryWaitTime,
		RetryMaxWaitTime:    defaultRetryMaxWaitTime,
	}

	flag.StringVar(&config.ServerAddress, "a", defaultServerAddress, "ADDRESS")
	flag.DurationVar(&config.ReportInterval, "r", defaultReportInterval, "REPORT_INTERVAL")
	flag.DurationVar(&config.PollInterval, "p", defaultPollInterval, "POLL_INTERVAL")
	flag.Parse()

	if err := env.Parse(&config); err != nil {
		log.Fatalf("Failed to parse config options: %v", err)
	}

	agent, err := agent.NewAgent(config)
	if err != nil {
		log.Fatalf("Failed to create an agent: %v", err)
	}

	if err := agent.Run(ctx); err != nil {
		log.Fatalf("Failed to run an agent: %v", err)
	}
}
