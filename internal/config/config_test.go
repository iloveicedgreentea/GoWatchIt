package config

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	s := GetString("homeAssistant.port")
	assert.Equal(t, "8123", s)
}