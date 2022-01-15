package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetric_MetricFromGauge(t *testing.T) {
	type want struct {
		Name        string
		Type        MetricType
		GaugeValue  Gauge
		StringValue string
	}
	tests := []struct {
		name        string
		metricName  string
		metricValue Gauge
		want        want
	}{
		{
			name:        "gauge metric name 1",
			metricName:  "metric1",
			metricValue: Gauge(1.05),
			want: want{
				Name:        "metric1",
				Type:        MetricTypeGauge,
				GaugeValue:  Gauge(1.05),
				StringValue: "1.05",
			},
		},
		{
			name:        "gauge metric name 2",
			metricName:  "metric2",
			metricValue: Gauge(2),
			want: want{
				Name:        "metric2",
				Type:        MetricTypeGauge,
				GaugeValue:  Gauge(2),
				StringValue: "2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric := MetricFromGauge(tt.metricName, tt.metricValue)
			assert.Equal(t, tt.want.Name, metric.Name)
			assert.Equal(t, tt.want.Type, metric.Type)
			assert.Equal(t, tt.want.GaugeValue, metric.GaugeValue)
			assert.Equal(t, tt.want.StringValue, metric.StringValue())
		})
	}
}

func TestMetric_MetricFromCounter(t *testing.T) {
	type want struct {
		Name         string
		Type         MetricType
		CounterValue Counter
		StringValue  string
	}
	tests := []struct {
		name        string
		metricName  string
		metricValue Counter
		want        want
	}{
		{
			name:        "counter metric name 1",
			metricName:  "metric1",
			metricValue: Counter(1),
			want: want{
				Name:         "metric1",
				Type:         MetricTypeCounter,
				CounterValue: Counter(1),
				StringValue:  "1",
			},
		},
		{
			name:        "counter metric name 2",
			metricName:  "metric2",
			metricValue: Counter(2),
			want: want{
				Name:         "metric2",
				Type:         MetricTypeCounter,
				CounterValue: Counter(2),
				StringValue:  "2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric := MetricFromCounter(tt.metricName, tt.metricValue)
			assert.Equal(t, tt.want.Name, metric.Name)
			assert.Equal(t, tt.want.Type, metric.Type)
			assert.Equal(t, tt.want.CounterValue, metric.CounterValue)
			assert.Equal(t, tt.want.StringValue, metric.StringValue())
		})
	}
}

func TestMetric_MetricFromString(t *testing.T) {
	type want struct {
		Name         string
		Type         MetricType
		GaugeValue   Gauge
		CounterValue Counter
		StringValue  string
	}
	tests := []struct {
		name              string
		metricName        string
		metricType        MetricType
		metricStringValue string
		want              want
		wantErr           bool
	}{
		{
			name:              "gauge metric name 1",
			metricName:        "metric1",
			metricType:        MetricTypeGauge,
			metricStringValue: "0",
			want: want{
				Name:        "metric1",
				Type:        MetricTypeGauge,
				GaugeValue:  Gauge(0),
				StringValue: "0",
			},
			wantErr: false,
		},
		{
			name:              "gauge metric name 2",
			metricName:        "metric2",
			metricType:        MetricTypeGauge,
			metricStringValue: "1.0095",
			want: want{
				Name:        "metric2",
				Type:        MetricTypeGauge,
				GaugeValue:  Gauge(1.0095),
				StringValue: "1.0095",
			},
			wantErr: false,
		},
		{
			name:              "counter metric name 3",
			metricName:        "metric3",
			metricType:        MetricTypeCounter,
			metricStringValue: "0",
			want: want{
				Name:         "metric3",
				Type:         MetricTypeCounter,
				CounterValue: Counter(0),
				StringValue:  "0",
			},
			wantErr: false,
		},
		{
			name:              "counter metric name 4",
			metricName:        "metric4",
			metricType:        MetricTypeCounter,
			metricStringValue: "99999999",
			want: want{
				Name:         "metric4",
				Type:         MetricTypeCounter,
				CounterValue: Counter(99999999),
				StringValue:  "99999999",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric, err := MetricFromString(tt.metricName, tt.metricType, tt.metricStringValue)
			if !tt.wantErr {
				require.NoError(t, err)
				assert.Equal(t, tt.want.Name, metric.Name)
				assert.Equal(t, tt.want.Type, metric.Type)
				assert.Equal(t, tt.want.StringValue, metric.StringValue())

				switch tt.metricType {
				case MetricTypeGauge:
					assert.Equal(t, tt.want.GaugeValue, metric.GaugeValue)
				case MetricTypeCounter:
					assert.Equal(t, tt.want.CounterValue, metric.CounterValue)
				}
				return
			}

			assert.Error(t, err)
		})
	}
}

func TestMetricTypeValidate(t *testing.T) {
	tests := []struct {
		name    string
		value   MetricType
		wantErr bool
	}{
		{
			name:    "Valid MetricType",
			value:   MetricTypeCounter,
			wantErr: false,
		},
		{
			name:    "Empty MetricType",
			value:   MetricType(""),
			wantErr: true,
		},
		{
			name:    "Invalid MetricType",
			value:   MetricType("abrakadabra"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.value.Validate()
			if !tt.wantErr {
				require.NoError(t, err)
				return
			}
			assert.Error(t, err)
		})
	}
}

func TestGauge_String(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		want  string
	}{
		{
			name:  "Gauge 0",
			value: 0,
			want:  "0",
		},
		{
			name:  "Gauge minus 1",
			value: -1,
			want:  "-1",
		},
		{
			name:  "Gauge 1.5",
			value: 1.507,
			want:  "1.507",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := Gauge(tt.value)
			assert.Equal(t, tt.want, value.String())
		})
	}
}

func TestGauge_GaugeFromString(t *testing.T) {
}

func TestCounter_String(t *testing.T) {
	tests := []struct {
		name  string
		value int64
		want  string
	}{
		{
			name:  "Counter 0",
			value: 0,
			want:  "0",
		},
		{
			name:  "Counter minus 1",
			value: -1,
			want:  "-1",
		},
		{
			name:  "Counter 100",
			value: 100,
			want:  "100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := Counter(tt.value)
			assert.Equal(t, tt.want, value.String())
		})
	}
}

func TestCounter_CounterFromString(t *testing.T) {
}
