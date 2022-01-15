package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextKeyString(t *testing.T) {
	tests := []struct {
		name       string
		contextKey string
		want       string
	}{
		{
			name:       "new context key 1",
			contextKey: "test-ctx-01",
			want:       "yandex-practicum-devops-test-ctx-01",
		},
		{
			name:       "new context key 1",
			contextKey: "test-ctx-02",
			want:       "yandex-practicum-devops-test-ctx-02",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contextKey := ContextKey(tt.contextKey)
			assert.Equal(t, tt.want, contextKey.String())
		})
	}
}
