package agent

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
)

const (
	metricPostURL = "http://{host}:{port}/update/{metricType}/{metricName}/{metricValue}"
)

type Config struct {
	PollInterval        time.Duration
	ReportInterval      time.Duration
	ServerHost          string
	ServerPort          string
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	RetryCount          int
	RetryWaitTime       time.Duration
	RetryMaxWaitTime    time.Duration
}

type metrics struct {
	memStats    runtime.MemStats
	randomValue float64
	pollCount   int64
}

type Agent struct {
	config *Config
	client *resty.Client
	data   metrics
	wg     sync.WaitGroup
}

func minDuration(d1, d2 time.Duration) time.Duration {
	if d1 <= d2 {
		return d1
	}
	return d2
}

func (a *Agent) Run(ctx context.Context) error {
	if a.config.PollInterval <= 0 {
		msg := "invalid non-positive PollInterval=%v"
		return fmt.Errorf(msg, a.config.PollInterval)
	}
	if a.config.ReportInterval <= 0 {
		msg := "invalid non-positive ReportInterval=%v"
		return fmt.Errorf(msg, a.config.ReportInterval)
	}

	postTimeout := minDuration(a.config.ReportInterval, a.config.PollInterval)

	pollTicker := time.NewTicker(a.config.PollInterval)
	sendTicker := time.NewTicker(a.config.ReportInterval)
	for {
		select {
		case <-pollTicker.C:
			a.pollMetrics()
		case <-sendTicker.C:
			ctx2, cancel := context.WithTimeout(ctx, postTimeout)
			defer cancel()
			a.postMetrics(ctx2)
		case <-ctx.Done():
			return nil
		}
	}
}

func (a *Agent) pollMetrics() {
	runtime.ReadMemStats(&a.data.memStats)
	a.data.randomValue = rand.Float64()
	a.data.pollCount++
}

func (a *Agent) post(
	ctx context.Context,
	metricType model.MetricType,
	metricName model.MetricName,
	value fmt.Stringer,
) {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		request := a.client.R().
			SetContext(ctx).
			SetHeader("content-type", "text/plain").
			SetPathParams(map[string]string{
				"host":        a.config.ServerHost,
				"port":        a.config.ServerPort,
				"metricType":  metricType.String(),
				"metricName":  metricName.String(),
				"metricValue": value.String(),
			})

		request.Post(metricPostURL)
	}()
}

func (a *Agent) postMetrics(ctx context.Context) {
	gauge := model.MetricTypeGauge
	counter := model.MetricTypeCounter
	m := &a.data.memStats

	a.post(ctx, gauge, model.MetricNameAlloc, model.Gauge(m.Alloc))
	a.post(ctx, gauge, model.MetricNameBuckHashSys, model.Gauge(m.BuckHashSys))
	a.post(ctx, gauge, model.MetricNameFrees, model.Gauge(m.Frees))
	a.post(ctx, gauge, model.MetricNameGCCPUFraction, model.Gauge(m.GCCPUFraction))
	a.post(ctx, gauge, model.MetricNameGCSys, model.Gauge(m.GCSys))
	a.post(ctx, gauge, model.MetricNameHeapAlloc, model.Gauge(m.HeapAlloc))
	a.post(ctx, gauge, model.MetricNameHeapIdle, model.Gauge(m.HeapIdle))
	a.post(ctx, gauge, model.MetricNameHeapInuse, model.Gauge(m.HeapInuse))
	a.post(ctx, gauge, model.MetricNameHeapObjects, model.Gauge(m.HeapObjects))
	a.post(ctx, gauge, model.MetricNameHeapReleased, model.Gauge(m.HeapReleased))
	a.post(ctx, gauge, model.MetricNameHeapSys, model.Gauge(m.HeapSys))
	a.post(ctx, gauge, model.MetricNameLastGC, model.Gauge(m.LastGC))
	a.post(ctx, gauge, model.MetricNameLookups, model.Gauge(m.Lookups))
	a.post(ctx, gauge, model.MetricNameMCacheInuse, model.Gauge(m.MCacheInuse))
	a.post(ctx, gauge, model.MetricNameMCacheSys, model.Gauge(m.MCacheSys))
	a.post(ctx, gauge, model.MetricNameMSpanInuse, model.Gauge(m.MSpanInuse))
	a.post(ctx, gauge, model.MetricNameMSpanSys, model.Gauge(m.MSpanSys))
	a.post(ctx, gauge, model.MetricNameMallocs, model.Gauge(m.Mallocs))
	a.post(ctx, gauge, model.MetricNameNextGC, model.Gauge(m.NextGC))
	a.post(ctx, gauge, model.MetricNameNumForcedGC, model.Gauge(m.NumForcedGC))
	a.post(ctx, gauge, model.MetricNameNumGC, model.Gauge(m.NumGC))
	a.post(ctx, gauge, model.MetricNameOtherSys, model.Gauge(m.OtherSys))
	a.post(ctx, gauge, model.MetricNamePauseTotalNs, model.Gauge(m.PauseTotalNs))
	a.post(ctx, gauge, model.MetricNameStackInuse, model.Gauge(m.StackInuse))
	a.post(ctx, gauge, model.MetricNameStackSys, model.Gauge(m.StackSys))
	a.post(ctx, gauge, model.MetricNameSys, model.Gauge(m.Sys))
	a.post(ctx, gauge, model.MetricNameRandomValue, model.Gauge(a.data.randomValue))
	a.post(ctx, counter, model.MetricNamePollCount, model.Counter(a.data.pollCount))

	a.wg.Wait()
}

func NewAgent(config *Config) *Agent {
	if config == nil {
		return nil
	}

	t := &http.Transport{}
	t.MaxIdleConns = config.MaxIdleConns
	t.MaxIdleConnsPerHost = config.MaxIdleConnsPerHost

	httpClient := &http.Client{
		Transport: t,
	}

	client := resty.NewWithClient(httpClient)
	if client == nil {
		return nil
	}

	client.
		SetRetryCount(config.RetryCount).
		SetRetryWaitTime(config.RetryWaitTime).
		SetRetryMaxWaitTime(config.RetryMaxWaitTime)

	return &Agent{
		config: config,
		client: client,
	}
}
