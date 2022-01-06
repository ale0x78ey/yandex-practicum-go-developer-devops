package agent

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMinDuration(t *testing.T) {
	tests := []struct {
		name string
		d1   time.Duration
		d2   time.Duration
		want time.Duration
	}{
		{
			name: "d1 is less than d2",
			d1:   1 * time.Second,
			d2:   2 * time.Second,
			want: 1 * time.Second,
		},
		{
			name: "d1 is more than d2",
			d1:   3 * time.Second,
			d2:   2 * time.Second,
			want: 2 * time.Second,
		},
		{
			name: "d1 is equal to d2",
			d1:   5 * time.Second,
			d2:   5 * time.Second,
			want: 5 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			minDur := minDuration(tt.d1, tt.d2)
			assert.Equal(t, tt.want, minDur)
		})
	}
}
