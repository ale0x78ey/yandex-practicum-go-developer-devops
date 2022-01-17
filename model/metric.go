package model

import (
	"fmt"
	"strconv"
)

type (
	Metric struct {
		ID    string     `json:"id"`
		MType MetricType `json:"type"`
		Delta *Counter   `json:"delta,omitempty"`
		Value *Gauge     `json:"value,omitempty"`
	}

	MetricType string
	Gauge      float64
	Counter    int64
)

func MetricFromGauge(id string, value Gauge) Metric {
	return Metric{
		ID:    id,
		MType: MetricTypeGauge,
		Value: &value,
	}
}

func MetricFromCounter(id string, value Counter) Metric {
	return Metric{
		ID:    id,
		MType: MetricTypeCounter,
		Delta: &value,
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
		return Metric{}, fmt.Errorf("unkown MType: %s", metricType)
	}
}

func (m Metric) Validate() error {
	switch m.MType {
	case MetricTypeGauge:
		if m.Value == nil {
			return fmt.Errorf("invalid Value == nil for MType: %s", m.MType)
		}
	case MetricTypeCounter:
		if m.Delta == nil {
			return fmt.Errorf("invalid Delta == nil for MType: %s", m.MType)
		}
	default:
		return fmt.Errorf("unkown MType: %s", m.MType)
	}
	return nil
}

func (m Metric) String() string {
	switch m.MType {
	case MetricTypeGauge:
		return m.Value.String()
	case MetricTypeCounter:
		return m.Delta.String()
	default:
		return ""
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
		return fmt.Errorf("unkown MType: %s", t)
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
