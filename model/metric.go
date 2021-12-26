package model

import (
	"fmt"
	"strconv"
)

type (
	Gauge      int64
	Counter    float64
	MetricName string
	MetricType string
)

func (g Gauge) String() string {
	return strconv.FormatFloat(float64(g), 'E', -1, 64)
}

func (c Counter) String() string {
	return strconv.FormatInt(int64(c), 10)
}

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

func (t MetricType) String() string {
	return string(t)
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

func (n MetricName) String() string {
	return string(n)
}
