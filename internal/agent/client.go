package agent

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type ClientConfig struct {
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	ServerHost          string
	ServerPort          string
	ReportInterval      time.Duration
}

type Client interface {
	Run(ctx context.Context) error
	UpdateMetrics(metrics *Metrics)
}

type client struct {
	m       sync.Mutex
	config  ClientConfig
	client  *http.Client
	metrics Metrics
}

func (c *client) Run(ctx context.Context) error {
	ticker := time.NewTicker(c.config.ReportInterval)
	for {
		select {
		case <-ticker.C:
			ctx2, _ := context.WithTimeout(ctx, c.config.ReportInterval)
			c.sendMetrics(ctx2)
		case <-ctx.Done():
			log.Printf("Client.Run %v", ctx.Err())
			return nil
		}
	}
}

func (c *client) UpdateMetrics(metrics *Metrics) {
	if metrics == nil {
		log.Print("Client: UpdateMetrics metrics==nil")
		return
	}

	log.Print("Client.UpdateMetrics")
	c.m.Lock()
	defer c.m.Unlock()

	c.metrics = *metrics
}

func (c *client) post(ctx context.Context, url string) error {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		nil,
	)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "text/html; charset=utf-8")

	_, err = c.client.Do(request)
	return err
}

func (c *client) sendGauge(ctx context.Context, name string, value Gauge) {
	go func() {
		url := fmt.Sprintf(
			"https://%s:%s/update/gauge/%s/%v",
			c.config.ServerHost,
			c.config.ServerPort,
			name, value,
		)
		if err := c.post(ctx, url); err != nil {
			log.Printf("Client.sendGauge error: %v", err)
		}
	}()
}

func (c *client) sendCounter(ctx context.Context, name string, value Counter) {
	go func() {
		url := fmt.Sprintf(
			"https://%s:%s/update/counter/%s/%v",
			c.config.ServerHost,
			c.config.ServerPort,
			name, value,
		)
		if err := c.post(ctx, url); err != nil {
			log.Printf("Client.sendCounter error: %v", err)
		}
	}()
}

func (c *client) sendMetrics(ctx context.Context) {
	log.Print("Client.sendMetrics")
	c.m.Lock()
	defer c.m.Unlock()

	c.sendGauge(ctx, "Alloc", Gauge(c.metrics.MemStats.Alloc))
	c.sendGauge(ctx, "BuckHashSys", Gauge(c.metrics.MemStats.BuckHashSys))
	c.sendGauge(ctx, "Frees", Gauge(c.metrics.MemStats.Frees))
	c.sendGauge(ctx, "GCCPUFraction", Gauge(c.metrics.MemStats.GCCPUFraction))
	c.sendGauge(ctx, "GCSys", Gauge(c.metrics.MemStats.GCSys))
	c.sendGauge(ctx, "HeapAlloc", Gauge(c.metrics.MemStats.HeapAlloc))
	c.sendGauge(ctx, "HeapIdle", Gauge(c.metrics.MemStats.HeapIdle))
	c.sendGauge(ctx, "HeapInuse", Gauge(c.metrics.MemStats.HeapInuse))
	c.sendGauge(ctx, "HeapObjects", Gauge(c.metrics.MemStats.HeapObjects))
	c.sendGauge(ctx, "HeapReleased", Gauge(c.metrics.MemStats.HeapReleased))
	c.sendGauge(ctx, "HeapSys", Gauge(c.metrics.MemStats.HeapSys))
	c.sendGauge(ctx, "LastGC", Gauge(c.metrics.MemStats.LastGC))
	c.sendGauge(ctx, "Lookups", Gauge(c.metrics.MemStats.Lookups))
	c.sendGauge(ctx, "MCacheInuse", Gauge(c.metrics.MemStats.MCacheInuse))
	c.sendGauge(ctx, "MCacheSys", Gauge(c.metrics.MemStats.MCacheSys))
	c.sendGauge(ctx, "MSpanInuse", Gauge(c.metrics.MemStats.MSpanInuse))
	c.sendGauge(ctx, "MSpanSys", Gauge(c.metrics.MemStats.MSpanSys))
	c.sendGauge(ctx, "Mallocs", Gauge(c.metrics.MemStats.Mallocs))
	c.sendGauge(ctx, "NextGC", Gauge(c.metrics.MemStats.NextGC))
	c.sendGauge(ctx, "NumForcedGC", Gauge(c.metrics.MemStats.NumForcedGC))
	c.sendGauge(ctx, "NumGC", Gauge(c.metrics.MemStats.NumGC))
	c.sendGauge(ctx, "OtherSys", Gauge(c.metrics.MemStats.OtherSys))
	c.sendGauge(ctx, "PauseTotalNs", Gauge(c.metrics.MemStats.PauseTotalNs))
	c.sendGauge(ctx, "StackInuse", Gauge(c.metrics.MemStats.StackInuse))
	c.sendGauge(ctx, "StackSys", Gauge(c.metrics.MemStats.StackSys))
	c.sendGauge(ctx, "Sys", Gauge(c.metrics.MemStats.Sys))
	c.sendGauge(ctx, "RandomValue", c.metrics.RandomValue)

	c.sendCounter(ctx, "PollCount", c.metrics.PollCount)
}

func NewClient(config ClientConfig) (Client, error) {
	t := &http.Transport{}
	t.MaxIdleConns = config.MaxIdleConns
	t.MaxIdleConnsPerHost = config.MaxIdleConnsPerHost

	newClient := &client{
		config: config,
		client: &http.Client{
			Transport: t,
		},
	}

	return newClient, nil
}
