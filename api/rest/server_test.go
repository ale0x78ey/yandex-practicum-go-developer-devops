package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	_, err := NewServer()
	assert.NoError(t, err)
}
