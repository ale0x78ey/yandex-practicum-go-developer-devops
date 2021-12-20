package agent

import (
	"context"
	"log"
	"net/url"
	"sync"
	"time"
)

type ClientConfig struct {
	ServerHost     string
	ServerPort     string
	ReportInterval time.Duration
}

type Client interface {
	Run(context.Context) error
	UpdateMetrics(*Metrics)
}

type client struct {
	m         sync.Mutex
	config    ClientConfig
	serverURL url.URL
	metrics   Metrics
}

func (c *client) Run(ctx context.Context) error {
	// TODO: init client?
	ticker := time.NewTicker(c.config.ReportInterval)
	for {
		select {
		case <-ticker.C:
			c.sendMetrics()
		case <-ctx.Done():
			log.Printf("Client: %v", ctx.Err())
			return nil
		}
	}
}

func (c *client) UpdateMetrics(metrics *Metrics) {
	if metrics == nil {
		log.Print("Client: UpdateMetrics metrics==nil")
		return
	}

	log.Print("Client: UpdateMetrics")
	c.m.Lock()
	c.metrics = *metrics
	c.m.Unlock()
}

func (c *client) sendMetrics() {
	log.Print("Client: sendMetrics")
	c.m.Lock()
	defer c.m.Unlock()

	// TODO: ...
}

func NewClient(config ClientConfig) (Client, error) {
	newClient := &client{
		config: config,
	}
	return newClient, nil
}
