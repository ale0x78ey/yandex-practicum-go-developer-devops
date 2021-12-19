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

func processMetrics(metrics *mag.Metrics) {
	if metrics == nil {
		panic("The agent passed a nil object.")
	}

	log.Printf("Process metrics with PollCount=%v rand=%v", metrics.PollCount, metrics.RandomValue)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	agent, err := mag.NewAgent(mag.Config{
		PollInterval: 2 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to create an agent: %v", err)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	if err := agent.Run(ctx, mag.ConsumerFunc(processMetrics)); err != nil {
		log.Fatalf("Agent failed: %v", err)
	}
}
