package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetric_MetricFromGauge(t *testing.T) {
	type want struct {
		ID          MetricName
		MType       MetricType
		Value       Gauge
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
				ID:          MetricName("metric1"),
				MType:       MetricTypeGauge,
				Value:       Gauge(1.05),
				StringValue: "1.05",
			},
		},
		{
			name:        "gauge metric name 2",
			metricName:  "metric2",
			metricValue: Gauge(2),
			want: want{
				ID:          MetricName("metric2"),
				MType:       MetricTypeGauge,
				Value:       Gauge(2),
				StringValue: "2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric := MetricFromGauge(tt.metricName, tt.metricValue)
			assert.Equal(t, tt.want.ID, metric.ID)
			assert.Equal(t, tt.want.MType, metric.MType)
			assert.Equal(t, tt.want.Value, *metric.Value)
			assert.Equal(t, tt.want.StringValue, metric.String())
		})
	}
}

func TestMetric_MetricFromCounter(t *testing.T) {
	type want struct {
		ID          MetricName
		MType       MetricType
		Delta       Counter
		StringValue string
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
				ID:          MetricName("metric1"),
				MType:       MetricTypeCounter,
				Delta:       Counter(1),
				StringValue: "1",
			},
		},
		{
			name:        "counter metric name 2",
			metricName:  "metric2",
			metricValue: Counter(2),
			want: want{
				ID:          MetricName("metric2"),
				MType:       MetricTypeCounter,
				Delta:       Counter(2),
				StringValue: "2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric := MetricFromCounter(tt.metricName, tt.metricValue)
			assert.Equal(t, tt.want.ID, metric.ID)
			assert.Equal(t, tt.want.MType, metric.MType)
			assert.Equal(t, tt.want.Delta, *metric.Delta)
			assert.Equal(t, tt.want.StringValue, metric.String())
		})
	}
}

func TestMetric_MetricFromString(t *testing.T) {
	type want struct {
		ID          MetricName
		MType       MetricType
		Value       Gauge
		Delta       Counter
		StringValue string
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
				ID:          MetricName("metric1"),
				MType:       MetricTypeGauge,
				Value:       Gauge(0),
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
				ID:          MetricName("metric2"),
				MType:       MetricTypeGauge,
				Value:       Gauge(1.0095),
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
				ID:          MetricName("metric3"),
				MType:       MetricTypeCounter,
				Delta:       Counter(0),
				StringValue: "0",
			},
			wantErr: false,
		},
		{
			name:              "counter metric name 4",
			metricName:        "metric4",
			metricType:        MetricTypeCounter,
			metricStringValue: "99999999",
			want: want{
				ID:          MetricName("metric4"),
				MType:       MetricTypeCounter,
				Delta:       Counter(99999999),
				StringValue: "99999999",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric, err := MetricFromString(tt.metricName, tt.metricType, tt.metricStringValue)
			if !tt.wantErr {
				require.NoError(t, err)
				assert.Equal(t, tt.want.ID, metric.ID)
				assert.Equal(t, tt.want.MType, metric.MType)
				assert.Equal(t, tt.want.StringValue, metric.String())

				switch tt.metricType {
				case MetricTypeGauge:
					assert.Equal(t, tt.want.Value, *metric.Value)
				case MetricTypeCounter:
					assert.Equal(t, tt.want.Delta, *metric.Delta)
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
