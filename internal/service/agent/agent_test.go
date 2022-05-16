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
		name                string
		pollInterval        time.Duration
		reportInterval      time.Duration
		postWorkersPoolSize int
		wantErr             bool
	}{
		{
			name:                "positive PollInterval and ReportInterval",
			pollInterval:        1 * time.Second,
			reportInterval:      1,
			postWorkersPoolSize: 1,
			wantErr:             false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := NewAgent(Config{
				PollInterval:        tt.pollInterval,
				ReportInterval:      tt.reportInterval,
				PostWorkersPoolSize: tt.postWorkersPoolSize,
			}, "")
			require.Nil(t, err)
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
			defer cancel()
			err = a.Run(ctx)
			if !tt.wantErr {
				require.NoError(t, err)
				return
			}
			assert.Error(t, err)
		})
	}
}
