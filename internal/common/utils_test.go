package common

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestInsensitiveContains(t *testing.T) {
	assert := assert.New(t)
	assert.True(InsensitiveContains("Dts-HD MA 5.1 - English - Default", "DTS-HD MA 5.1"))

}