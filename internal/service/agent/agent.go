package agent

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

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
	DefaultPollMetricsBuffSize = 100
	DefaultPostWorkersPoolSize = 15
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
	PollMetricsBuffSize int
	PostWorkersPoolSize int
}

func (c Config) Validate() error {
	if c.PollInterval <= 0 {
		return fmt.Errorf("invalid non-positive PollInterval=%v", c.PollInterval)
	}
	if c.ReportInterval <= 0 {
		return fmt.Errorf("invalid non-positive ReportInterval=%v", c.ReportInterval)
	}
	if c.PollMetricsBuffSize < 0 {
		return fmt.Errorf("invalid negative PollMetricsBuffSize=%v", c.PollMetricsBuffSize)
	}
	if c.PostWorkersPoolSize <= 0 {
		return fmt.Errorf("invalid non-positive PostWorkersPoolSize=%v", c.PostWorkersPoolSize)
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
		metrics:   make(chan model.Metric, config.PollMetricsBuffSize),
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

func (a *Agent) pollMetrics(ctx context.Context, wg *sync.WaitGroup) {
	ticker := time.NewTicker(a.config.PollInterval)
	defer ticker.Stop()
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.pollMemStatsMetrics(ctx)
			a.pollRandomValueMetric(ctx)
			a.pollCountMetric(ctx)
			a.pollGopsutilMetrics(ctx)
		}
	}
}

func (a *Agent) queueMetric(ctx context.Context, metric model.Metric) {
	select {
	case <-ctx.Done():
		return
	case a.metrics <- metric:
		return
	}
}

func (a *Agent) pollMemStatsMetrics(ctx context.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	a.queueMetric(ctx, model.MetricFromGauge("Alloc", model.Gauge(m.Alloc)))
	a.queueMetric(ctx, model.MetricFromGauge("TotalAlloc", model.Gauge(m.TotalAlloc)))
	a.queueMetric(ctx, model.MetricFromGauge("BuckHashSys", model.Gauge(m.BuckHashSys)))
	a.queueMetric(ctx, model.MetricFromGauge("Frees", model.Gauge(m.Frees)))
	a.queueMetric(ctx, model.MetricFromGauge("GCCPUFraction", model.Gauge(m.GCCPUFraction)))
	a.queueMetric(ctx, model.MetricFromGauge("GCSys", model.Gauge(m.GCSys)))
	a.queueMetric(ctx, model.MetricFromGauge("HeapAlloc", model.Gauge(m.HeapAlloc)))
	a.queueMetric(ctx, model.MetricFromGauge("HeapIdle", model.Gauge(m.HeapIdle)))
	a.queueMetric(ctx, model.MetricFromGauge("HeapInuse", model.Gauge(m.HeapInuse)))
	a.queueMetric(ctx, model.MetricFromGauge("HeapObjects", model.Gauge(m.HeapObjects)))
	a.queueMetric(ctx, model.MetricFromGauge("HeapReleased", model.Gauge(m.HeapReleased)))
	a.queueMetric(ctx, model.MetricFromGauge("HeapSys", model.Gauge(m.HeapSys)))
	a.queueMetric(ctx, model.MetricFromGauge("LastGC", model.Gauge(m.LastGC)))
	a.queueMetric(ctx, model.MetricFromGauge("Lookups", model.Gauge(m.Lookups)))
	a.queueMetric(ctx, model.MetricFromGauge("MCacheInuse", model.Gauge(m.MCacheInuse)))
	a.queueMetric(ctx, model.MetricFromGauge("MCacheSys", model.Gauge(m.MCacheSys)))
	a.queueMetric(ctx, model.MetricFromGauge("MSpanInuse", model.Gauge(m.MSpanInuse)))
	a.queueMetric(ctx, model.MetricFromGauge("MSpanSys", model.Gauge(m.MSpanSys)))
	a.queueMetric(ctx, model.MetricFromGauge("Mallocs", model.Gauge(m.Mallocs)))
	a.queueMetric(ctx, model.MetricFromGauge("NextGC", model.Gauge(m.NextGC)))
	a.queueMetric(ctx, model.MetricFromGauge("NumForcedGC", model.Gauge(m.NumForcedGC)))
	a.queueMetric(ctx, model.MetricFromGauge("NumGC", model.Gauge(m.NumGC)))
	a.queueMetric(ctx, model.MetricFromGauge("OtherSys", model.Gauge(m.OtherSys)))
	a.queueMetric(ctx, model.MetricFromGauge("PauseTotalNs", model.Gauge(m.PauseTotalNs)))
	a.queueMetric(ctx, model.MetricFromGauge("StackInuse", model.Gauge(m.StackInuse)))
	a.queueMetric(ctx, model.MetricFromGauge("StackSys", model.Gauge(m.StackSys)))
	a.queueMetric(ctx, model.MetricFromGauge("Sys", model.Gauge(m.Sys)))
}

func (a *Agent) pollRandomValueMetric(ctx context.Context) {
	randomValue := rand.Float64()
	a.queueMetric(ctx, model.MetricFromGauge("RandomValue", model.Gauge(randomValue)))
}

func (a *Agent) pollCountMetric(ctx context.Context) {
	pollCount := atomic.AddInt64(&a.pollCount, 1)
	a.queueMetric(ctx, model.MetricFromCounter("PollCount", model.Counter(pollCount)))
}

func (a *Agent) pollGopsutilMetrics(ctx context.Context) {
	if err := a.pollGopsutilMemoryMetrics(ctx); err != nil {
		log.Printf("failed to poll gopsutil memory metrics: %v", err)
	}

	if err := a.pollGopsutilCPUMetrics(ctx); err != nil {
		log.Printf("failed to poll gopsutil cpu metrics: %v", err)
	}
}

func (a *Agent) pollGopsutilMemoryMetrics(ctx context.Context) error {
	v, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return err
	}

	a.queueMetric(ctx, model.MetricFromGauge("TotalMemory", model.Gauge(v.Total)))
	a.queueMetric(ctx, model.MetricFromGauge("FreeMemory", model.Gauge(v.Free)))

	return nil
}

func (a *Agent) pollGopsutilCPUMetrics(ctx context.Context) error {
	times, err := cpu.TimesWithContext(ctx, true)
	if err != nil {
		return err
	}

	for i, timesStat := range times {
		a.queueMetric(
			ctx,
			model.MetricFromGauge(
				fmt.Sprintf("CPUutilization%d", i),
				model.Gauge(timesStat.User+timesStat.System),
			))
	}

	return nil
}

func (a *Agent) postMetrics(ctx context.Context, wg *sync.WaitGroup) {
	ticker := time.NewTicker(a.config.ReportInterval)
	defer ticker.Stop()
	defer wg.Done()

	done := false
	signal := sync.NewCond(&sync.Mutex{})

	for i := 0; i < a.config.PostWorkersPoolSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				signal.L.Lock()
				if !done {
					signal.Wait()
					signal.L.Unlock()

					for {
						select {
						case <-ctx.Done():
							return
						case metric := <-a.metrics:
							if err := a.postOneMetric(ctx, metric); err != nil {
								log.Printf("failed to post %v: %v", metric, err)
							}
						default:
							break
						}
					}

					continue
				}

				signal.L.Unlock()
				return
			}
		}()
	}

	for {
		select {
		case <-ctx.Done():
			signal.L.Lock()
			done = true
			signal.Broadcast()
			signal.L.Unlock()
			return
		case <-ticker.C:
			signal.L.Lock()
			signal.Broadcast()
			signal.L.Unlock()
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
