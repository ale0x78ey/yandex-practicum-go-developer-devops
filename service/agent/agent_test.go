package agent

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name           string
		pollInterval   time.Duration
		reportInterval time.Duration
		wantErr        bool
	}{
		{
			name:           "positive PollInterval and ReportInterval",
			pollInterval:   1,
			reportInterval: 1,
			wantErr:        false,
		},
		{
			name:           "non-positive PollInterval",
			pollInterval:   0,
			reportInterval: 1,
			wantErr:        true,
		},
		{
			name:           "non-positive ReportInterval",
			pollInterval:   1,
			reportInterval: 0,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := agent{
				config: AgentConfig{
					PollInterval:   tt.pollInterval,
					ReportInterval: tt.reportInterval,
				},
			}
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			err := a.Run(ctx)
			if !tt.wantErr {
				require.NoError(t, err)
				return
			}
			assert.Error(t, err)
		})
	}
}

func TestPollMetrics_PollCount(t *testing.T) {
	tests := []struct {
		name      string
		pollCount int64
		want      int64
	}{
		{
			name:      "pollCount 0",
			pollCount: 0,
			want:      1,
		},
		{
			name:      "pollCount 1",
			pollCount: 1,
			want:      2,
		},
		{
			name:      "pollCount 2",
			pollCount: 2,
			want:      3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := agent{
				data: metrics{
					pollCount: tt.pollCount,
				},
			}
			a.pollMetrics()
			assert.Equal(t, tt.want, a.data.pollCount)
		})
	}
}

func TestPost(t *testing.T) {
}

func TestSend(t *testing.T) {
}

func TestSendMetrics(t *testing.T) {
}

func TestNewAgent(t *testing.T) {
	_, err := NewAgent(AgentConfig{})
	assert.NoError(t, err)
}
