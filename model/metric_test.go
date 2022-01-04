package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
