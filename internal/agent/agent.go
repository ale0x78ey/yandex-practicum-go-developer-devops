package agent

import (
	"context"
	"fmt"
	"log"
	"time"
)

type AgentConfig struct {
	PollInterval time.Duration
}

type Agent interface {
	Run(ctx context.Context, consumer MetricsConsumer) error
}

type agent struct {
	config    AgentConfig
	pollCount Counter
}

func (a *agent) Run(ctx context.Context, consumer MetricsConsumer) error {
	log.Print("Agent: Run")
	if a.config.PollInterval <= 0 {
		msg := "Invalid non-positive PollInterval=%v"
		return fmt.Errorf(msg, a.config.PollInterval)
	}

	ticker := time.NewTicker(a.config.PollInterval)
	for {
		select {
		case <-ticker.C:
			metrics := makeMetrics(a.pollCount)
			a.pollCount++
			consumer.Consume(&metrics)
		case <-ctx.Done():
			log.Printf("Agent: %v", ctx.Err())
			return nil
		}
	}
}

func NewAgent(config AgentConfig) (Agent, error) {
	newAgent := &agent{
		config: config,
	}
	return newAgent, nil
}
