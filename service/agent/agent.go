package agent

import (
	"context"
	"errors"
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
	PostTimeout         time.Duration
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
	wg     sync.WaitGroup
}

func NewAgent(config Config) (*Agent, error) {
	if config.PostTimeout <= 0 {
		config.PostTimeout = minDuration(config.ReportInterval, config.PollInterval)
	}

	t := &http.Transport{}
	t.MaxIdleConns = config.MaxIdleConns
	t.MaxIdleConnsPerHost = config.MaxIdleConnsPerHost

	httpClient := &http.Client{
		Transport: t,
	}

	client := resty.NewWithClient(httpClient)
	if client == nil {
		return nil, errors.New("resty client wasn't created")
	}

	client.
		SetRetryCount(config.RetryCount).
		SetRetryWaitTime(config.RetryWaitTime).
		SetRetryMaxWaitTime(config.RetryMaxWaitTime)

	a := &Agent{
		config: config,
		client: client,
	}

	return a, nil
}

func (a *Agent) Run(ctx context.Context) error {
	if a.config.PollInterval <= 0 {
		return fmt.Errorf("invalid non-positive PollInterval=%v",
			a.config.PollInterval)
	}
	if a.config.ReportInterval <= 0 {
		return fmt.Errorf("invalid non-positive ReportInterval=%v",
			a.config.ReportInterval)
	}

	pollTicker := time.NewTicker(a.config.PollInterval)
	defer pollTicker.Stop()

	sendTicker := time.NewTicker(a.config.ReportInterval)
	defer sendTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			a.pollMetrics()
		case <-sendTicker.C:
			a.postMetrics(ctx)
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

func (a *Agent) postMetrics(ctx context.Context) {
	ctx2, cancel := context.WithTimeout(ctx, a.config.PostTimeout)
	defer cancel()

	m := &a.data.memStats

	a.post(ctx2, model.MetricFromGauge("Alloc", model.Gauge(m.Alloc)))
	a.post(ctx2, model.MetricFromGauge("BuckHashSys", model.Gauge(m.BuckHashSys)))
	a.post(ctx2, model.MetricFromGauge("Frees", model.Gauge(m.Frees)))
	a.post(ctx2, model.MetricFromGauge("GCCPUFraction", model.Gauge(m.GCCPUFraction)))
	a.post(ctx2, model.MetricFromGauge("GCSys", model.Gauge(m.GCSys)))
	a.post(ctx2, model.MetricFromGauge("HeapAlloc", model.Gauge(m.HeapAlloc)))
	a.post(ctx2, model.MetricFromGauge("HeapIdle", model.Gauge(m.HeapIdle)))
	a.post(ctx2, model.MetricFromGauge("HeapInuse", model.Gauge(m.HeapInuse)))
	a.post(ctx2, model.MetricFromGauge("HeapObjects", model.Gauge(m.HeapObjects)))
	a.post(ctx2, model.MetricFromGauge("HeapReleased", model.Gauge(m.HeapReleased)))
	a.post(ctx2, model.MetricFromGauge("HeapSys", model.Gauge(m.HeapSys)))
	a.post(ctx2, model.MetricFromGauge("LastGC", model.Gauge(m.LastGC)))
	a.post(ctx2, model.MetricFromGauge("Lookups", model.Gauge(m.Lookups)))
	a.post(ctx2, model.MetricFromGauge("MCacheInuse", model.Gauge(m.MCacheInuse)))
	a.post(ctx2, model.MetricFromGauge("MCacheSys", model.Gauge(m.MCacheSys)))
	a.post(ctx2, model.MetricFromGauge("MSpanInuse", model.Gauge(m.MSpanInuse)))
	a.post(ctx2, model.MetricFromGauge("MSpanSys", model.Gauge(m.MSpanSys)))
	a.post(ctx2, model.MetricFromGauge("Mallocs", model.Gauge(m.Mallocs)))
	a.post(ctx2, model.MetricFromGauge("NextGC", model.Gauge(m.NextGC)))
	a.post(ctx2, model.MetricFromGauge("NumForcedGC", model.Gauge(m.NumForcedGC)))
	a.post(ctx2, model.MetricFromGauge("NumGC", model.Gauge(m.NumGC)))
	a.post(ctx2, model.MetricFromGauge("OtherSys", model.Gauge(m.OtherSys)))
	a.post(ctx2, model.MetricFromGauge("PauseTotalNs", model.Gauge(m.PauseTotalNs)))
	a.post(ctx2, model.MetricFromGauge("StackInuse", model.Gauge(m.StackInuse)))
	a.post(ctx2, model.MetricFromGauge("StackSys", model.Gauge(m.StackSys)))
	a.post(ctx2, model.MetricFromGauge("Sys", model.Gauge(m.Sys)))
	a.post(ctx2, model.MetricFromGauge("RandomValue", model.Gauge(a.data.randomValue)))
	a.post(ctx2, model.MetricFromCounter("PollCount", model.Counter(a.data.pollCount)))

	a.wg.Wait()
}

func (a *Agent) post(ctx context.Context, metric model.Metric) {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		request := a.client.R().
			SetContext(ctx).
			SetHeader("content-type", "text/plain").
			SetPathParams(map[string]string{
				"host":        a.config.ServerHost,
				"port":        a.config.ServerPort,
				"metricName":  metric.Name,
				"metricType":  metric.Type.String(),
				"metricValue": metric.StringValue(),
			})

		request.Post(metricPostURL)
	}()
}
