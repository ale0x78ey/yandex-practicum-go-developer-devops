package model

import (
	"strconv"
)

type Metric interface {
	Type() string
	Name() string
	StringValue() string
}

type Gauge struct {
	Value float64
	name  string
}

type Counter struct {
	Value int64
	name  string
}

func (g Gauge) Type() string {
	return "gauge"
}

func (g Gauge) Name() string {
	return g.name
}

func (g Gauge) StringValue() string {
	return strconv.FormatFloat(g.Value, 'E', -1, 64)
}

func GaugeFromFloat64(name string, value float64) Gauge {
	return Gauge{name: name, Value: value}
}

func GaugeFromUInt32(name string, value uint32) Gauge {
	return Gauge{name: name, Value: float64(value)}
}

func GaugeFromUInt64(name string, value uint64) Gauge {
	return Gauge{name: name, Value: float64(value)}
}

func (c Counter) Type() string {
	return "counter"
}

func (c Counter) Name() string {
	return c.name
}

func (c Counter) StringValue() string {
	return strconv.FormatInt(c.Value, 10)
}

func CounterFromInt64(name string, value int64) Counter {
	return Counter{name: name, Value: value}
}
