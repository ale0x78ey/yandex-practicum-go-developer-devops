package agent

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
)

const (
	metricPostUrl = "http://{host}:{port}/update/{metricType}/{metricName}/{metricValue}"
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
	config Config
	client *resty.Client
	data   metrics
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

	pollTicker := time.NewTicker(a.config.PollInterval)
	sendTicker := time.NewTicker(a.config.ReportInterval)
	for {
		select {
		case <-pollTicker.C:
			a.pollMetrics()
		case <-sendTicker.C:
			ctx2, cancel := context.WithTimeout(ctx, a.config.ReportInterval)
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
	value string,
) {
	// go func() {
	request := a.client.R().SetContext(ctx).SetPathParams(map[string]string{
		"host":        a.config.ServerHost,
		"port":        a.config.ServerPort,
		"metricType":  metricType.String(),
		"metricName":  metricName.String(),
		"metricValue": value,
	})

	request.Post(metricPostUrl)
	// }()
}

func (a *Agent) postMetrics(ctx context.Context) {
	gauge := model.MetricTypeGauge
	counter := model.MetricTypeCounter
	m := &a.data.memStats

	a.post(ctx, gauge, model.MetricNameAlloc, model.Gauge(m.Alloc).String())
	a.post(ctx, gauge, model.MetricNameBuckHashSys, model.Gauge(m.BuckHashSys).String())
	a.post(ctx, gauge, model.MetricNameFrees, model.Gauge(m.Frees).String())
	a.post(ctx, gauge, model.MetricNameGCCPUFraction, model.Gauge(m.GCCPUFraction).String())
	a.post(ctx, gauge, model.MetricNameGCSys, model.Gauge(m.GCSys).String())
	a.post(ctx, gauge, model.MetricNameHeapAlloc, model.Gauge(m.HeapAlloc).String())
	a.post(ctx, gauge, model.MetricNameHeapIdle, model.Gauge(m.HeapIdle).String())
	a.post(ctx, gauge, model.MetricNameHeapInuse, model.Gauge(m.HeapInuse).String())
	a.post(ctx, gauge, model.MetricNameHeapObjects, model.Gauge(m.HeapObjects).String())
	a.post(ctx, gauge, model.MetricNameHeapReleased, model.Gauge(m.HeapReleased).String())
	a.post(ctx, gauge, model.MetricNameHeapSys, model.Gauge(m.HeapSys).String())
	a.post(ctx, gauge, model.MetricNameLastGC, model.Gauge(m.LastGC).String())
	a.post(ctx, gauge, model.MetricNameLookups, model.Gauge(m.Lookups).String())
	a.post(ctx, gauge, model.MetricNameMCacheInuse, model.Gauge(m.MCacheInuse).String())
	a.post(ctx, gauge, model.MetricNameMCacheSys, model.Gauge(m.MCacheSys).String())
	a.post(ctx, gauge, model.MetricNameMSpanInuse, model.Gauge(m.MSpanInuse).String())
	a.post(ctx, gauge, model.MetricNameMSpanSys, model.Gauge(m.MSpanSys).String())
	a.post(ctx, gauge, model.MetricNameMallocs, model.Gauge(m.Mallocs).String())
	a.post(ctx, gauge, model.MetricNameNextGC, model.Gauge(m.NextGC).String())
	a.post(ctx, gauge, model.MetricNameNumForcedGC, model.Gauge(m.NumForcedGC).String())
	a.post(ctx, gauge, model.MetricNameNumGC, model.Gauge(m.NumGC).String())
	a.post(ctx, gauge, model.MetricNameOtherSys, model.Gauge(m.OtherSys).String())
	a.post(ctx, gauge, model.MetricNamePauseTotalNs, model.Gauge(m.PauseTotalNs).String())
	a.post(ctx, gauge, model.MetricNameStackInuse, model.Gauge(m.StackInuse).String())
	a.post(ctx, gauge, model.MetricNameStackSys, model.Gauge(m.StackSys).String())
	a.post(ctx, gauge, model.MetricNameSys, model.Gauge(m.Sys).String())
	a.post(ctx, gauge, model.MetricNameRandomValue, model.Gauge(a.data.randomValue).String())
	a.post(ctx, counter, model.MetricNamePollCount, model.Counter(a.data.pollCount).String())

	// TODO: Wait posts before return?
}

func NewAgent(config Config) *Agent {
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
