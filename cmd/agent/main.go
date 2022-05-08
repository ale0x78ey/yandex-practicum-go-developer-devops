package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os/signal"
	"syscall"
	"time"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/config"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/service/agent"
)

const (
	updateURLFormat = "http://%s/update/"
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

	cfg := config.LoadAgentConfig()
	if cfg.Agent == nil {
		log.Fatalf("Missing config for agent")
	}

	agent, err := agent.NewAgent(
		*cfg.Agent,
		fmt.Sprintf(updateURLFormat, cfg.Http.ServerAddress),
	)
	if err != nil {
		log.Fatalf("Failed to create an agent: %v", err)
	}

	agent.Run(ctx)
}
