package agent

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
)

type AgentConfig struct {
	PollInterval        time.Duration
	ReportInterval      time.Duration
	ServerHost          string
	ServerPort          string
	MaxIdleConns        int
	MaxIdleConnsPerHost int
}

type Agent interface {
	Run(ctx context.Context) error
}

type metrics struct {
	memStats    runtime.MemStats
	randomValue float64
	pollCount   int64
}

type agent struct {
	config AgentConfig
	client *http.Client
	data   metrics
}

func (a *agent) Run(ctx context.Context) error {
	if a.config.PollInterval <= 0 {
		msg := "Invalid non-positive PollInterval=%v"
		return fmt.Errorf(msg, a.config.PollInterval)
	}

	pollTicker := time.NewTicker(a.config.PollInterval)
	sendTicker := time.NewTicker(a.config.ReportInterval)
	for {
		select {
		case <-pollTicker.C:
			a.pollMetrics()
		case <-sendTicker.C:
			ctx2, _ := context.WithTimeout(ctx, a.config.ReportInterval)
			a.sendMetrics(ctx2)
		case <-ctx.Done():
			log.Printf("Agent.Run done: %v", ctx.Err())
			return nil
		}
	}
}

func (a *agent) pollMetrics() {
	runtime.ReadMemStats(&a.data.memStats)
	a.data.randomValue = rand.Float64()
	a.data.pollCount++
}

func (a *agent) post(ctx context.Context, url string) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		log.Printf("Agent.post %s error: %v", url, err)
		return
	}

	request.Header.Set("Content-Type", "text/html; charset=utf-8")
	if _, err = a.client.Do(request); err != nil {
		log.Printf("Agent.post %s error: %v", url, err)
	}

	log.Printf("Agent.post %s", url)
}

func (a *agent) send(ctx context.Context, value model.Metric) {
	go func() {
		url := fmt.Sprintf(
			"http://%s:%s/update/%s/%s/%s",
			a.config.ServerHost,
			a.config.ServerPort,
			value.Type(),
			value.Name(),
			value.StringValue(),
		)
		a.post(ctx, url)
	}()
}

func (a *agent) sendMetrics(ctx context.Context) {
	d := &a.data
	a.send(ctx, model.GaugeFromUInt64("Alloc", d.memStats.Alloc))
	a.send(ctx, model.GaugeFromUInt64("BuckHashSys", d.memStats.BuckHashSys))
	a.send(ctx, model.GaugeFromUInt64("Frees", d.memStats.Frees))
	a.send(ctx, model.GaugeFromFloat64("GCCPUFraction", d.memStats.GCCPUFraction))
	a.send(ctx, model.GaugeFromUInt64("GCSys", d.memStats.GCSys))
	a.send(ctx, model.GaugeFromUInt64("HeapAlloc", d.memStats.HeapAlloc))
	a.send(ctx, model.GaugeFromUInt64("HeapIdle", d.memStats.HeapIdle))
	a.send(ctx, model.GaugeFromUInt64("HeapInuse", d.memStats.HeapInuse))
	a.send(ctx, model.GaugeFromUInt64("HeapObjects", d.memStats.HeapObjects))
	a.send(ctx, model.GaugeFromUInt64("HeapReleased", d.memStats.HeapReleased))
	a.send(ctx, model.GaugeFromUInt64("HeapSys", d.memStats.HeapSys))

	a.send(ctx, model.GaugeFromUInt64("LastGC", d.memStats.LastGC))
	a.send(ctx, model.GaugeFromUInt64("Lookups", d.memStats.Lookups))
	a.send(ctx, model.GaugeFromUInt64("MCacheInuse", d.memStats.MCacheInuse))
	a.send(ctx, model.GaugeFromUInt64("MCacheSys", d.memStats.MCacheSys))
	a.send(ctx, model.GaugeFromUInt64("MSpanInuse", d.memStats.MSpanInuse))
	a.send(ctx, model.GaugeFromUInt64("MSpanSys", d.memStats.MSpanSys))
	a.send(ctx, model.GaugeFromUInt64("Mallocs", d.memStats.Mallocs))
	a.send(ctx, model.GaugeFromUInt64("NextGC", d.memStats.NextGC))
	a.send(ctx, model.GaugeFromUInt32("NumForcedGC", d.memStats.NumForcedGC))
	a.send(ctx, model.GaugeFromUInt32("NumGC", d.memStats.NumGC))
	a.send(ctx, model.GaugeFromUInt64("OtherSys", d.memStats.OtherSys))
	a.send(ctx, model.GaugeFromUInt64("PauseTotalNs", d.memStats.PauseTotalNs))
	a.send(ctx, model.GaugeFromUInt64("StackInuse", d.memStats.StackInuse))
	a.send(ctx, model.GaugeFromUInt64("StackSys", d.memStats.StackSys))
	a.send(ctx, model.GaugeFromUInt64("Sys", d.memStats.Sys))
	a.send(ctx, model.GaugeFromFloat64("RandomValue", d.randomValue))

	a.send(ctx, model.CounterFromInt64("PollCount", d.pollCount))
}

func NewAgent(config AgentConfig) (Agent, error) {
	t := &http.Transport{}
	t.MaxIdleConns = config.MaxIdleConns
	t.MaxIdleConnsPerHost = config.MaxIdleConnsPerHost

	client := &http.Client{
		Transport: t,
	}

	newAgent := &agent{
		config: config,
		client: client,
	}

	return newAgent, nil
}
