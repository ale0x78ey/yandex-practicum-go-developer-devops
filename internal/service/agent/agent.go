package agent

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/internal/model"
)

const (
	DefaultMaxIdleConns        = 15
	DefaultMaxIdleConnsPerHost = 15
	DefaultRetryCount          = 1
	DefaultRetryWaitTime       = 100 * time.Millisecond
	DefaultRetryMaxWaitTime    = 900 * time.Millisecond
	DefaultPollInterval        = 02 * time.Second
	DefaultReportInterval      = 10 * time.Second
)

type Config struct {
	PollInterval        time.Duration `env:"POLL_INTERVAL"`
	ReportInterval      time.Duration `env:"REPORT_INTERVAL"`
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	RetryCount          int
	RetryWaitTime       time.Duration
	RetryMaxWaitTime    time.Duration
	Key                 string `env:"KEY"`
}

func (c Config) Validate() error {
	if c.PollInterval <= 0 {
		return fmt.Errorf("invalid non-positive PollInterval=%v", c.PollInterval)
	}
	if c.ReportInterval <= 0 {
		return fmt.Errorf("invalid non-positive ReportInterval=%v", c.ReportInterval)
	}

	return nil
}

type Agent struct {
	config    Config
	client    *resty.Client
	updateURL string
	metrics   chan model.Metric

	pollCount int64
}

func NewAgent(config Config, updateURL string) (*Agent, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	t := &http.Transport{}
	t.MaxIdleConns = config.MaxIdleConns
	t.MaxIdleConnsPerHost = config.MaxIdleConnsPerHost

	httpClient := &http.Client{
		Transport: t,
	}

	client := resty.NewWithClient(httpClient)
	client.
		SetRetryCount(config.RetryCount).
		SetRetryWaitTime(config.RetryWaitTime).
		SetRetryMaxWaitTime(config.RetryMaxWaitTime)

	a := &Agent{
		config:    config,
		client:    client,
		updateURL: updateURL,
		metrics:   make(chan model.Metric),
	}

	return a, nil
}

func (a *Agent) Run(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	defer wg.Wait()

	go a.pollMetrics(ctx, wg)
	go a.postMetrics(ctx, wg)

	return nil
}

func (a *Agent) pollMemStatsMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	a.metrics <- model.MetricFromGauge("Alloc", model.Gauge(m.Alloc))
	a.metrics <- model.MetricFromGauge("TotalAlloc", model.Gauge(m.TotalAlloc))
	a.metrics <- model.MetricFromGauge("BuckHashSys", model.Gauge(m.BuckHashSys))
	a.metrics <- model.MetricFromGauge("Frees", model.Gauge(m.Frees))
	a.metrics <- model.MetricFromGauge("GCCPUFraction", model.Gauge(m.GCCPUFraction))
	a.metrics <- model.MetricFromGauge("GCSys", model.Gauge(m.GCSys))
	a.metrics <- model.MetricFromGauge("HeapAlloc", model.Gauge(m.HeapAlloc))
	a.metrics <- model.MetricFromGauge("HeapIdle", model.Gauge(m.HeapIdle))
	a.metrics <- model.MetricFromGauge("HeapInuse", model.Gauge(m.HeapInuse))
	a.metrics <- model.MetricFromGauge("HeapObjects", model.Gauge(m.HeapObjects))
	a.metrics <- model.MetricFromGauge("HeapReleased", model.Gauge(m.HeapReleased))
	a.metrics <- model.MetricFromGauge("HeapSys", model.Gauge(m.HeapSys))
	a.metrics <- model.MetricFromGauge("LastGC", model.Gauge(m.LastGC))
	a.metrics <- model.MetricFromGauge("Lookups", model.Gauge(m.Lookups))
	a.metrics <- model.MetricFromGauge("MCacheInuse", model.Gauge(m.MCacheInuse))
	a.metrics <- model.MetricFromGauge("MCacheSys", model.Gauge(m.MCacheSys))
	a.metrics <- model.MetricFromGauge("MSpanInuse", model.Gauge(m.MSpanInuse))
	a.metrics <- model.MetricFromGauge("MSpanSys", model.Gauge(m.MSpanSys))
	a.metrics <- model.MetricFromGauge("Mallocs", model.Gauge(m.Mallocs))
	a.metrics <- model.MetricFromGauge("NextGC", model.Gauge(m.NextGC))
	a.metrics <- model.MetricFromGauge("NumForcedGC", model.Gauge(m.NumForcedGC))
	a.metrics <- model.MetricFromGauge("NumGC", model.Gauge(m.NumGC))
	a.metrics <- model.MetricFromGauge("OtherSys", model.Gauge(m.OtherSys))
	a.metrics <- model.MetricFromGauge("PauseTotalNs", model.Gauge(m.PauseTotalNs))
	a.metrics <- model.MetricFromGauge("StackInuse", model.Gauge(m.StackInuse))
	a.metrics <- model.MetricFromGauge("StackSys", model.Gauge(m.StackSys))
	a.metrics <- model.MetricFromGauge("Sys", model.Gauge(m.Sys))
}

func (a *Agent) pollRandomValueMetric() {
	randomValue := rand.Float64()
	a.metrics <- model.MetricFromGauge("RandomValue", model.Gauge(randomValue))
}

func (a *Agent) pollCountMetric() {
	a.pollCount++
	a.metrics <- model.MetricFromCounter("PollCount", model.Counter(a.pollCount))
}

func (a *Agent) pollMetrics(ctx context.Context, wg *sync.WaitGroup) {
	ticker := time.NewTicker(a.config.PollInterval)
	defer ticker.Stop()
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.pollMemStatsMetrics()
			a.pollRandomValueMetric()
			a.pollCountMetric()
		}
	}
}

func (a *Agent) postOneMetric(ctx context.Context, metric model.Metric) error {
	if err := metric.UpdateHash(a.config.Key); err != nil {
		return err
	}

	_, err := a.client.R().
		SetContext(ctx).
		SetBody(metric).
		Post(a.updateURL)

	return err
}

func (a *Agent) postMetrics(ctx context.Context, wg *sync.WaitGroup) {
	ticker := time.NewTicker(a.config.ReportInterval)
	defer ticker.Stop()
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for metric := range a.metrics {
				if err := a.postOneMetric(ctx, metric); err != nil {
					log.Printf("failed to post %v: %v", metric, err)
				}
			}
		}
	}
}
