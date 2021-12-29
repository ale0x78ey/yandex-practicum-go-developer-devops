package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_SaveMetric(t *testing.T) {
}

func TestServer_LoadMetric(t *testing.T) {
}

func TestNewServer(t *testing.T) {
	assert.NotNil(t, NewServer())
}
