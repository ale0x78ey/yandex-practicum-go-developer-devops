package model

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/internal/common"
)

type (
	Metric struct {
		ID    MetricName `json:"id"`
		MType MetricType `json:"type"`
		Delta *Counter   `json:"delta,omitempty"`
		Value *Gauge     `json:"value,omitempty"`
		Hash  string     `json:"hash,omitempty"`
	}

	MetricName string
	MetricType string
	Gauge      float64
	Counter    int64
)

func MetricFromGauge(id string, value Gauge) Metric {
	return Metric{
		ID:    MetricName(id),
		MType: MetricTypeGauge,
		Value: &value,
	}
}

func MetricFromCounter(id string, value Counter) Metric {
	return Metric{
		ID:    MetricName(id),
		MType: MetricTypeCounter,
		Delta: &value,
	}
}

func MetricFromString(metricName string, metricType MetricType, value string) (Metric, error) {
	var metric Metric
	switch metricType {
	case MetricTypeGauge:
		gaugeValue, err := GaugeFromString(value)
		if err != nil {
			return Metric{}, err
		}
		metric = MetricFromGauge(metricName, gaugeValue)

	case MetricTypeCounter:
		counterValue, err := CounterFromString(value)
		if err != nil {
			return Metric{}, err
		}
		metric = MetricFromCounter(metricName, counterValue)

	default:
		return Metric{}, fmt.Errorf("unknown MetricType: %s", metricType)
	}

	if err := metric.Validate(); err != nil {
		return Metric{}, err
	}

	return metric, nil
}

func (m Metric) Validate() error {
	if err := m.ID.Validate(); err != nil {
		return err
	}

	if err := m.MType.Validate(); err != nil {
		return err
	}

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
		return fmt.Errorf("unknown MetricType: %s", m.MType)
	}

	return nil
}

func (m Metric) ValidateHash(key string) (bool, error) {
	hash, err := m.ProcessHash(key)
	if err != nil {
		return false, err
	}
	return hash == m.Hash, nil
}

func (m Metric) String() string {
	switch m.MType {
	case MetricTypeGauge:
		return (*m.Value).String()
	case MetricTypeCounter:
		return (*m.Delta).String()
	default:
		return ""
	}
}

func (m Metric) ProcessHash(key string) (string, error) {
	var data string

	switch m.MType {
	case MetricTypeGauge:
		data = fmt.Sprintf("%s:%s:%f", m.ID, m.MType, float64(*m.Value))
	case MetricTypeCounter:
		data = fmt.Sprintf("%s:%s:%d", m.ID, m.MType, int64(*m.Delta))
	default:
		return "", fmt.Errorf("unkown MetricType: %s", m.MType)
	}

	hash, err := common.Hash([]byte(data), []byte(key))
	if err != nil {
		return "", err
	}

	return hash, nil
}

func (m *Metric) UpdateHash(key string) error {
	if key == "" {
		return nil
	}

	hash, err := m.ProcessHash(key)
	if err != nil {
		return err
	}

	m.Hash = hash

	return nil
}

func (t MetricName) Validate() error {
	if t == "" {
		return errors.New("invalid empty MetricName")
	}
	return nil
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
		return fmt.Errorf("unknown MetricType: %s", t)
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
