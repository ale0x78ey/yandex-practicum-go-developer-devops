package agent

import (
	"log"
	"math/rand"
	"runtime"
)

type Gauge float64

type Counter int64

type Metrics struct {
	// MemStats is the memory allocator statistics.
	MemStats runtime.MemStats

	// PollCount is the number of previous polls.
	PollCount Counter

	// RandomValue is just a random value.
	RandomValue Gauge
}

func makeMetrics(pollCount Counter) Metrics {
	metrics := Metrics{
		PollCount:   pollCount,
		RandomValue: Gauge(rand.Float64()),
	}
	runtime.ReadMemStats(&metrics.MemStats)
	return metrics
}

type MetricsConsumer interface {
	Consume(*Metrics)
}

type MetricsConsumerFunc func(*Metrics)

func (f MetricsConsumerFunc) Consume(metrics *Metrics) {
	log.Print("Metrics.Consume")
	f(metrics)
}
