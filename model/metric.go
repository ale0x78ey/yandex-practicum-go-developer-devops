package model

import (
	"fmt"
	"strconv"
)

type (
	Metric interface {
		Type() MetricType
		Name() MetricName
		StringValue() string
	}

	Gauge struct {
		Value float64
		name  MetricName
	}

	Counter struct {
		Value int64
		name  MetricName
	}

	MetricName string
	MetricType string
)

const (
	MetricTypeGauge   MetricType = "gauge"
	MetricTypeCounter MetricType = "counter"
)

func (t MetricType) Validate() error {
	switch t {
	case MetricTypeGauge, MetricTypeCounter:
		return nil
	default:
		return fmt.Errorf("unkown MetricType: %s", t)
	}
}

const (
	MetricNameAlloc         MetricName = "Alloc"
	MetricNameBuckHashSys   MetricName = "BuckHashSys"
	MetricNameFrees         MetricName = "Frees"
	MetricNameGCCPUFraction MetricName = "GCCPUFraction"
	MetricNameGCSys         MetricName = "GCSys"
	MetricNameHeapAlloc     MetricName = "HeapAlloc"
	MetricNameHeapIdle      MetricName = "HeapIdle"
	MetricNameHeapInuse     MetricName = "HeapInuse"
	MetricNameHeapObjects   MetricName = "HeapObjects"
	MetricNameHeapReleased  MetricName = "HeapReleased"
	MetricNameHeapSys       MetricName = "HeapSys"
	MetricNameLastGC        MetricName = "LastGC"
	MetricNameLookups       MetricName = "Lookups"
	MetricNameMCacheInuse   MetricName = "MCacheInuse"
	MetricNameMCacheSys     MetricName = "MCacheSys"
	MetricNameMSpanInuse    MetricName = "MSpanInuse"
	MetricNameMSpanSys      MetricName = "MSpanSys"
	MetricNameMallocs       MetricName = "Mallocs"
	MetricNameNextGC        MetricName = "NextGC"
	MetricNameNumForcedGC   MetricName = "NumForcedGC"
	MetricNameNumGC         MetricName = "NumGC"
	MetricNameOtherSys      MetricName = "OtherSys"
	MetricNamePauseTotalNs  MetricName = "PauseTotalNs"
	MetricNameStackInuse    MetricName = "StackInuse"
	MetricNameStackSys      MetricName = "StackSys"
	MetricNameSys           MetricName = "Sys"
	MetricNameRandomValue   MetricName = "RandomValue"
	MetricNamePollCount     MetricName = "PollCount"
)

func (n MetricName) Validate() error {
	switch n {
	case
		MetricNameAlloc,
		MetricNameBuckHashSys,
		MetricNameFrees,
		MetricNameGCCPUFraction,
		MetricNameGCSys,
		MetricNameHeapAlloc,
		MetricNameHeapIdle,
		MetricNameHeapInuse,
		MetricNameHeapObjects,
		MetricNameHeapReleased,
		MetricNameHeapSys,
		MetricNameLastGC,
		MetricNameLookups,
		MetricNameMCacheInuse,
		MetricNameMCacheSys,
		MetricNameMSpanInuse,
		MetricNameMSpanSys,
		MetricNameMallocs,
		MetricNameNextGC,
		MetricNameNumForcedGC,
		MetricNameNumGC,
		MetricNameOtherSys,
		MetricNamePauseTotalNs,
		MetricNameStackInuse,
		MetricNameStackSys,
		MetricNameSys,
		MetricNameRandomValue,
		MetricNamePollCount:
		return nil
	default:
		return fmt.Errorf("unkown MetricName: %s", n)
	}
}

func (g Gauge) Type() MetricType {
	return "gauge"
}

func (g Gauge) Name() MetricName {
	return g.name
}

func (g Gauge) StringValue() string {
	return strconv.FormatFloat(g.Value, 'E', -1, 64)
}

func GaugeFromFloat64(name string, value float64) Gauge {
	return Gauge{name: MetricName(name), Value: value}
}

func GaugeFromUInt32(name string, value uint32) Gauge {
	return Gauge{name: MetricName(name), Value: float64(value)}
}

func GaugeFromUInt64(name string, value uint64) Gauge {
	return Gauge{name: MetricName(name), Value: float64(value)}
}

func (c Counter) Type() MetricType {
	return "counter"
}

func (c Counter) Name() MetricName {
	return c.name
}

func (c Counter) StringValue() string {
	return strconv.FormatInt(c.Value, 10)
}

func CounterFromInt64(name string, value int64) Counter {
	return Counter{name: MetricName(name), Value: value}
}
