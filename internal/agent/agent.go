package agent

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Consumer interface {
	Consume(*Metrics)
}

type ConsumerFunc func(*Metrics)

func (f ConsumerFunc) Consume(metrics *Metrics) {
	log.Print("Consume")
	f(metrics)
}

type Config struct {
	PollInterval time.Duration
}

type Agent interface {
	Run(ctx context.Context, consumer Consumer) error
}

type agent struct {
	config    Config
	pollCount Counter
}

func (a *agent) Run(ctx context.Context, consumer Consumer) error {
	log.Print("Run agent")
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

	panic("The unreachable point was executed.")
}

func NewAgent(config Config) (Agent, error) {
	newAgent := &agent{
		config: config,
	}
	return newAgent, nil
}
