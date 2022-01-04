package model

import (
	"fmt"
	"strconv"
)

type (
	Metric struct {
		Name         string
		Type         MetricType
		GaugeValue   Gauge
		CounterValue Counter
	}

	MetricType string
	Gauge      float64
	Counter    int64
)

func MetricFromGauge(name string, value Gauge) Metric {
	return Metric{
		Name:       name,
		Type:       MetricTypeGauge,
		GaugeValue: value,
	}
}

func MetricFromCounter(name string, value Counter) Metric {
	return Metric{
		Name:         name,
		Type:         MetricTypeCounter,
		CounterValue: value,
	}
}

func MetricFromString(metricName string, metricType MetricType, value string) (Metric, error) {
	switch metricType {
	case MetricTypeGauge:
		gaugeValue, err := GaugeFromString(value)
		if err != nil {
			return Metric{}, err
		}
		return MetricFromGauge(metricName, gaugeValue), nil

	case MetricTypeCounter:
		counterValue, err := CounterFromString(value)
		if err != nil {
			return Metric{}, err
		}
		return MetricFromCounter(metricName, counterValue), nil

	default:
		return Metric{}, fmt.Errorf("unkown MetricType: %s", metricType)
	}
}

func (m *Metric) StringValue() string {
	switch m.Type {
	case MetricTypeGauge:
		return m.GaugeValue.String()
	case MetricTypeCounter:
		return m.CounterValue.String()
	default:
		return fmt.Sprintf("<unkown>")
	}
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

func (g Gauge) String() string {
	return strconv.FormatFloat(float64(g), 'f', -1, 64)
}

func GaugeFromString(value string) (Gauge, error) {
	g, err := strconv.ParseFloat(value, 64)
	return Gauge(g), err
}

func (c Counter) String() string {
	return strconv.FormatInt(int64(c), 10)
}

func CounterFromString(value string) (Counter, error) {
	c, err := strconv.Atoi(value)
	return Counter(c), err
}
